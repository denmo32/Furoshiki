package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"image/color"
)

// textWidget は、テキストを持つウィジェット（Button, Labelなど）が満たすインターフェースです。
// 【提案1対応】component.Builderが要求するcomponent.Buildableインターフェースを埋め込むことで、
// 型制約を満たすように修正しました。
type textWidget interface {
	component.Buildable
	SetText(string)
	SetWrapText(bool) // 折り返し設定メソッドを追加
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

// WrapText は、ウィジェットの幅を超えるテキストを自動的に折り返すかどうかを設定します。
func (b *Builder[T, W]) WrapText(wrap bool) T {
	b.Widget.SetWrapText(wrap)
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