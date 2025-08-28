package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"image"
	"log"
	"runtime/debug"
)

// ScrollView is a container widget that allows scrolling its content.
// It has been refactored to use a component-based architecture.
type ScrollView struct {
	// Component-based fields
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Interaction
	*component.Visibility
	*component.Dirty

	// ScrollView specific fields
	container         *container.Container
	layout            layout.Layout
	contentContainer  component.Widget
	vScrollBar        component.ScrollBarWidget
	scrollY           float64
	contentHeight     int
	ScrollSensitivity float64
	hasBeenLaidOut    bool
}

// --- Interface implementation verification ---
var _ component.Widget = (*ScrollView)(nil)
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)
var _ container.Scroller = (*ScrollView)(nil)
var _ event.EventTarget = (*ScrollView)(nil)
var _ component.EventProcessor = (*ScrollView)(nil)

// newScrollView creates a new component-based ScrollView.
func newScrollView() (*ScrollView, error) {
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}

	// Initialize components
	sv.Node = component.NewNode(sv)
	sv.Transform = component.NewTransform()
	sv.LayoutProperties = component.NewLayoutProperties()
	sv.Appearance = component.NewAppearance(sv)
	sv.Interaction = component.NewInteraction(sv)
	sv.Visibility = component.NewVisibility(sv)
	sv.Dirty = component.NewDirty()

	internalContainer, err := container.NewContainer()
	if err != nil {
		return nil, err
	}
	sv.container = internalContainer
	sv.container.SetClipsChildren(true)
	sv.container.SetParent(sv)

	sv.layout = &layout.ScrollViewLayout{}

	vScrollBar, err := NewScrollBarBuilder().Build()
	if err != nil {
		return nil, err
	}
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar)

	sv.AddEventHandler(event.MouseScroll, sv.onMouseScroll)

	return sv, nil
}

// onMouseScroll handles the mouse scroll event to scroll the content.
func (sv *ScrollView) onMouseScroll(e *event.Event) event.Propagation {
	scrollAmount := e.ScrollY * sv.ScrollSensitivity
	sv.scrollY -= scrollAmount
	sv.MarkDirty(true)
	return event.StopPropagation
}

// SetContent sets the scrollable content container.
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		sv.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content
	if content != nil {
		sv.AddChild(content)
		sv.AddChild(sv.vScrollBar) // Re-add scrollbar to keep it on top
	}
	sv.MarkDirty(true)
}

// SetStyle sets the style for both the ScrollView and its internal container.
func (sv *ScrollView) SetStyle(s style.Style) {
	sv.Appearance.SetStyle(s)
	if sv.container != nil {
		sv.container.SetStyle(s)
	}
}

// Update controls the specialized layout calculation for the ScrollView.
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

	for _, child := range sv.container.GetChildren() {
		child.Update()
	}

	if sv.IsDirty() {
		sv.ClearDirty()
	}
}

// Draw delegates the drawing to the internal container which handles clipping.
func (sv *ScrollView) Draw(info component.DrawInfo) {
	if !sv.IsVisible() || !sv.hasBeenLaidOut {
		return
	}
	sv.container.Draw(info)
}

// Cleanup releases resources used by the ScrollView and its children.
func (sv *ScrollView) Cleanup() {
	sv.container.Cleanup()
	sv.SetParent(nil)
}

// MarkDirty marks the widget as needing a redraw or relayout.
// ScrollView is a layout boundary, so it doesn't propagate the dirty flag up.
func (sv *ScrollView) MarkDirty(relayout bool) {
	sv.Dirty.MarkDirty(relayout)
}

// SetPosition sets the position for the ScrollView and its internal container.
func (sv *ScrollView) SetPosition(x, y int) {
	if !sv.hasBeenLaidOut {
		sv.hasBeenLaidOut = true
	}
	if currX, currY := sv.GetPosition(); currX != x || currY != y {
		sv.Transform.SetPosition(x, y)
		if sv.container != nil {
			sv.container.SetPosition(x, y)
		}
		sv.MarkDirty(false)
	}
}

// SetSize sets the size for the ScrollView and its internal container.
func (sv *ScrollView) SetSize(width, height int) {
	if w, h := sv.GetSize(); w != width || h != height {
		sv.Transform.SetSize(width, height)
		if sv.container != nil {
			sv.container.SetSize(width, height)
		}
		sv.MarkDirty(true)
	}
}

// SetMinSize sets the minimum size for the ScrollView.
// This is a no-op for ScrollView as its size is determined by its container.
func (sv *ScrollView) SetMinSize(width, height int) {
	// No-op
}

// GetMinSize returns the minimum size for the ScrollView.
func (sv *ScrollView) GetMinSize() (int, int) {
	// ScrollView itself doesn't have an intrinsic minimum size.
	// It's determined by its container and layout.
	return 0, 0
}

// HitTest checks for hits first on the ScrollView itself, then on its internal container.
func (sv *ScrollView) HitTest(x, y int) component.Widget {
	if !sv.IsVisible() || sv.IsDisabled() {
		return nil
	}
	wx, wy := sv.GetPosition()
	wwidth, wheight := sv.GetSize()
	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if !rect.Empty() && (image.Point{X: x, Y: y}.In(rect)) {
		if target := sv.container.HitTest(x, y); target != nil {
			return target
		}
		return sv
	}
	return nil
}

// HandleEvent processes events, similar to the Button's implementation.
func (sv *ScrollView) HandleEvent(e *event.Event) {
	if handlers, exists := sv.GetEventHandlers()[e.Type]; exists {
		for _, handler := range handlers {
			if e.Handled {
				break
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf(`Recovered from panic in event handler: %v
%s`, r, debug.Stack())
					}
				}()
				if handler(e) == event.StopPropagation {
					e.Handled = true
				}
			}()
		}
	}

	if e != nil && !e.Handled && sv.GetParent() != nil {
		if processor, ok := sv.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// --- Method delegation to internal container ---
func (sv *ScrollView) AddChild(child component.Widget)    { sv.container.AddChild(child) }
func (sv *ScrollView) RemoveChild(child component.Widget) { sv.container.RemoveChild(child) }
func (sv *ScrollView) GetChildren() []component.Widget    { return sv.container.GetChildren() }

// --- layout.ScrollViewer interface ---
func (sv *ScrollView) GetLayout() layout.Layout           { return sv.layout }
func (sv *ScrollView) SetLayout(l layout.Layout)          { sv.layout = l }
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