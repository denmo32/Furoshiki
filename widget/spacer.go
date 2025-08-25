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

// Draw は何もしません。Spacerは視覚的な表現を持たないため、
// LayoutableWidgetの背景描画処理をオーバーライドして無効化します。
func (s *Spacer) Draw(screen *ebiten.Image) {
	// Spacer is invisible and should not draw anything.
}

// --- SpacerBuilder ---

// SpacerBuilder はSpacerを構築するためのビルダーです。
// component.Builderを埋め込むことで、Flexなどの共通プロパティ設定をサポートします。
type SpacerBuilder struct {
	component.Builder[*SpacerBuilder, *Spacer]
}

// NewSpacerBuilder は新しいSpacerBuilderを作成します。
func NewSpacerBuilder() *SpacerBuilder {
	// Spacerインスタンスを作成します。
	s := &Spacer{}
	// 【改善】LayoutableWidgetを初期化し、Initメソッドでself参照を設定します。
	s.LayoutableWidget = component.NewLayoutableWidget()
	s.Init(s)

	b := &SpacerBuilder{}
	// 汎用ビルダーを初期化し、メソッドチェーンを可能にします。
	b.Init(b, s)
	return b
}

// Build は最終的なSpacerウィジェットを返します。
// 内部で汎用のcomponent.BuilderのBuildメソッドを呼び出します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}