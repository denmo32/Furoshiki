package component

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// --- TextWidget ---
// TextWidgetは、テキスト表示に関連する共通の機能（テキスト内容、スタイル、最小サイズ計算）を提供します。
// ButtonやLabelなど、テキストを持つウィジェットはこれを埋め込みます。
type TextWidget struct {
	*LayoutableWidget
	text string
}

// NewTextWidget は新しいTextWidgetを生成します。
func NewTextWidget(text string) *TextWidget {
	return &TextWidget{
		LayoutableWidget: NewLayoutableWidget(),
		text:             text,
	}
}

// Text はウィジェットのテキストを取得します。
func (t *TextWidget) Text() string {
	return t.text
}

// SetText はウィジェットのテキストを設定し、ダーティフラグを立てます。
func (t *TextWidget) SetText(text string) {
	if t.text != text {
		t.text = text
		// テキスト変更は最小サイズに影響する可能性があるため再レイアウトが必要
		t.MarkDirty(true)
	}
}

// Draw はTextWidgetを描画します。LayoutableWidgetのDrawをオーバーライドしてテキストを追加描画します。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	if !t.isVisible {
		return
	}
	// 背景描画は基本のDrawを呼び出す
	t.LayoutableWidget.Draw(screen)

	// テキストを描画
	// プロパティへの直接アクセスではなく、ゲッターメソッドを使用して一貫性を保つ
	x, y := t.GetPosition()
	width, height := t.GetSize()
	// [改善] GetStyle()が値型を返すようになったため、戻り値を変数に受けて使用します。
	style := t.GetStyle()
	DrawAlignedText(screen, t.text, image.Rect(x, y, x+width, y+height), style)
}

// CalculateMinSize は、現在のテキストとスタイルに基づいて最小サイズを計算します。
func (t *TextWidget) CalculateMinSize() (int, int) {
	style := t.GetStyle()
	if t.text != "" && style.Font != nil {
		bounds := text.BoundString(style.Font, t.text)
		minWidth := bounds.Dx() + style.Padding.Left + style.Padding.Right
		metrics := style.Font.Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + style.Padding.Top + style.Padding.Bottom

		// 既存の最小サイズより大きい場合はそれを優先
		if t.minWidth > minWidth {
			minWidth = t.minWidth
		}
		if t.minHeight > minHeight {
			minHeight = t.minHeight
		}

		return minWidth, minHeight
	}
	// テキストがない場合でも設定済みの最小サイズを返す
	return t.minWidth, t.minHeight
}