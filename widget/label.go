package widget

import (
	"furoshiki/component"
	"furoshiki/theme"
)

// Labelはテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*component.TextWidget
}

// newLabelは、ラベルウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//       常にNewLabelBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newLabel(text string) (*Label, error) {
	label := &Label{}
	label.TextWidget = component.NewTextWidget(text)
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
	label, err := newLabel("")
	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	b.AddError(err)
	return b
}

// Build は、最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}