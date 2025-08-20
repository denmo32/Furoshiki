package widget

import (
	"furoshiki/component"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpacerはFlexLayout内で余白を埋めるために使用される、描画されないウィジェットです。
// 主にFlexプロパティと組み合わせて、要素間のスペースを動的に調整するために使います。
type Spacer struct {
	*component.LayoutableWidget
}

// NewSpacer は新しいSpacerウィジェットを作成します。
func NewSpacer() *Spacer {
	s := &Spacer{}
	// Spacer自身をselfとして渡してLayoutableWidgetを初期化します。
	s.LayoutableWidget = component.NewLayoutableWidget(s)
	return s
}

// Draw は何もしません。Spacerは視覚的な表現を持たないため、
// LayoutableWidgetの背景描画処理をオーバーライドして無効化します。
func (s *Spacer) Draw(screen *ebiten.Image) {
	// Spacer is invisible and should not draw anything.
}

// --- SpacerBuilder ---

// SpacerBuilder はSpacerを構築するためのビルダーです。
type SpacerBuilder struct {
	spacer *Spacer
	errors []error
}

// NewSpacerBuilder は新しいSpacerBuilderを作成します。
func NewSpacerBuilder() *SpacerBuilder {
	return &SpacerBuilder{
		spacer: NewSpacer(),
	}
}

// Flex はSpacerの伸縮係数を設定します。
func (b *SpacerBuilder) Flex(flex int) *SpacerBuilder {
	b.spacer.SetFlex(flex)
	return b
}

// Build は最終的なSpacerウィジェットを返します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	// Spacerには複雑な検証はないため、エラーチェックはシンプルです。
	if len(b.errors) > 0 {
		// 現状エラーを追加するパスはないが、将来のために残す
		return nil, b.errors[0]
	}
	return b.spacer, nil
}