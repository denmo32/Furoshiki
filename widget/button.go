package widget

import (
	"image"
	"image/color"
	"log"           // [追加] ログ出力のために追加
	"runtime/debug" // [追加] スタックトレース取得のために追加

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
	component.DrawStyledBackground(screen, x, y, width, height, styleToUse)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), styleToUse)
}

// [削除] HitTestメソッドは、component.LayoutableWidgetの汎用的な実装で十分なため、削除します。
// LayoutableWidgetは初期化時に具象ウィジェット(self)への参照を受け取り、
// HitTestが成功した際にその参照を返すため、具象型でのオーバーライドは不要です。

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
	// HoverStyleの呼び出し順問題を解決するため、ホバー用の差分スタイルを保持します。
	hoverStyleDiff *style.Style
}

// NewButtonBuilder は、デフォルトのスタイルで初期化されたButtonBuilderを返します。
// [修正] 初期化をself参照パターンに合わせ、スタイル設定をポインタ対応にします。
func NewButtonBuilder() *ButtonBuilder {
	// ヘルパー関数で値をポインタ化
	ptrFloat32 := func(v float32) *float32 { return &v }
	// [修正] 具象型の値からcolor.Color型の変数を作成し、そのアドレスを渡す
	bgColor := color.Color(color.RGBA{R: 220, G: 220, B: 220, A: 255})
	textColor := color.Color(color.Black)
	borderColor := color.Color(color.Gray{Y: 150})

	defaultStyle := style.Style{
		Background:  &bgColor,
		TextColor:   &textColor,
		BorderColor: &borderColor,
		BorderWidth: ptrFloat32(1),
		Padding: &style.Insets{
			Top: 5, Right: 10, Bottom: 5, Left: 10,
		},
	}
	// まずボタンインスタンスを作成
	button := &Button{}
	// 次に、ボタン自身をselfとして渡してTextWidgetを初期化
	button.TextWidget = component.NewTextWidget(button, "")

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
					// [改善] パニック発生時に、より詳細なデバッグ情報（スタックトレース）をログに出力します。
					log.Printf("Recovered from panic in button click handler: %v\n%s", r, debug.Stack())
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