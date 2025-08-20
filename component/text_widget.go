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

// [改良] calculateContentMinSize は、現在のテキストとスタイルに基づいてコンテンツが表示されるべき最小サイズを計算します。
// このメソッドはGetMinSizeから内部的に呼び出されます。
func (t *TextWidget) calculateContentMinSize() (int, int) {
	s := t.GetStyle()
	// テキストとフォントが存在する場合のみ、コンテンツに基づいたサイズを計算します。
	if t.text != "" && s.Font != nil && *s.Font != nil {
		bounds := text.BoundString(*s.Font, t.text)

		// パディングの値を取得（nilの場合はゼロ値として扱います）
		padding := style.Insets{}
		if s.Padding != nil {
			padding = *s.Padding
		}

		// テキストの幅と高さにパディングを加えたものをコンテンツの最小サイズとします。
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		metrics := (*s.Font).Metrics()
		contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

		return contentMinWidth, contentMinHeight
	}
	// テキストがない場合は、コンテンツの最小サイズは0です。
	return 0, 0
}

// [改良] GetMinSize は、ウィジェットが表示されるべき最小サイズを返します。
// LayoutableWidgetのGetMinSizeをオーバーライドし、テキストコンテンツのサイズを考慮に入れます。
// 最終的な最小サイズは、「コンテンツから計算されるサイズ」と「ユーザーが明示的に設定した最小サイズ」のうち、大きい方になります。
func (t *TextWidget) GetMinSize() (int, int) {
	// コンテンツ（テキストとパディング）から要求される最小サイズを計算します。
	contentMinWidth, contentMinHeight := t.calculateContentMinSize()

	// ユーザーが .MinSize() で明示的に設定した最小サイズを取得します。
	userMinWidth, userMinHeight := t.LayoutableWidget.GetMinSize()

	// 両者を比較し、各次元で大きい方を最終的な最小サイズとします。
	finalMinWidth := contentMinWidth
	if userMinWidth > contentMinWidth {
		finalMinWidth = userMinWidth
	}

	finalMinHeight := contentMinHeight
	if userMinHeight > contentMinHeight {
		finalMinHeight = userMinHeight
	}

	return finalMinWidth, finalMinHeight
}