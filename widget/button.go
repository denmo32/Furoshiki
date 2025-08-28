package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2/text"
)

// Button is a clickable UI element, refactored to use composition over inheritance.
type Button struct {
	// UPDATE: 複数の共通コンポーネントをWidgetCoreに集約
	*component.WidgetCore
	*component.Appearance
	*component.Interaction
	*component.Text
}

// --- Interface implementation verification ---
// UPDATE: 新しい複合インターフェース TextBasedWidget の実装を検証
var _ component.TextBasedWidget = (*Button)(nil)
var _ event.EventTarget = (*Button)(nil)

// newButton creates a new component-based Button.
func newButton(text string) (*Button, error) {
	b := &Button{}
	// UPDATE: WidgetCoreと、Buttonに特有のコンポーネントを初期化
	b.WidgetCore = component.NewWidgetCore(b)
	b.Appearance = component.NewAppearance(b)
	b.Interaction = component.NewInteraction(b)
	b.Text = component.NewText(b, text)

	t := theme.GetCurrent()
	b.SetStyle(t.Button.Normal)
	b.SetStyleForState(component.StateHovered, t.Button.Hovered)
	b.SetStyleForState(component.StatePressed, t.Button.Pressed)
	b.SetStyleForState(component.StateDisabled, t.Button.Disabled)

	b.SetSize(100, 40)
	return b, nil
}

// --- Interface implementations ---

// GetNode, GetLayoutProperties are now inherited from WidgetCore.
func (b *Button) Update()  {}
func (b *Button) Cleanup() { b.SetParent(nil) }

func (b *Button) Draw(info component.DrawInfo) {
	if !b.IsVisible() || !b.HasBeenLaidOut() {
		return
	}
	x, y := b.GetPosition()
	width, height := b.GetSize()
	finalX := x + info.OffsetX
	finalY := y + info.OffsetY

	styleToUse := b.GetStyleForState(b.CurrentState())

	component.DrawStyledBackground(info.Screen, finalX, finalY, width, height, styleToUse)
	finalRect := image.Rect(finalX, finalY, finalX+width, finalY+height)
	component.DrawAlignedText(info.Screen, b.Text.Text(), finalRect, styleToUse, b.WrapText())
}

// UPDATE: MarkDirty, SetPosition, SetSize, SetMinSizeはWidgetCoreに実装されているため削除。

// GetMinSizeは、コンテンツ（テキスト）が要求する最小サイズと、
// ユーザーが設定した最小サイズのうち、大きい方を返します。
// これにより、テキストがはみ出さないことが保証されます。
func (b *Button) GetMinSize() (int, int) {
	contentMinWidth, contentMinHeight := b.calculateContentMinSize()
	userMinWidth, userMinHeight := b.WidgetCore.GetMinSize()
	return max(contentMinWidth, userMinWidth), max(contentMinHeight, userMinHeight)
}

func (b *Button) GetHeightForWidth(width int) int {
	// ラベルと同様に、折り返しが無効な場合はコンテンツの最小高さに依存します。
	if !b.WrapText() {
		_, h := b.calculateContentMinSize()
		return h
	}

	// テキストとスタイル情報を取得します。
	s := b.ReadOnlyStyle()
	if b.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		// フォントがない場合は高さを計算できません。
		return 0
	}

	// パディングを考慮して、テキストが描画される実際の幅を計算します。
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	contentWidth := width - padding.Left - padding.Right

	// コンテンツ幅が0以下の場合、最小高さでフォールバックします。
	if contentWidth <= 0 {
		_, h := b.calculateContentMinSize()
		return h
	}

	// 折り返されたテキストの高さを計算し、垂直パディングを加算します。
	_, requiredHeight := component.CalculateWrappedText(*s.Font, b.Text.Text(), contentWidth)
	return requiredHeight + padding.Top + padding.Bottom
}

// calculateContentMinSizeは、テキストとパディングに基づいて、
// ボタンがコンテンツを表示するために最低限必要なサイズを計算します。
func (b *Button) calculateContentMinSize() (int, int) {
	s := b.ReadOnlyStyle()
	if b.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0, 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	metrics := (*s.Font).Metrics()
	// 最小の高さは、フォントの高さと垂直パディングの合計です。
	contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom

	if b.WrapText() {
		// 折り返しが有効な場合、最小幅は最も長い単語の幅によって決まります。
		// これにより、単語の途中で改行されるのを防ぎます。
		longestWord := ""
		words := utils.SplitIntoWords(b.Text.Text())
		for _, word := range words {
			if len(word) > len(longestWord) {
				longestWord = word
			}
		}
		if longestWord == "" {
			longestWord = b.Text.Text()
		}
		bounds := text.BoundString(*s.Font, longestWord)
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	} else {
		// 折り返しが無効な場合、最小幅はテキスト全体の幅です。
		bounds := text.BoundString(*s.Font, b.Text.Text())
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	}
}

func (b *Button) HitTest(x, y int) component.Widget {
	if !b.IsVisible() || b.IsDisabled() {
		return nil
	}
	wx, wy := b.GetPosition()
	wwidth, wheight := b.GetSize()
	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if rect.Empty() {
		return nil
	}
	if !(image.Point{X: x, Y: y}.In(rect)) {
		return nil
	}
	return b
}

// --- EventTarget and EventProcessor Implementation ---
func (b *Button) HandleEvent(e *event.Event) {
	b.Interaction.TriggerHandlers(e)

	// 親ウィジェットへのイベント伝播
	if e != nil && !e.Handled && b.GetParent() != nil {
		if processor, ok := b.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// UPDATE: SetRequestedPositionとGetRequestedPositionはWidgetCoreに実装されているため削除。

// --- ButtonBuilder ---

type ButtonBuilder struct {
	component.Builder[*ButtonBuilder, *Button]
}

// NewButtonBuilderは新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	button, err := newButton("")
	b := &ButtonBuilder{}
	b.Init(b, button)
	b.AddError(err)
	return b
}

func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build()
}

// --- Button-specific Builder Methods ---

func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.Widget.SetText(text)
	return b
}

func (b *ButtonBuilder) WrapText(wrap bool) *ButtonBuilder {
	b.Widget.SetWrapText(wrap)
	return b
}