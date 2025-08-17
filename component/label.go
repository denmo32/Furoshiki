package component

import (
	"fmt"
	"furoshiki/style"
	"image/color"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
type Label struct {
	*TextWidget
}

// --- LabelBuilder ---
type LabelBuilder struct {
	label  *Label
	errors []error
}

func NewLabelBuilder() *LabelBuilder {
	label := &Label{
		TextWidget: NewTextWidget(""),
	}
	label.width = 100
	label.height = 30
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

func (b *LabelBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.label.calculateMinSize()
	b.label.SetMinSize(minWidth, minHeight)
}

func (b *LabelBuilder) CalculateMinSize() *LabelBuilder {
	// この呼び出しはBuild時に自動的に行われるため、通常は不要です。
	return b
}

func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.label.SetText(text)
	return b
}

func (b *LabelBuilder) Size(width, height int) *LabelBuilder {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("label size must be non-negative, got %dx%d", width, height))
		return b
	}
	b.label.SetSize(width, height)
	return b
}

func (b *LabelBuilder) Style(s style.Style) *LabelBuilder {
	existingStyle := b.label.GetStyle()
	b.label.SetStyle(style.Merge(*existingStyle, s))
	return b
}

func (b *LabelBuilder) Flex(flex int) *LabelBuilder {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("label flex must be non-negative, got %d", flex))
		return b
	}
	b.label.SetFlex(flex)
	return b
}

func (b *LabelBuilder) Build() (*Label, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("label build errors: %v", b.errors)
	}
	b.calculateMinSizeInternal()
	return b.label, nil
}
