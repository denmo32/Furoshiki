package widget

import (
	"furoshiki/component"
	"furoshiki/theme"
)

// Labelはテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*component.TextWidget
}

// NewLabelは、ラベルウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: 内部のInit呼び出しが失敗する可能性があるため、コンストラクタはerrorを返すように変更されました。
func NewLabel(text string) (*Label, error) {
	label := &Label{}
	label.TextWidget = component.NewTextWidget(text)
	// NOTE: Initがエラーを返すようになったため、エラーをチェックし、呼び出し元に伝播させます。
	if err := label.Init(label); err != nil {
		return nil, err
	}

	t := theme.GetCurrent()
	label.SetStyle(t.Label.Default)
	label.SetSize(100, 30)

	return label, nil
}

// --- LabelBuilder ---
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	label, err := NewLabel("")
	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	// NOTE: コンストラクタで発生した初期化エラーをビルダーに追加します。
	b.AddError(err)
	return b
}

// Build は、最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}