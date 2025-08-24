package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/layout"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollView は、独立したコンテナとして再実装されました。
type ScrollView struct {
	*component.LayoutableWidget
	children         []component.Widget
	contentContainer component.Widget
	vScrollBar       component.ScrollBarWidget
	clipsChildren    bool
	offscreenImage   *ebiten.Image
	scrollY          float64
	contentHeight    int
	// 【追加】レイアウトマネージャを保持するためのフィールド
	layout layout.Layout
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)

// NewScrollView はScrollViewのインスタンスを生成します。
func NewScrollView() *ScrollView {
	sv := &ScrollView{
		clipsChildren: true,
	}
	sv.LayoutableWidget = component.NewLayoutableWidget(sv)
	// 【修正】sv自身に実装されたSetLayoutメソッドを呼び出します。
	sv.SetLayout(&layout.ScrollViewLayout{})

	vScrollBar, _ := NewScrollBarBuilder().Build()
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar)

	return sv
}

// SetContent は、スクロールさせるコンテンツコンテナを設定します。
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		sv.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content
	if content != nil {
		// コンテンツはスクロールバーより手前に描画するため、子リストの先頭に追加
		sv.children = append([]component.Widget{content}, sv.children...)
		content.SetParent(sv)
		sv.MarkDirty(true)
	}
}

// Update はウィジェットの状態を更新します。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}
	if sv.IsDirty() {
		if sv.NeedsRelayout() {
			// 自身のレイアウト（ScrollViewLayout）を実行
			if l := sv.GetLayout(); l != nil {
				l.Layout(sv)
			}
		}
		sv.ClearDirty()
	}
	// 子要素のUpdateを呼び出す
	for _, child := range sv.children {
		child.Update()
	}
}

// 【追加】SetLayoutは、このコンテナが使用するレイアウトマネージャを設定します。
func (sv *ScrollView) SetLayout(l layout.Layout) {
	sv.layout = l
}

// 【修正】GetLayoutは、このコンテナが保持しているレイアウトマネージャを返します。
func (sv *ScrollView) GetLayout() layout.Layout {
	return sv.layout
}

// DrawはScrollViewと子要素を描画します。
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	if !sv.IsVisible() || !sv.HasBeenLaidOut() {
		return
	}

	x, y := sv.GetPosition()
	w, h := sv.GetSize()

	// 背景描画
	component.DrawStyledBackground(screen, x, y, w, h, sv.GetStyle())

	// クリッピング描画
	if sv.clipsChildren {
		sv.drawWithClipping(screen)
	} else {
		for _, child := range sv.children {
			child.Draw(screen)
		}
	}
}

// drawWithClipping はクリッピングを行いながら描画します。
func (sv *ScrollView) drawWithClipping(screen *ebiten.Image) {
	containerX, containerY := sv.GetPosition()
	containerWidth, containerHeight := sv.GetSize()
	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	if sv.offscreenImage == nil || sv.offscreenImage.Bounds().Dx() != containerWidth || sv.offscreenImage.Bounds().Dy() != containerHeight {
		if sv.offscreenImage != nil {
			sv.offscreenImage.Deallocate()
		}
		sv.offscreenImage = ebiten.NewImage(containerWidth, containerHeight)
	}
	sv.offscreenImage.Clear()

	// 子要素をオフスクリーンに描画
	for _, child := range sv.children {
		originalX, originalY := child.GetPosition()
		relativeX := originalX - containerX
		relativeY := originalY - containerY

		child.SetPosition(relativeX, relativeY)
		child.Draw(sv.offscreenImage)
		child.SetPosition(originalX, originalY)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(containerX), float64(containerY))
	screen.DrawImage(sv.offscreenImage, opts)
}

// HandleEvent はイベントを処理します。
func (sv *ScrollView) HandleEvent(e *event.Event) {
	// 自身のカスタムハンドラをまず呼び出す
	sv.LayoutableWidget.HandleEvent(e)
	if e.Handled {
		return
	}
	// マウスホイールイベントを処理
	if e.Type == event.MouseScroll {
		scrollAmount := e.ScrollY * 20
		sv.scrollY -= scrollAmount
		sv.MarkDirty(true)
		e.Handled = true
	}
	// 子要素へのイベント伝播はLayoutableWidget.HandleEventが親を辿る形で行うのでここでは不要
}

// HitTest
func (sv *ScrollView) HitTest(x, y int) component.Widget {
	if !sv.IsVisible() {
		return nil
	}
	for i := len(sv.children) - 1; i >= 0; i-- {
		child := sv.children[i]
		if !child.IsVisible() {
			continue
		}
		if target := child.HitTest(x, y); target != nil {
			return target
		}
	}
	if sv.LayoutableWidget.HitTest(x, y) != nil {
		return sv
	}
	return nil
}

// --- component.Container interface implementation ---
func (sv *ScrollView) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	if oldParent := child.GetParent(); oldParent != nil {
		oldParent.RemoveChild(child)
	}
	child.SetParent(sv)
	sv.children = append(sv.children, child)
	sv.MarkDirty(true)
}

func (sv *ScrollView) RemoveChild(child component.Widget) {
	if child == nil {
		return
	}
	for i, currentChild := range sv.children {
		if currentChild == child {
			sv.children = append(sv.children[:i], sv.children[i+1:]...)
			child.SetParent(nil)
			sv.MarkDirty(true)
			return
		}
	}
}

func (sv *ScrollView) GetChildren() []component.Widget {
	return sv.children
}

// Cleanup
func (sv *ScrollView) Cleanup() {
	for _, child := range sv.children {
		child.Cleanup()
	}
	sv.children = nil
	if sv.offscreenImage != nil {
		sv.offscreenImage.Deallocate()
		sv.offscreenImage = nil
	}
	sv.LayoutableWidget.Cleanup()
}

// --- layout.ScrollViewer interface implementation ---
func (sv *ScrollView) GetContentContainer() component.Widget {
	return sv.contentContainer
}
func (sv *ScrollView) GetVScrollBar() component.ScrollBarWidget {
	return sv.vScrollBar
}
func (sv *ScrollView) GetScrollY() float64 {
	return sv.scrollY
}
func (sv *ScrollView) SetScrollY(y float64) {
	sv.scrollY = y
}
func (sv *ScrollView) SetContentHeight(h int) {
	sv.contentHeight = h
}
func (sv *ScrollView) GetPadding() layout.Insets {
	style := sv.GetStyle()
	if style.Padding != nil {
		return layout.Insets{
			Top:    style.Padding.Top,
			Right:  style.Padding.Right,
			Bottom: style.Padding.Bottom,
			Left:   style.Padding.Left,
		}
	}
	return layout.Insets{}
}