package component

import (
	"furoshiki/style"
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
// [修正] 第一引数に、このTextWidgetを埋め込む具象ウィジェット(self)への参照を取るように変更します。
func NewTextWidget(self Widget, text string) *TextWidget {
	return &TextWidget{
		LayoutableWidget: NewLayoutableWidget(self),
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
// [修正] 埋め込み先のDrawメソッドを呼び出すのをやめ、このメソッド内で背景とテキストの両方を描画するように変更します。
// これにより、描画ロジックがこのウィジェット内で完結し、意図しない動作を防ぎます。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	if !t.isVisible {
		return
	}
	// ゲッターメソッドを使用してプロパティを取得
	x, y := t.GetPosition()
	width, height := t.GetSize()
	style := t.GetStyle()

	// 最初に背景と境界線を描画
	DrawStyledBackground(screen, x, y, width, height, style)

	// 次にテキストを描画
	DrawAlignedText(screen, t.text, image.Rect(x, y, x+width, y+height), style)
}

// CalculateMinSize は、現在のテキストとスタイルに基づいて最小サイズを計算します。
// [修正] スタイルのPaddingとFontがポインタになったため、nilチェックを追加します。
func (t *TextWidget) CalculateMinSize() (int, int) {
	s := t.GetStyle()
	// [改善] s.Fontがポインタになったため、nilチェックを追加。
	if t.text != "" && s.Font != nil && *s.Font != nil {
		bounds := text.BoundString(*s.Font, t.text)

		// パディングの値を取得（nilの場合はゼロ値として扱う）
		padding := s.Padding
		if s.Padding == nil {
			// [修正] style.Paddingがnilの場合、ゼロ値のstyle.Insetsを作成して計算を続行
			padding = &style.Insets{}
		}

		minWidth := bounds.Dx() + padding.Left + padding.Right
		metrics := (*s.Font).Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

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