package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
// Button は、クリック可能なUI要素です。TextWidgetを拡張し、ホバー状態のスタイル管理機能を追加します。
type Button struct {
	*component.TextWidget
	hoverStyle style.Style // マウスホバー時のスタイル。値型で保持することでnilチェックを不要にします。
}

// SetHovered はホバー状態を設定し、再描画を要求します。
// component.LayoutableWidgetのSetHoveredを直接呼び出すことで、ダーティフラグが適切に設定されます。
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

	// 描画時にホバー状態を確認し、適切なスタイルを選択します。
	styleToUse := b.GetStyle()
	if b.IsHovered() {
		styleToUse = b.hoverStyle
	}

	// 選択したスタイルで背景とテキストを描画します。
	component.DrawStyledBackground(screen, x, y, width, height, styleToUse)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), styleToUse)
}

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
	// HoverStyleの呼び出し順問題を解決するため、ホバー用の差分スタイルを保持します。
	hoverStyleDiff *style.Style
}

// NewButtonBuilder は新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	defaultStyle := style.Style{
		Background:  style.PColor(color.RGBA{R: 220, G: 220, B: 220, A: 255}),
		TextColor:   style.PColor(color.Black),
		BorderColor: style.PColor(color.Gray{Y: 150}),
		BorderWidth: style.PFloat32(1),
		Padding: style.PInsets(style.Insets{
			Top: 5, Right: 10, Bottom: 5, Left: 10,
		}),
		TextAlign:     style.PTextAlignType(style.TextAlignCenter),
		VerticalAlign: style.PVerticalAlignType(style.VerticalAlignMiddle),
	}
	// まずボタンインスタンスを作成
	button := &Button{}
	// 次に、ボタン自身をselfとして渡してTextWidgetを初期化
	button.TextWidget = component.NewTextWidget(button, "")

	button.SetSize(100, 40)
	button.SetStyle(defaultStyle)
	// hoverStyleの最終的な設定は、全てのスタイル設定が終わった後のBuild()メソッドで行います。
	// これにより、ロジックがBuild()内に集約され、より見通しが良くなります。

	b := &ButtonBuilder{}
	b.Builder.Init(b, button)
	return b
}

// OnClick は、ボタンがクリックされたときに実行されるイベントハンドラを設定します。
// ハンドラは event.Event を引数として受け取るため、クリック座標などの詳細情報にアクセスできます。
// イベント情報が不要な場合は、引数を無視して `func(_ event.Event) { ... }` のように記述できます。
func (b *ButtonBuilder) OnClick(handler event.EventHandler) *ButtonBuilder {
	if handler != nil {
		// イベントハンドラ内のパニック回復処理は、component.LayoutableWidgetのHandleEventメソッドが
		// 一括して担うため、ここでの二重の回復処理は不要です。
		b.Widget.AddEventHandler(event.EventClick, handler)
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

	// 汎用ビルダーのBuildを呼び出して、最終的なエラーチェックなどを行います。
	return b.Builder.Build("Button")
}