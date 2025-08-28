package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
)

// Spacer is a non-drawing widget used to fill space in a FlexLayout.
type Spacer struct {
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance  // NOTE: For Buildable interface
	*component.Interaction // NOTE: For Buildable interface
	*component.Visibility
	*component.Dirty

	// UPDATE: hasBeenLaidOutフィールドはVisibilityコンポーネントに統合されたため削除されました。
	// hasBeenLaidOut bool
}

// --- Interface implementation verification ---
var _ component.Widget = (*Spacer)(nil)

// NOTE: Spacerがcomponent.Builderで利用可能になるために、Buildableインターフェースを満たす必要があります。
var _ component.Buildable = (*Spacer)(nil)
var _ component.NodeOwner = (*Spacer)(nil)
var _ component.LayoutPropertiesOwner = (*Spacer)(nil)
var _ component.VisibilityOwner = (*Spacer)(nil)
var _ component.DirtyManager = (*Spacer)(nil)
var _ component.AbsolutePositioner = (*Spacer)(nil)
var _ component.EventProcessor = (*Spacer)(nil)

// newSpacer creates a new component-based Spacer.
func newSpacer() (*Spacer, error) {
	s := &Spacer{}
	s.Node = component.NewNode(s)
	s.Transform = component.NewTransform()
	s.LayoutProperties = component.NewLayoutProperties()
	// NOTE: Spacerはスタイルやインタラクションを持ちませんが、
	//       Buildableインターフェースを満たすためにコンポーネントを初期化します。
	s.Appearance = component.NewAppearance(s)
	s.Interaction = component.NewInteraction(s)
	s.Visibility = component.NewVisibility(s)
	s.Dirty = component.NewDirty()
	return s, nil
}

// --- Interface implementations ---

func (s *Spacer) GetNode() *component.Node                         { return s.Node }
func (s *Spacer) GetLayoutProperties() *component.LayoutProperties { return s.LayoutProperties }
func (s *Spacer) Update()                                          {}
func (s *Spacer) Cleanup()                                         { s.SetParent(nil) }
func (s *Spacer) Draw(info component.DrawInfo)                     {} // Spacer is not drawn

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
	// UPDATE: レイアウト済み状態の管理をVisibilityコンポーネントに委譲します。
	if !s.HasBeenLaidOut() {
		s.SetLaidOut(true)
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

func (s *Spacer) SetMinSize(width, height int) {} // Spacer has no min size
func (s *Spacer) GetMinSize() (int, int) {
	return 0, 0
}

func (s *Spacer) HitTest(x, y int) component.Widget {
	return nil // Spacer is not interactive
}

// HandleEvent is a dummy implementation to satisfy the EventProcessor interface.
func (s *Spacer) HandleEvent(e *event.Event) {}

// --- AbsolutePositioner and other Buildable interface implementations ---
func (s *Spacer) SetRequestedPosition(x, y int) {
	s.Transform.SetRequestedPosition(x, y)
	s.MarkDirty(true)
}

func (s *Spacer) GetRequestedPosition() (int, int) {
	return s.Transform.GetRequestedPosition()
}

// NOTE: 以下のメソッドは component.Buildable インターフェースを満たすために実装されています。
func (s *Spacer) SetFlex(flex int)                  { s.LayoutProperties.SetFlex(flex) }
func (s *Spacer) GetFlex() int                      { return s.LayoutProperties.GetFlex() }
func (s *Spacer) SetLayoutBoundary(isBoundary bool) { s.LayoutProperties.SetLayoutBoundary(isBoundary) }
func (s *Spacer) SetLayoutData(data any)            { s.LayoutProperties.SetLayoutData(data) }
func (s *Spacer) GetLayoutData() any                { return s.LayoutProperties.GetLayoutData() }

// NOTE: Spacerはスタイルを持たないが、インターフェースを満たすためにダミーメソッドを実装
func (s *Spacer) SetStyle(style style.Style) {}
func (s *Spacer) GetStyle() style.Style      { return style.Style{} }
func (s *Spacer) ReadOnlyStyle() style.Style { return style.Style{} }

// --- SpacerBuilder ---

// 【提案3対応】SpacerBuilderは汎用のcomponent.Builderを埋め込むように変更されました。
type SpacerBuilder struct {
	component.Builder[*SpacerBuilder, *Spacer]
}

// NewSpacerBuilderは新しいSpacerBuilderを生成します。
func NewSpacerBuilder() *SpacerBuilder {
	spacer, err := newSpacer()
	b := &SpacerBuilder{}
	b.Init(b, spacer)
	b.AddError(err)
	return b
}

// Build は、最終的なSpacerを構築して返します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}

// NOTE: Size, Flex, AssignTo などの共通メソッドは component.Builder に
//       実装されているため、ここからは削除されました。