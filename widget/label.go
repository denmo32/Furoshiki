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
func NewLabel(text string) *Label {
	label := &Label{}
	label.TextWidget = component.NewTextWidget(text)
	label.Init(label) // LayoutableWidgetの初期化

	t := theme.GetCurrent()
	label.SetStyle(t.Label.Default)
	label.SetSize(100, 30)

	return label
}

// --- LabelBuilder ---
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	label := NewLabel("")
	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	return b
}

// Build は、最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}
