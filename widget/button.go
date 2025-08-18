package widget

import (
	"fmt"
	"image"
	"image/color"

	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
// Button は、クリック可能なUI要素です。TextWidgetを拡張し、ホバー状態のスタイル管理機能を追加します。
type Button struct {
	*component.TextWidget
	hoverStyle   *style.Style
	currentStyle *style.Style // 描画時に使用するスタイルをキャッシュ
}

// SetHovered はホバー状態を設定し、描画スタイルを更新します。
func (b *Button) SetHovered(hovered bool) {
	b.TextWidget.SetHovered(hovered)
	if b.IsHovered() && b.hoverStyle != nil {
		b.currentStyle = b.hoverStyle
	} else {
		b.currentStyle = b.GetStyle()
	}
}

// Draw はButtonを描画します。キャッシュされたスタイルを使用します。
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	x, y := b.GetPosition()
	width, height := b.GetSize()
	text := b.Text()
	// キャッシュされたスタイルで描画
	component.DrawStyledBackground(screen, x, y, width, height, *b.currentStyle)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), *b.currentStyle)
}

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
}

// NewButtonBuilder は、デフォルトのスタイルで初期化されたButtonBuilderを返します。
func NewButtonBuilder() *ButtonBuilder {
	defaultStyle := style.Style{
		Background:  color.RGBA{R: 220, G: 220, B: 220, A: 255},
		TextColor:   color.Black,
		BorderColor: color.Gray{Y: 150},
		BorderWidth: 1,
		Padding:     style.Insets{Top: 5, Right: 10, Bottom: 5, Left: 10},
	}
	button := &Button{
		TextWidget: component.NewTextWidget(""),
	}
	button.SetSize(100, 40)
	button.SetStyle(defaultStyle)
	// 初期状態のスタイルをキャッシュ
	button.currentStyle = button.GetStyle()

	b := &ButtonBuilder{}
	b.Builder.Init(b, button)
	return b
}

// OnClick は、ボタンがクリックされたときに実行されるイベントハンドラを設定します。
func (b *ButtonBuilder) OnClick(onClick func()) *ButtonBuilder {
	if onClick != nil {
		b.Widget.AddEventHandler(event.EventClick, func(e event.Event) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in button click handler: %v\n", r)
				}
			}()
			onClick()
		})
	}
	return b
}

// Style はボタンの基本スタイルを設定します。
// NOTE: このメソッドはベースビルダーのStyleをオーバーライドして、currentStyleも更新します。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	b.Builder.Style(s)
	b.Widget.currentStyle = b.Widget.GetStyle()
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるときのスタイルを設定します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	mergedHoverStyle := style.Merge(*b.Widget.GetStyle(), s)
	b.Widget.hoverStyle = &mergedHoverStyle
	return b
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build("Button")
}