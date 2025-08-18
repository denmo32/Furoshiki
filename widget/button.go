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
	hoverStyle *style.Style
}

// SetHovered はホバー状態を設定し、再描画を要求します。
// 実際のスタイルの選択はDrawメソッドで行われます。
func (b *Button) SetHovered(hovered bool) {
	b.TextWidget.SetHovered(hovered)
}

// Draw はButtonを描画します。ホバー状態に応じて適切なスタイルを適用します。
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	x, y := b.GetPosition()
	width, height := b.GetSize()
	text := b.Text()

	// 描画時にホバー状態を確認し、適切なスタイルを選択する
	styleToUse := b.GetStyle()
	if b.IsHovered() && b.hoverStyle != nil {
		styleToUse = b.hoverStyle
	}

	// 選択したスタイルで描画
	component.DrawStyledBackground(screen, x, y, width, height, *styleToUse)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), *styleToUse)
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
// NOTE: このメソッドはベースビルダーのStyleを呼び出します。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	b.Builder.Style(s)
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるときのスタイルを設定します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	// ベーススタイルにホバースタイルをマージして、完全なホバースタイルを生成
	mergedHoverStyle := style.Merge(*b.Widget.GetStyle(), s)
	b.Widget.hoverStyle = &mergedHoverStyle
	return b
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build("Button")
}