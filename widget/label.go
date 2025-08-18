package widget

import (
	"fmt"
	"image/color"

	"furoshiki/core"
	"furoshiki/style"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
// 主にテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*core.TextWidget
}

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	label  *Label
	errors []error
}

// NewLabelBuilder は、デフォルトのスタイルで初期化されたLabelBuilderを返します。
func NewLabelBuilder() *LabelBuilder {
	label := &Label{
		TextWidget: core.NewTextWidget(""),
	}
	label.SetSize(100, 30)
	defaultStyle := style.Style{
		Background: color.Transparent,
		TextColor:  color.Black,
		Padding:    style.Insets{Top: 2, Right: 5, Bottom: 2, Left: 5},
	}
	label.SetStyle(defaultStyle)

	return &LabelBuilder{
		label: label,
	}
}

// calculateMinSizeInternal は、ラベルのテキストとパディングに基づいて最小サイズを計算し、設定します。
func (b *LabelBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.label.CalculateMinSize()
	b.label.SetMinSize(minWidth, minHeight)
}

// CalculateMinSize は、ラベルの最小サイズを計算します。
// この呼び出しはBuild時に自動的に行われるため、通常はユーザーが呼び出す必要はありません。
func (b *LabelBuilder) CalculateMinSize() *LabelBuilder {
	return b
}

// Text はラベルに表示されるテキストを設定します。
func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.label.SetText(text)
	return b
}

// Size はラベルのサイズを設定します。
func (b *LabelBuilder) Size(width, height int) *LabelBuilder {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("label size must be non-negative, got %dx%d", width, height))
		return b
	}
	b.label.SetSize(width, height)
	return b
}

// Style はラベルのスタイルを設定します。
func (b *LabelBuilder) Style(s style.Style) *LabelBuilder {
	existingStyle := b.label.GetStyle()
	b.label.SetStyle(style.Merge(*existingStyle, s))
	return b
}

// Flex は、親がFlexLayoutの場合にラベルがどのように伸縮するかを設定します。
func (b *LabelBuilder) Flex(flex int) *LabelBuilder {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("label flex must be non-negative, got %d", flex))
		return b
	}
	b.label.SetFlex(flex)
	return b
}

// Build は、設定に基づいて最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("label build errors: %v", b.errors)
	}
	b.calculateMinSizeInternal()
	return b.label, nil
}
