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
// 第一引数には、このTextWidgetを埋め込む具象ウィジェット(self)への参照を取ります。
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
		// テキスト変更は最小サイズに影響し、レイアウトが変わる可能性があるため再レイアウトを要求します。
		t.MarkDirty(true)
	}
}

// Draw はTextWidgetを描画します。
// このメソッド内で背景とテキストの両方を描画することで、描画ロジックがこのウィジェット内で完結します。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	if !t.isVisible {
		return
	}
	// ゲッターメソッドを使用してプロパティを取得
	x, y := t.GetPosition()
	width, height := t.GetSize()
	currentStyle := t.GetStyle()

	// 最初に背景と境界線を描画します。
	DrawStyledBackground(screen, x, y, width, height, currentStyle)

	// 次にその上にテキストを描画します。
	DrawAlignedText(screen, t.text, image.Rect(x, y, x+width, y+height), currentStyle)
}

// CalculateMinSize は、現在のテキストとスタイルに基づいてウィジェットが表示されるべき最小サイズを計算します。
// この値はレイアウトシステムによって利用されます。
func (t *TextWidget) CalculateMinSize() (int, int) {
	s := t.GetStyle()
	// テキストとフォントが存在する場合のみ、コンテンツに基づいたサイズを計算します。
	if t.text != "" && s.Font != nil && *s.Font != nil {
		bounds := text.BoundString(*s.Font, t.text)

		// パディングの値を取得（nilの場合はゼロ値として扱います）
		padding := s.Padding
		if s.Padding == nil {
			padding = &style.Insets{}
		}

		// テキストの幅と高さにパディングを加えたものを最小サイズとします。
		minWidth := bounds.Dx() + padding.Left + padding.Right
		metrics := (*s.Font).Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

		// ユーザーによって明示的に設定された最小サイズ(minWidth, minHeight)がある場合は、
		// 計算値と比べて大きい方を最終的な最小サイズとします。
		if t.minWidth > minWidth {
			minWidth = t.minWidth
		}
		if t.minHeight > minHeight {
			minHeight = t.minHeight
		}

		return minWidth, minHeight
	}
	// テキストがない場合でも、ユーザー設定の最小サイズは尊重します。
	return t.minWidth, t.minHeight
}