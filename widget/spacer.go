package widget

import (
	"errors"
	"furoshiki/component"
)

// Spacer is a non-drawing widget used to fill space in a FlexLayout.
type Spacer struct {
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Visibility
	*component.Dirty

	hasBeenLaidOut bool
}

// --- Interface implementation verification ---
var _ component.Widget = (*Spacer)(nil)
var _ component.NodeOwner = (*Spacer)(nil)
var _ component.LayoutPropertiesOwner = (*Spacer)(nil)
var _ component.VisibilityOwner = (*Spacer)(nil)
var _ component.DirtyManager = (*Spacer)(nil)
var _ component.AbsolutePositioner = (*Spacer)(nil)

// newSpacer creates a new component-based Spacer.
func newSpacer() (*Spacer, error) {
	s := &Spacer{}
	s.Node = component.NewNode(s)
	s.Transform = component.NewTransform()
	s.LayoutProperties = component.NewLayoutProperties()
	s.Visibility = component.NewVisibility(s)
	s.Dirty = component.NewDirty()
	return s, nil
}

// --- Interface implementations ---

func (s *Spacer) GetNode() *component.Node                   { return s.Node }
func (s *Spacer) GetLayoutProperties() *component.LayoutProperties { return s.LayoutProperties }
func (s *Spacer) Update()                                    {}
func (s *Spacer) Cleanup()                                   { s.SetParent(nil) }
func (s *Spacer) Draw(info component.DrawInfo)               {} // Spacer is not drawn

func (s *Spacer) MarkDirty(relayout bool) {
	s.Dirty.MarkDirty(relayout)
	if relayout && !s.IsLayoutBoundary() {
		if parent := s.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (s *Spacer) SetPosition(x, y int) {
	if !s.hasBeenLaidOut {
		s.hasBeenLaidOut = true
	}
	if posX, posY := s.GetPosition(); posX != x || posY != y {
		s.Transform.SetPosition(x, y)
		s.MarkDirty(false)
	}
}

func (s *Spacer) SetSize(width, height int) {
	if w, h := s.GetSize(); w != width || h != height {
		s.Transform.SetSize(width, height)
		s.MarkDirty(true)
	}
}

// Spacer has no intrinsic minimum size.
func (s *Spacer) GetMinSize() (int, int) {
	return 0, 0
}

func (s *Spacer) HitTest(x, y int) component.Widget {
	return nil // Spacer is not interactive
}

// --- AbsolutePositioner Implementation ---
func (s *Spacer) SetRequestedPosition(x, y int) {
	s.Transform.SetRequestedPosition(x, y)
	s.MarkDirty(true)
}

func (s *Spacer) GetRequestedPosition() (int, int) {
	return s.Transform.GetRequestedPosition()
}

// --- SpacerBuilder ---
type SpacerBuilder struct {
	spacer *Spacer
	errors []error
}

func NewSpacerBuilder() *SpacerBuilder {
	spacer, err := newSpacer()
	b := &SpacerBuilder{spacer: spacer}
	if err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *SpacerBuilder) Build() (*Spacer, error) {
	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}
	b.spacer.MarkDirty(true)
	return b.spacer, nil
}

func (b *SpacerBuilder) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// --- Builder Methods ---

func (b *SpacerBuilder) Size(width, height int) *SpacerBuilder {
	b.spacer.SetSize(width, height)
	return b
}

func (b *SpacerBuilder) Flex(flex int) *SpacerBuilder {
	b.spacer.SetFlex(flex)
	return b
}

func (b *SpacerBuilder) AssignTo(target **Spacer) *SpacerBuilder {
	if target == nil {
		b.errors = append(b.errors, errors.New("AssignTo target cannot be nil"))
		return b
	}
	*target = b.spacer
	return b
}
