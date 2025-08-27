package component

import (
	"furoshiki/style"
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// --- TextWidget ---
// TextWidgetは、テキスト表示に関連する共通の機能（テキスト内容、スタイル、最小サイズ計算）を提供します。
// ButtonやLabelなど、テキストを持つウィジェットはこれを埋め込みます。
type TextWidget struct {
	*LayoutableWidget
	text     string
	wrapText bool // テキストを折り返すかどうか
}

// コンパイル時にインターフェースの実装を検証します。
var _ HeightForWider = (*TextWidget)(nil)

// NewTextWidget は新しいTextWidgetを生成します。
func NewTextWidget(text string) *TextWidget {
	tw := &TextWidget{
		LayoutableWidget: NewLayoutableWidget(),
		text:             text,
		wrapText:         false, // デフォルトでは折り返さない
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

// SetWrapText はテキストの折り返し設定を変更します。
func (t *TextWidget) SetWrapText(wrap bool) {
	if t.wrapText != wrap {
		t.wrapText = wrap
		// 折り返し設定の変更はレイアウトに影響するため、再レイアウトを要求します。
		t.MarkDirty(true)
	}
}

// GetHeightForWidth は、HeightForWiderインターフェースの実装です。
// 指定された幅に基づいて、テキストを折り返した場合に必要となる高さを計算します。
func (t *TextWidget) GetHeightForWidth(width int) int {
	if !t.wrapText {
		// 折り返しが無効な場合、通常の最小高さを返します。
		_, h := t.calculateContentMinSize()
		return h
	}

	s := t.ReadOnlyStyle()
	if t.text == "" || s.Font == nil || *s.Font == nil {
		return 0
	}

	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}

	contentWidth := width - padding.Left - padding.Right
	if contentWidth <= 0 {
		_, h := t.calculateContentMinSize() // 幅がない場合は1行の高さを返す
		return h
	}

	// drawing_helpers.goの公開関数を呼び出します。
	_, requiredHeight := CalculateWrappedText(*s.Font, t.text, contentWidth)
	return requiredHeight + padding.Top + padding.Bottom
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
	// 描画ヘルパーに折り返し情報を渡します。
	DrawAlignedText(screen, t.text, image.Rect(x, y, x+width, y+height), styleToUse, t.wrapText)
}

// Draw はTextWidgetを描画します。
// このメソッドは、ウィジェット自身の現在のスタイルを使用して、共通の描画ロジック(DrawWithStyle)を呼び出します。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	// NOTE: パフォーマンス向上のため、ディープコピーを行うGetStyle()の代わりに
	//       シャローコピーを行うReadOnlyStyle()を使用します。
	currentStyle := t.ReadOnlyStyle()
	t.DrawWithStyle(screen, currentStyle)
}

// calculateContentMinSize は、現在のテキストとスタイルに基づいてコンテンツが表示されるべき最小サイズを計算します。
func (t *TextWidget) calculateContentMinSize() (int, int) {
	s := t.ReadOnlyStyle()
	if t.text == "" || s.Font == nil || *s.Font == nil {
		// テキストがない場合は、コンテンツの最小サイズは0です。
		return 0, 0
	}

	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}

	metrics := (*s.Font).Metrics()
	contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

	if t.wrapText {
		// 折り返しが有効な場合、最小幅は最も長い単語の幅になります。
		longestWord := ""
		words := strings.Fields(t.text)
		for _, word := range words {
			if len(word) > len(longestWord) {
				longestWord = word
			}
		}
		if longestWord == "" {
			longestWord = t.text // 空白を含まない長い単一の単語の場合
		}
		bounds := text.BoundString(*s.Font, longestWord)
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	} else {
		// 折り返しが無効な場合、最小幅はテキスト全体の幅になります。
		bounds := text.BoundString(*s.Font, t.text)
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	}
}