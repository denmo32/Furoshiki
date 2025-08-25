package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"image/color"
)

// textWidget は、テキストを持つウィジェット（Button, Labelなど）が満たすインターフェースです。
type textWidget interface {
	component.Widget
	SetText(string)
}

// Builder は、テキストを持つウィジェットのための汎用ビルダーです。
type Builder[T any, W textWidget] struct {
	component.Builder[T, W]
}

// Init はビルダーを初期化します。
func (b *Builder[T, W]) Init(self T, widget W) {
	b.Builder.Init(self, widget)
}

// Text はウィジェットのテキスト内容を設定します。
func (b *Builder[T, W]) Text(text string) T {
	b.Widget.SetText(text)
	return b.Self
}

// TextColor はウィジェットのテキスト色を設定します。
func (b *Builder[T, W]) TextColor(c color.Color) T {
	return b.Style(style.Style{TextColor: style.PColor(c)})
}

// TextAlign はテキストの水平方向の揃え位置を設定します。
func (b *Builder[T, W]) TextAlign(align style.TextAlignType) T {
	return b.Style(style.Style{TextAlign: style.PTextAlignType(align)})
}

// VerticalAlign はテキストの垂直方向の揃え位置を設定します。
func (b *Builder[T, W]) VerticalAlign(align style.VerticalAlignType) T {
	return b.Style(style.Style{VerticalAlign: style.PVerticalAlignType(align)})
}
