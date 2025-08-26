package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollView は、コンテンツをスクロール表示するためのコンテナウィジェットです。
// このウィジェットは、基本的なウィジェット機能(LayoutableWidget)と、
// 子要素の管理とクリッピング描画を行う内部コンテナ(container.Container)を「合成」しています。
// この設計により、ScrollView独自の複雑な更新・描画ロジックと、汎用コンテナのロジックを
// 明確に分離し、コードの堅牢性を高めています。
type ScrollView struct {
	*component.LayoutableWidget
	container         *container.Container
	layout            layout.Layout
	contentContainer  component.Widget
	vScrollBar        component.ScrollBarWidget
	scrollY           float64
	contentHeight     int
	ScrollSensitivity float64
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)
var _ container.Scroller = (*ScrollView)(nil)

// NewScrollView はScrollViewのインスタンスを生成します。
func NewScrollView() *ScrollView {
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}
	sv.LayoutableWidget = component.NewLayoutableWidget()
	sv.container = container.NewContainer()
	sv.container.SetClipsChildren(true)

	sv.Init(sv)                            // ScrollView自身のLayoutableWidgetを初期化
	sv.container.SetParent(sv)             // 内部コンテナの親をScrollView自身に設定
	sv.layout = &layout.ScrollViewLayout{} // ScrollView専用のレイアウトを設定

	vScrollBar, _ := NewScrollBarBuilder().Build()
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar) // 子要素は内部コンテナに追加される

	return sv
}

// SetContent は、スクロールさせるコンテンツコンテナを設定します。
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		sv.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content
	if content != nil {
		sv.AddChild(content)
		sv.AddChild(sv.vScrollBar) // スクロールバーが最前面に来るように再追加
	}
	sv.MarkDirty(true)
}

// Update はScrollViewの状態を更新します。
// 汎用コンテナのUpdateとは異なり、ScrollView専用のレイアウト計算を制御します。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}

	// 自身のサイズと位置を内部コンテナに常に同期させます。
	w, h := sv.GetSize()
	sv.container.SetSize(w, h)
	x, y := sv.GetPosition()
	sv.container.SetPosition(x, y)

	// ScrollView自身が再レイアウトを要求されている場合のみ、専用のレイアウトを実行します。
	if sv.NeedsRelayout() {
		if sv.layout != nil {
			sv.layout.Layout(sv)
		}
	}

	// 内部コンテナの子要素のUpdateを再帰的に呼び出します。
	// ScrollViewLayout内でコンテンツコンテナ(sv.contentContainer)のUpdateは既に
	// 呼び出されているため、二重呼び出しを避けることで効率化します。
	// スクロールバーのような他の子要素のUpdateはここで呼び出す必要があります。
	for _, child := range sv.container.GetChildren() {
		if child == sv.contentContainer {
			continue
		}
		child.Update()
	}

	if sv.IsDirty() {
		sv.ClearDirty()
	}
}

// Draw はScrollViewを描画します。描画はクリッピング機能を持つ内部コンテナに完全に委譲します。
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	if !sv.IsVisible() {
		return
	}
	sv.container.Draw(screen)
}

// HandleEvent はScrollViewのイベントを処理します。
func (sv *ScrollView) HandleEvent(e *event.Event) {
	if !e.Handled && e.Type == event.MouseScroll {
		scrollAmount := e.ScrollY * sv.ScrollSensitivity
		sv.scrollY -= scrollAmount
		sv.MarkDirty(true)
		e.Handled = true
	}
	sv.LayoutableWidget.HandleEvent(e)
}

// HitTest は、指定された座標がヒットするウィジェットを探します。
func (sv *ScrollView) HitTest(x, y int) component.Widget {
	if sv.LayoutableWidget.HitTest(x, y) != nil {
		// 自身がヒット範囲内であれば、次に内部コンテナの子要素をテストします。
		if target := sv.container.HitTest(x, y); target != nil {
			return target
		}
		return sv // 子にヒットしなければScrollView自身を返す
	}
	return nil
}

// --- メソッドの委譲 ---
func (sv *ScrollView) AddChild(child component.Widget)    { sv.container.AddChild(child) }
func (sv *ScrollView) RemoveChild(child component.Widget) { sv.container.RemoveChild(child) }
func (sv *ScrollView) GetChildren() []component.Widget    { return sv.container.GetChildren() }
func (sv *ScrollView) GetLayout() layout.Layout           { return sv.layout }
func (sv *ScrollView) SetLayout(l layout.Layout)          { sv.layout = l }
func (sv *ScrollView) GetPadding() layout.Insets          { return sv.container.GetPadding() }

// --- layout.ScrollViewer interface ---
func (sv *ScrollView) GetContentContainer() component.Widget    { return sv.contentContainer }
func (sv *ScrollView) GetVScrollBar() component.ScrollBarWidget { return sv.vScrollBar }
func (sv *ScrollView) GetScrollY() float64                      { return sv.scrollY }
func (sv *ScrollView) SetContentHeight(h int)                   { sv.contentHeight = h }
func (sv *ScrollView) SetScrollY(y float64) {
	if sv.scrollY != y {
		sv.scrollY = y
		sv.MarkDirty(false)
	}
}

// --- container.Scroller interface ---
func (sv *ScrollView) GetScrollOffset() (x, y int) { return 0, -int(sv.scrollY) }
