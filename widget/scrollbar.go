package widget

import (
	"errors"
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ScrollBar is a widget that indicates the state of a scrollable area.
type ScrollBar struct {
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Interaction
	*component.Visibility
	*component.Dirty

	hasBeenLaidOut bool
	contentRatio   float64
	scrollRatio    float64
}

// --- Interface implementation verification ---
var _ component.Widget = (*ScrollBar)(nil)
var _ component.ScrollBarWidget = (*ScrollBar)(nil)
var _ component.NodeOwner = (*ScrollBar)(nil)
var _ component.LayoutPropertiesOwner = (*ScrollBar)(nil)
var _ component.VisibilityOwner = (*ScrollBar)(nil)
var _ component.DirtyManager = (*ScrollBar)(nil)
var _ component.AbsolutePositioner = (*ScrollBar)(nil)
var _ event.EventTarget = (*ScrollBar)(nil)
var _ component.EventProcessor = (*ScrollBar)(nil)

// newScrollBar creates a new component-based ScrollBar.
func newScrollBar() (*ScrollBar, error) {
	s := &ScrollBar{}
	s.Node = component.NewNode(s)
	s.Transform = component.NewTransform()
	s.LayoutProperties = component.NewLayoutProperties()
	s.Appearance = component.NewAppearance(s)
	s.Interaction = component.NewInteraction(s)
	s.Visibility = component.NewVisibility(s)
	s.Dirty = component.NewDirty()

	// Default styles
	s.SetStyle(style.Style{
		Background:  style.PColor(color.RGBA{220, 220, 220, 255}), // Track color
		BorderColor: style.PColor(color.RGBA{180, 180, 180, 255}), // Thumb color
	})

	s.SetSize(10, 100)
	return s, nil
}

// --- Interface implementations ---

func (s *ScrollBar) GetNode() *component.Node                   { return s.Node }
func (s *ScrollBar) GetLayoutProperties() *component.LayoutProperties { return s.LayoutProperties }
func (s *ScrollBar) Update()                                    {}
func (s *ScrollBar) Cleanup()                                   { s.SetParent(nil) }
func (s *ScrollBar) HasBeenLaidOut() bool                       { return s.hasBeenLaidOut }

func (s *ScrollBar) Draw(info component.DrawInfo) {
	if !s.IsVisible() || !s.hasBeenLaidOut {
		return
	}
	x, y := s.GetPosition()
	width, height := s.GetSize()

	finalX := float32(x + info.OffsetX)
	finalY := float32(y + info.OffsetY)

	st := s.ReadOnlyStyle()
	trackColor := color.RGBA{220, 220, 220, 255} // Default color
	if st.Background != nil {
		r, g, b, a := (*st.Background).RGBA()
		trackColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	thumbColor := color.RGBA{180, 180, 180, 255} // Default color
	if st.BorderColor != nil {
		r, g, b, a := (*st.BorderColor).RGBA()
		thumbColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}

	// Draw track
	vector.DrawFilledRect(info.Screen, finalX, finalY, float32(width), float32(height), trackColor, false)

	if s.contentRatio >= 1.0 {
		return
	}

	// Draw thumb
	thumbHeight := float32(float64(height) * s.contentRatio)
	minThumbHeight := float32(10)
	if thumbHeight < minThumbHeight {
		thumbHeight = minThumbHeight
	}
	if height < int(minThumbHeight) {
		return
	}

	thumbYRange := float32(height) - thumbHeight
	thumbY := finalY + thumbYRange*float32(s.scrollRatio)

	vector.DrawFilledRect(info.Screen, finalX, thumbY, float32(width), thumbHeight, thumbColor, false)
}

func (s *ScrollBar) MarkDirty(relayout bool) {
	s.Dirty.MarkDirty(relayout)
	if relayout && !s.IsLayoutBoundary() {
		if parent := s.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (s *ScrollBar) SetPosition(x, y int) {
	if !s.hasBeenLaidOut {
		s.hasBeenLaidOut = true
	}
	if posX, posY := s.GetPosition(); posX != x || posY != y {
		s.Transform.SetPosition(x, y)
		s.MarkDirty(false)
	}
}

func (s *ScrollBar) SetSize(width, height int) {
	if w, h := s.GetSize(); w != width || h != height {
		s.Transform.SetSize(width, height)
		s.MarkDirty(true)
	}
}

func (s *ScrollBar) GetMinSize() (int, int) {
	return 0, 0
}

func (s *ScrollBar) HitTest(x, y int) component.Widget {
	// TODO: Implement hit testing for thumb dragging
	return nil
}

func (s *ScrollBar) HandleEvent(e *event.Event) {
	// TODO: Implement event handling for thumb dragging
}

func (s *ScrollBar) SetRatios(contentRatio, scrollRatio float64) {
	if s.contentRatio != contentRatio || s.scrollRatio != scrollRatio {
		s.contentRatio = contentRatio
		s.scrollRatio = scrollRatio
		s.MarkDirty(false)
	}
}

// --- AbsolutePositioner Implementation ---
func (s *ScrollBar) SetRequestedPosition(x, y int) {
	s.Transform.SetRequestedPosition(x, y)
	s.MarkDirty(true)
}

func (s *ScrollBar) GetRequestedPosition() (int, int) {
	return s.Transform.GetRequestedPosition()
}

// --- ScrollBarBuilder ---
type ScrollBarBuilder struct {
	scrollBar *ScrollBar
	errors    []error
}

func NewScrollBarBuilder() *ScrollBarBuilder {
	s, err := newScrollBar()
	b := &ScrollBarBuilder{scrollBar: s}
	if err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *ScrollBarBuilder) Build() (*ScrollBar, error) {
	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}
	b.scrollBar.MarkDirty(true)
	return b.scrollBar, nil
}

func (b *ScrollBarBuilder) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

func (b *ScrollBarBuilder) TrackColor(c color.Color) *ScrollBarBuilder {
	st := b.scrollBar.GetStyle()
	st.Background = style.PColor(c)
	b.scrollBar.SetStyle(st)
	return b
}

func (b *ScrollBarBuilder) ThumbColor(c color.Color) *ScrollBarBuilder {
	st := b.scrollBar.GetStyle()
	st.BorderColor = style.PColor(c) // Using BorderColor for the thumb
	b.scrollBar.SetStyle(st)
	return b
}
