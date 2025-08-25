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

// 【新規追加】NewSpacerは、Spacerウィジェットの新しいインスタンスを生成し、初期化します。
func NewSpacer() *Spacer {
	s := &Spacer{}
	s.LayoutableWidget = component.NewLayoutableWidget()
	s.Init(s) // LayoutableWidgetの初期化
	return s
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
	// 【改善】新しいNewSpacerコンストラクタを呼び出して、初期化ロジックを集約します。
	s := NewSpacer()

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