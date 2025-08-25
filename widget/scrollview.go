package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollView は、コンテンツをスクロール表示するためのコンテナウィジェットです。
// 【リファクタリング】基本的なウィジェット機能を component.LayoutableWidget から継承し、
// 子要素の管理は内部の container.Container インスタンスに委譲します。
// これにより、埋め込みによる複雑なメソッドオーバーライドを回避し、コードの堅牢性を高めます。
type ScrollView struct {
	*component.LayoutableWidget // ポインタとして埋め込む

	// 内部コンテナ: 子ウィジェット（コンテンツやスクロールバー）の管理と描画を担当
	container *container.Container
	// 【修正】レイアウトマネージャを自身で保持する
	layout layout.Layout

	// スクロールビュー固有のプロパティ
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
	// ScrollView自身のインスタンスを生成
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}

	// LayoutableWidgetを正しく初期化
	sv.LayoutableWidget = component.NewLayoutableWidget()

	// 1. 内部で利用するコンテナを生成
	internalContainer := container.NewContainer()
	internalContainer.SetClipsChildren(true)
	sv.container = internalContainer

	// 2. ScrollView自身のLayoutableWidgetを初期化
	sv.LayoutableWidget.Init(sv)

	// 3. 内部コンテナの親をScrollView自身に設定
	sv.container.SetParent(sv)

	// 4. 【修正】ScrollView自身のレイアウトマネージャを設定
	sv.layout = &layout.ScrollViewLayout{}

	// 5. 垂直スクロールバーを作成し、内部コンテナに追加
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
		sv.AddChild(content)
		sv.AddChild(sv.vScrollBar)
	}
	sv.MarkDirty(true)
}

// Update はScrollViewの状態を更新します。
// 【重要修正】このメソッドは、コンテナのUpdateロジックに依存せず、
// ScrollView専用のレイアウト計算と子の更新を直接制御します。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}

	// 自身のサイズと位置を内部コンテナに反映
	w, h := sv.GetSize()
	sv.container.SetSize(w, h)
	x, y := sv.GetPosition()
	sv.container.SetPosition(x, y)

	// ScrollViewがダーティなら、自身のレイアウトを実行
	if sv.NeedsRelayout() {
		if sv.layout != nil {
			// Layoutメソッドには、インターフェースを満たすsv自身を渡す
			sv.layout.Layout(sv)
		}
	}

	// 内部コンテナのレイアウト処理は行わず、子要素のUpdateのみを呼び出す
	for _, child := range sv.container.GetChildren() {
		child.Update()
	}

	// 最後に自身のダーティフラグをクリア
	if sv.IsDirty() {
		sv.ClearDirty()
	}
}

// Draw はScrollViewを描画します。
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	if !sv.IsVisible() {
		return
	}
	// 描画はクリッピング機能を持つ内部コンテナに完全に委譲
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
		if target := sv.container.HitTest(x, y); target != nil {
			return target
		}
		return sv
	}
	return nil
}

// --- component.Container interface implementation (Delegation) ---

func (sv *ScrollView) AddChild(child component.Widget) {
	sv.container.AddChild(child)
}

func (sv *ScrollView) RemoveChild(child component.Widget) {
	sv.container.RemoveChild(child)
}

func (sv *ScrollView) GetChildren() []component.Widget {
	return sv.container.GetChildren()
}

// 【修正】GetLayout/SetLayoutは自身のフィールドを操作
func (sv *ScrollView) GetLayout() layout.Layout {
	return sv.layout
}

func (sv *ScrollView) SetLayout(l layout.Layout) {
	sv.layout = l
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
	if sv.scrollY != y {
		sv.scrollY = y
		sv.MarkDirty(false)
	}
}
func (sv *ScrollView) SetContentHeight(h int) {
	sv.contentHeight = h
}

func (sv *ScrollView) GetPadding() layout.Insets {
	return sv.container.GetPadding()
}

// --- container.Scroller interface implementation ---

func (sv *ScrollView) GetScrollOffset() (x, y int) {
	return 0, -int(sv.scrollY)
}