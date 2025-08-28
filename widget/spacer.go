package widget

import (
	"furoshiki/component"
)

// Spacer is a non-drawing widget used to fill space in a FlexLayout.
type Spacer struct {
	// UPDATE: Spacerが必要とする最小限の共通コンポーネントをWidgetCoreから取得
	*component.WidgetCore
}

// --- Interface implementation verification ---
// UPDATE: Spacerはスタイルやインタラクションを持たないため、StandardWidgetは実装しません。
//         Widgetインターフェースと、レイアウトに必要な基本的なインターフェースのみを実装します。
var _ component.Widget = (*Spacer)(nil)
var _ component.NodeOwner = (*Spacer)(nil)
var _ component.LayoutPropertiesOwner = (*Spacer)(nil)
var _ component.VisibilityOwner = (*Spacer)(nil)
var _ component.DirtyManager = (*Spacer)(nil)
var _ component.AbsolutePositioner = (*Spacer)(nil)

// newSpacer creates a new component-based Spacer.
func newSpacer() (*Spacer, error) {
	s := &Spacer{}
	// UPDATE: WidgetCoreのみを初期化
	s.WidgetCore = component.NewWidgetCore(s)
	return s, nil
}

// --- Interface implementations ---

// GetNode, GetLayoutProperties are now inherited from WidgetCore.
func (s *Spacer) Update()                {}
func (s *Spacer) Cleanup()               { s.SetParent(nil) }
func (s *Spacer) Draw(info component.DrawInfo) {} // Spacer is not drawn

// UPDATE: MarkDirty, SetPosition, SetSize, SetMinSize, GetMinSize, SetRequestedPosition,
// GetRequestedPositionはすべてWidgetCoreに実装されているため、このファイルからは削除されました。
// Spacerはコンテンツを持たないため、GetMinSizeはWidgetCoreのデフォルト実装(0,0)で十分です。

func (s *Spacer) HitTest(x, y int) component.Widget {
	return nil // Spacer is not interactive
}

// --- SpacerBuilder ---

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

func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}