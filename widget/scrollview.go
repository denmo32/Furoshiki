package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"image"

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

// newScrollView はScrollViewのインスタンスを生成します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//       常にNewScrollViewBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newScrollView() (*ScrollView, error) {
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}
	sv.LayoutableWidget = component.NewLayoutableWidget()

	// NOTE: ScrollViewの依存コンポーネントを、それぞれの公開コンストラクタ/ビルダー経由で生成します。
	internalContainer, err := container.NewContainer()
	if err != nil {
		return nil, err
	}
	sv.container = internalContainer
	sv.container.SetClipsChildren(true)

	if err := sv.Init(sv); err != nil {
		return nil, err
	}
	sv.container.SetParent(sv)
	sv.layout = &layout.ScrollViewLayout{}

	// NOTE: ScrollBarの生成にビルダーを使用することで、安全な初期化を保証します。
	vScrollBar, err := NewScrollBarBuilder().Build()
	if err != nil {
		return nil, err
	}
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar)

	return sv, nil
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

// SetStyle はScrollViewと、その描画を担当する内部コンテナの両方にスタイルを設定します。
// これにより、ビルダーで設定された枠線などが正しく描画されるようになります。
func (sv *ScrollView) SetStyle(s style.Style) {
	sv.LayoutableWidget.SetStyle(s)
	if sv.container != nil {
		sv.container.SetStyle(s)
	}
}

// Update はScrollViewの状態を更新します。
// 新しいレイアウトシステムでは、レイアウト計算はトップダウンのMeasure/Arrangeパスで
// 処理されるため、このUpdateメソッドはレイアウト以外の更新（アニメーションなど）のみを行います。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}

	// 自身のサイズと位置を内部コンテナに常に同期させます。
	// これは、親のArrangeパスで設定されたScrollViewのジオメトリを、
	// 描画を担当する内部コンテナに反映させるために必要です。
	w, h := sv.GetSize()
	sv.container.SetSize(w, h)
	x, y := sv.GetPosition()
	sv.container.SetPosition(x, y)

	// レイアウトロジックはMeasure/Arrangeパスに移行しました。
	// このメソッドは子の更新のみを再帰的に呼び出します。
	for _, child := range sv.container.GetChildren() {
		child.Update()
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

// Measure implements the component.LayoutManager interface.
// It delegates the measurement to its assigned layout object.
func (sv *ScrollView) Measure(availableSize image.Point) image.Point {
	if sv.layout != nil {
		return sv.layout.Measure(sv, availableSize)
	}
	minW, minH := sv.GetMinSize()
	return image.Point{X: minW, Y: minH}
}

// Arrange implements the component.LayoutManager interface.
// It delegates the arrangement of its children to its assigned layout object.
func (sv *ScrollView) Arrange(finalBounds image.Rectangle) error {
	if sv.layout != nil {
		// The actual positioning of the ScrollView itself is done by its parent's Arrange pass.
		// This delegates the arrangement of the content and scrollbar to the ScrollViewLayout.
		return sv.layout.Arrange(sv, finalBounds)
	}
	return nil
}