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
	hoverStyle style.Style // ポインタから値型に変更。Nilチェックが不要になり、安全性が向上します。
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
	if b.IsHovered() {
		styleToUse = b.hoverStyle
	}

	// 選択したスタイルで描画
	// [修正] styleToUseは値型なので、ポインタ参照(*)は不要です。
	component.DrawStyledBackground(screen, x, y, width, height, styleToUse)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), styleToUse)
}

// HitTest は、指定された座標がボタンの領域内にあるかを判定します。
// component.LayoutableWidgetの基本的なテストを呼び出し、ヒットした場合は
// LayoutableWidgetではなく、具象型であるButton自身を返します。
// これにより、イベントシステムが正しいウィジェットインスタンスを扱えるようになります。
func (b *Button) HitTest(x, y int) component.Widget {
	// 埋め込まれたLayoutableWidgetのHitTestを呼び出して、基本的な境界チェックを行います
	if b.LayoutableWidget.HitTest(x, y) != nil {
		// ヒットした場合、インターフェースを満たす具象型であるButton自身(*b)を返します
		return b
	}
	return nil
}

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
	// HoverStyleの呼び出し順問題を解決するため、ホバー用の差分スタイルを保持します。
	hoverStyleDiff *style.Style
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
	// 初期状態では、ホバースタイルは基本スタイルと同じ
	button.hoverStyle = defaultStyle

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
// ここでは差分スタイルを保存するだけに留め、Build時にマージすることで、
// Style() と HoverStyle() の呼び出し順に依存しないようにします。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	b.hoverStyleDiff = &s
	return b
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	// Buildが呼び出された時点で、最終的な基本スタイルとホバー差分スタイルをマージします。
	if b.hoverStyleDiff != nil {
		finalBaseStyle := b.Widget.GetStyle()
		b.Widget.hoverStyle = style.Merge(finalBaseStyle, *b.hoverStyleDiff)
	} else {
		// HoverStyleが指定されなかった場合、ホバースタイルは基本スタイルと同じにします。
		b.Widget.hoverStyle = b.Widget.GetStyle()
	}

	return b.Builder.Build("Button")
}