package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"image"
	"log"
)

// ScrollView is a container widget that allows scrolling its content.
type ScrollView struct {
	*component.WidgetCore
	*component.Appearance
	*component.Interaction

	container         *container.Container
	layout            layout.Layout
	contentContainer  component.Widget
	vScrollBar        component.ScrollBarWidget
	scrollY           float64
	contentHeight     int
	ScrollSensitivity float64
}

// --- Interface implementation verification ---
var _ component.StandardWidget = (*ScrollView)(nil)
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)
var _ container.Scroller = (*ScrollView)(nil)
var _ event.EventTarget = (*ScrollView)(nil)

// newScrollView creates a new component-based ScrollView.
func newScrollView() (*ScrollView, error) {
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}

	sv.WidgetCore = component.NewWidgetCore(sv)
	sv.Appearance = component.NewAppearance(sv)
	sv.Interaction = component.NewInteraction(sv)

	internalContainer, err := container.NewContainer()
	if err != nil {
		return nil, err
	}
	sv.container = internalContainer
	sv.container.SetClipsChildren(true)
	sv.AddChild(sv.container) // 内部コンテナを子として追加

	sv.layout = &layout.ScrollViewLayout{}

	vScrollBar, err := NewScrollBarBuilder().Build()
	if err != nil {
		return nil, err
	}
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar) // スクロールバーも子として追加

	sv.AddEventHandler(event.MouseScroll, sv.onMouseScroll)

	return sv, nil
}

func (sv *ScrollView) onMouseScroll(e *event.Event) event.Propagation {
	scrollAmount := e.ScrollY * sv.ScrollSensitivity
	sv.scrollY -= scrollAmount
	sv.MarkDirty(true)
	return event.StopPropagation
}

// SetContent sets the scrollable content container.
// The content is added to the internal clipped container, not the ScrollView itself.
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		// 古いコンテンツを内部コンテナから削除
		sv.container.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content
	if content != nil {
		// 新しいコンテンツを内部コンテナに追加
		sv.container.AddChild(content)
	}
	sv.MarkDirty(true)
}

func (sv *ScrollView) SetStyle(s style.Style) {
	sv.Appearance.SetStyle(s)
	if sv.container != nil {
		// 内部コンテナのスタイルも同期させる
		sv.container.SetStyle(s)
	}
}

func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}

	if sv.NeedsRelayout() {
		if sv.layout != nil {
			if err := sv.layout.Layout(sv); err != nil {
				log.Printf("Error during ScrollView layout: %v", err)
			}
		}
	}

	// ScrollView自身の子（内部コンテナとスクロールバー）のUpdateを呼び出す
	for _, child := range sv.GetChildren() {
		child.Update()
	}

	if sv.IsDirty() {
		sv.ClearDirty()
	}
}

func (sv *ScrollView) Draw(info component.DrawInfo) {
	if !sv.IsVisible() || !sv.HasBeenLaidOut() {
		return
	}
	// ScrollView自体は何も描画せず、子要素（内部コンテナとスクロールバー）の描画に任せます。
	// これにより、クリッピングが正しく機能します。
	for _, child := range sv.GetChildren() {
		child.Draw(info)
	}
}

func (sv *ScrollView) Cleanup() {
	for _, child := range sv.GetChildren() {
		child.Cleanup()
	}
	sv.GetNode().ClearChildren()
	sv.SetParent(nil)
}

// MarkDirty marks the widget as needing a redraw or relayout.
// ScrollView is a layout boundary, so it doesn't propagate the dirty flag up.
func (sv *ScrollView) MarkDirty(relayout bool) {
	sv.Dirty.MarkDirty(relayout)
}

// SetPosition sets the position for the ScrollView. Its children's positions
// will be updated by the layout system.
func (sv *ScrollView) SetPosition(x, y int) {
	if !sv.HasBeenLaidOut() {
		sv.SetLaidOut(true)
	}
	if currX, currY := sv.GetPosition(); currX != x || currY != y {
		sv.Transform.SetPosition(x, y)
		// [BUGFIX] ScrollViewの位置が変更された際に、内部のクリッピング用コンテナの位置も
		//          同期させる必要があります。これを怠ると、ScrollViewが(0,0)以外の場所に
		//          配置された場合に、コンテンツの描画位置がずれて見えなくなります。
		if sv.container != nil {
			sv.container.SetPosition(x, y)
		}
		// 子要素の位置はレイアウトシステムによって管理されるため、
		// ここで子の位置を直接更新する代わりに再レイアウトを要求します。
		sv.MarkDirty(true)
	}
}

// SetSize sets the size for the ScrollView. The internal container's size
// will be updated as well.
func (sv *ScrollView) SetSize(width, height int) {
	if w, h := sv.GetSize(); w != width || h != height {
		sv.Transform.SetSize(width, height)
		// ScrollViewのサイズが変わると、内部コンテナのサイズも追従する必要があります。
		if sv.container != nil {
			sv.container.SetSize(width, height)
		}
		sv.MarkDirty(true)
	}
}

func (sv *ScrollView) HitTest(x, y int) component.Widget {
	if !sv.IsVisible() || sv.IsDisabled() {
		return nil
	}
	wx, wy := sv.GetPosition()
	wwidth, wheight := sv.GetSize()
	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if !rect.Empty() && (image.Point{X: x, Y: y}.In(rect)) {
		// ヒットテストは描画順と逆（手前が先）に行うのが一般的なので、
		// スクロールバーを先にテストします。
		if sv.vScrollBar != nil {
			if target := sv.vScrollBar.HitTest(x, y); target != nil {
				return target
			}
		}
		// 次に内部コンテナをテストします。
		if sv.container != nil {
			if target := sv.container.HitTest(x, y); target != nil {
				return target
			}
		}
		// どこにもヒットしなければ、ScrollView自身を返します（スクロールイベントのため）。
		return sv
	}
	return nil
}

func (sv *ScrollView) HandleEvent(e *event.Event) {
	sv.Interaction.TriggerHandlers(e)

	if e != nil && !e.Handled && sv.GetParent() != nil {
		if processor, ok := sv.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// --- Container interface implementation ---
// ScrollViewは特殊なコンテナです。AddChild/RemoveChild/GetChildrenは
// ScrollView自身の直接の子（内部コンテナとスクロールバー）を操作します。
// スクロールされるコンテンツはSetContent経由で内部コンテナに追加されます。
func (sv *ScrollView) AddChild(child component.Widget)    { sv.Node.AddChild(child) }
func (sv *ScrollView) RemoveChild(child component.Widget) { sv.Node.RemoveChild(child); child.Cleanup() }
// UPDATE: コンパイルエラーを修正。Nodeが返す[]NodeOwnerを[]Widgetに変換します。
func (sv *ScrollView) GetChildren() []component.Widget {
	nodeOwners := sv.Node.GetChildren()
	widgets := make([]component.Widget, len(nodeOwners))
	for i, owner := range nodeOwners {
		// AddChildでWidgetしか受け付けないため、このアサーションは安全です。
		widgets[i] = owner.(component.Widget)
	}
	return widgets
}

// --- layout.ScrollViewer interface ---
func (sv *ScrollView) GetLayout() layout.Layout { return sv.layout }
func (sv *ScrollView) SetLayout(l layout.Layout) { sv.layout = l }
func (sv *ScrollView) GetPadding() layout.Insets {
	style := sv.ReadOnlyStyle()
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