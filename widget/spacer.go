package widget

import (
	"furoshiki/component"
	"github.com/hajimehoshi/ebiten/v2"
)

// SpacerはFlexLayout内で余白を埋めるために使用される、描画されないウィジェットです。
type Spacer struct {
	*component.LayoutableWidget
}

// newSpacerは、Spacerウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//       常にNewSpacerBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newSpacer() (*Spacer, error) {
	s := &Spacer{}
	s.LayoutableWidget = component.NewLayoutableWidget()
	if err := s.Init(s); err != nil {
		return nil, err
	}
	return s, nil
}

// Draw は何もしません。Spacerは視覚的な表現を持たないためです。
func (s *Spacer) Draw(screen *ebiten.Image) {}

// --- SpacerBuilder ---
type SpacerBuilder struct {
	component.Builder[*SpacerBuilder, *Spacer]
}

// NewSpacerBuilder は新しいSpacerBuilderを作成します。
func NewSpacerBuilder() *SpacerBuilder {
	s, err := newSpacer()
	b := &SpacerBuilder{}
	b.Init(b, s)
	b.AddError(err)
	return b
}

// Build は最終的なSpacerウィジェットを返します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}