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
// NOTE: 内部のInit呼び出しが失敗する可能性があるため、コンストラクタはerrorを返すように変更されました。
func NewSpacer() (*Spacer, error) {
	s := &Spacer{}
	s.LayoutableWidget = component.NewLayoutableWidget()
	// NOTE: Initがエラーを返すようになったため、エラーをチェックし、呼び出し元に伝播させます。
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
	s, err := NewSpacer()
	b := &SpacerBuilder{}
	b.Init(b, s)
	// NOTE: コンストラクタで発生した初期化エラーをビルダーに追加します。
	b.AddError(err)
	return b
}

// Build は最終的なSpacerウィジェットを返します。
func (b *SpacerBuilder) Build() (*Spacer, error) {
	return b.Builder.Build()
}
