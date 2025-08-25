package widget

import (
	"furoshiki/component"
	"github.com/hajimehoshi/ebiten/v2"
)

// SpacerはFlexLayout内で余白を埋めるために使用される、描画されないウィジェットです。
type Spacer struct {
	*component.LayoutableWidget
}

// NewSpacerは、Spacerウィジェットの新しいインスタンスを生成し、初期化します。
func NewSpacer() *Spacer {
	s := &Spacer{}
	s.LayoutableWidget = component.NewLayoutableWidget()
	s.Init(s)
	return s
}

// Draw は何もしません。Spacerは視覚的な表現を持たないためです。
func (s *Spacer) Draw(screen *ebiten.Image) {}

// --- SpacerBuilder ---
type SpacerBuilder struct {
	component.Builder[*SpacerBuilder, *Spacer]
}

// NewSpacerBuilder は新しいSpacerBuilderを作成します。
func NewSpacerBuilder() *SpacerBuilder {
	s := NewSpacer()
	b := &SpacerBuilder{}
	b.Init(b, s)
	return b
}

// Build は最終的なSpacerウィジェットを返します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}
