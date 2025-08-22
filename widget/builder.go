package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"image/color"
)

// textWidget is an interface that text-based widgets like Button and Label satisfy.
type textWidget interface {
	component.Widget
	SetText(string)
}

// Builder is a generic builder for text-based widgets.
// It embeds the component.Builder and adds text-specific methods.
type Builder[T any, W textWidget] struct {
	component.Builder[T, W]
}

// Init initializes the builder.
func (b *Builder[T, W]) Init(self T, widget W) {
	b.Builder.Init(self, widget)
}

// Text sets the widget's text content.
func (b *Builder[T, W]) Text(text string) T {
	b.Widget.SetText(text)
	return b.Self
}

// TextColor sets the widget's text color.
func (b *Builder[T, W]) TextColor(c color.Color) T {
	return b.Style(style.Style{TextColor: style.PColor(c)})
}

// TextAlign sets the horizontal alignment of the text.
func (b *Builder[T, W]) TextAlign(align style.TextAlignType) T {
	return b.Style(style.Style{TextAlign: style.PTextAlignType(align)})
}

// VerticalAlign sets the vertical alignment of the text.
func (b *Builder[T, W]) VerticalAlign(align style.VerticalAlignType) T {
	return b.Style(style.Style{VerticalAlign: style.PVerticalAlignType(align)})
}