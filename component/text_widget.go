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
func NewTextWidget(text string) *TextWidget {
	tw := &TextWidget{
		LayoutableWidget: NewLayoutableWidget(),
		text:             text,
	}

	// コンテンツの最小サイズを計算する責務を、クロージャとしてLayoutableWidgetに委譲します。
	// これにより、最小サイズ決定のロジックが基底ウィジェットに集約され、
	// TextWidgetはGetMinSizeメソッドをオーバーライドする必要がなくなります。
	tw.LayoutableWidget.contentMinSizeFunc = tw.calculateContentMinSize

	return tw
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

// DrawWithStyleは、指定されたスタイルを用いてウィジェットの背景とテキストを描画する共通ロジックです。
// 通常のDrawメソッドと分離することで、Buttonのように状態に応じてスタイルを切り替える必要のある
// 具象ウィジェットが、描画ロジックを再利用しやすくなります。
func (t *TextWidget) DrawWithStyle(screen *ebiten.Image, styleToUse style.Style) {
	// IsVisible() に加えてレイアウト済みかもチェックします。
	// これにより、ウィジェットがUIツリーに追加されてから最初のレイアウト計算が完了するまでの1フレーム間、
	// 意図せず (0,0) 座標に描画されてしまうのを防ぎます。
	if !t.IsVisible() || !t.HasBeenLaidOut() {
		return
	}

	x, y := t.GetPosition()
	width, height := t.GetSize()

	DrawStyledBackground(screen, x, y, width, height, styleToUse)
	DrawAlignedText(screen, t.text, image.Rect(x, y, x+width, y+height), styleToUse)
}

// Draw はTextWidgetを描画します。
// このメソッドは、ウィジェット自身の現在のスタイルを使用して、共通の描画ロジック(DrawWithStyle)を呼び出します。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	currentStyle := t.GetStyle()
	t.DrawWithStyle(screen, currentStyle)
}

// calculateContentMinSize は、現在のテキストとスタイルに基づいてコンテンツが表示されるべき最小サイズを計算します。
func (t *TextWidget) calculateContentMinSize() (int, int) {
	s := t.GetStyle()
	if t.text != "" && s.Font != nil && *s.Font != nil {
		bounds := text.BoundString(*s.Font, t.text)

		padding := style.Insets{}
		if s.Padding != nil {
			padding = *s.Padding
		}

		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		metrics := (*s.Font).Metrics()
		contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

		return contentMinWidth, contentMinHeight
	}
	// テキストがない場合は、コンテンツの最小サイズは0です。
	return 0, 0
}
