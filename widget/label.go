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

// Labelはテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	// UPDATE: 複数の共通コンポーネントをWidgetCoreに集約
	*component.WidgetCore
	*component.Appearance
	*component.Interaction
	*component.Text
}

// --- インターフェース実装の検証 ---
// UPDATE: 新しい複合インターフェース TextBasedWidget の実装を検証
var _ component.TextBasedWidget = (*Label)(nil)

// newLabelは、新しいコンポーネントベースのLabelを生成します。
func newLabel(text string) (*Label, error) {
	l := &Label{}
	// UPDATE: WidgetCoreと、Labelに特有のコンポーネントを初期化
	l.WidgetCore = component.NewWidgetCore(l)
	l.Appearance = component.NewAppearance(l)
	l.Interaction = component.NewInteraction(l) // Label is not interactive, but this is for consistency and event propagation
	l.Text = component.NewText(l, text)

	t := theme.GetCurrent()
	l.SetStyle(t.Label.Default)

	l.SetSize(100, 30)
	return l, nil
}

// --- インターフェース実装 ---

// GetNode, GetLayoutProperties are now inherited from WidgetCore.
func (l *Label) Update()  {}
func (l *Label) Cleanup() { l.SetParent(nil) }

func (l *Label) Draw(info component.DrawInfo) {
	if !l.IsVisible() || !l.HasBeenLaidOut() {
		return
	}
	x, y := l.GetPosition()
	width, height := l.GetSize()
	finalX := x + info.OffsetX
	finalY := y + info.OffsetY

	styleToUse := l.GetStyle()

	component.DrawStyledBackground(info.Screen, finalX, finalY, width, height, styleToUse)
	finalRect := image.Rect(finalX, finalY, finalX+width, finalY+height)
	component.DrawAlignedText(info.Screen, l.Text.Text(), finalRect, styleToUse, l.WrapText())
}

// UPDATE: MarkDirty, SetPosition, SetSize, SetMinSizeはWidgetCoreに実装されているため削除。

// GetMinSizeは、コンテンツ（テキスト）が要求する最小サイズと、
// ユーザーが設定した最小サイズのうち、大きい方を返します。
func (l *Label) GetMinSize() (int, int) {
	contentMinWidth, contentMinHeight := l.calculateContentMinSize()
	userMinWidth, userMinHeight := l.WidgetCore.GetMinSize()
	return max(contentMinWidth, userMinWidth), max(contentMinHeight, userMinHeight)
}

func (l *Label) GetHeightForWidth(width int) int {
	if !l.WrapText() {
		_, h := l.calculateContentMinSize()
		return h
	}
	s := l.ReadOnlyStyle()
	if l.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	contentWidth := width - padding.Left - padding.Right
	if contentWidth <= 0 {
		_, h := l.calculateContentMinSize()
		return h
	}
	_, requiredHeight := component.CalculateWrappedText(*s.Font, l.Text.Text(), contentWidth)
	return requiredHeight + padding.Top + padding.Bottom
}

// calculateContentMinSizeは、テキストとパディングに基づいて、
// ラベルがコンテンツを表示するために最低限必要なサイズを計算します。
func (l *Label) calculateContentMinSize() (int, int) {
	s := l.ReadOnlyStyle()
	if l.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0, 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	metrics := (*s.Font).Metrics()
	contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom
	if l.WrapText() {
		longestWord := ""
		words := utils.SplitIntoWords(l.Text.Text())
		for _, word := range words {
			if len(word) > len(longestWord) {
				longestWord = word
			}
		}
		if longestWord == "" {
			longestWord = l.Text.Text()
		}
		bounds := text.BoundString(*s.Font, longestWord)
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	} else {
		bounds := text.BoundString(*s.Font, l.Text.Text())
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	}
}

func (l *Label) HitTest(x, y int) component.Widget {
	// Label is not interactive by default.
	return nil
}

// HandleEventはEventProcessorインターフェースを満たすための実装です。
// Labelは自身ではイベントを処理しませんが、イベントを親へ伝播させる役割を持ちます。
func (l *Label) HandleEvent(e *event.Event) {
	if e != nil && !e.Handled && l.GetParent() != nil {
		if processor, ok := l.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// UPDATE: SetRequestedPositionとGetRequestedPositionはWidgetCoreに実装されているため削除。

// AddEventHandlerとRemoveEventHandlerはEventProcessorインターフェースを満たすために必要です。
// Interactionコンポーネントに処理を委譲します。
func (l *Label) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	l.Interaction.AddEventHandler(eventType, handler)
}
func (l *Label) RemoveEventHandler(eventType event.EventType) {
	l.Interaction.RemoveEventHandler(eventType)
}

// --- LabelBuilder ---

type LabelBuilder struct {
	component.Builder[*LabelBuilder, *Label]
}

// NewLabelBuilderは新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	label, err := newLabel("")
	b := &LabelBuilder{}
	b.Init(b, label)
	b.AddError(err)
	return b
}

func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}

// --- Label-specific Builder Methods ---

func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.Widget.SetText(text)
	return b
}

func (b *LabelBuilder) WrapText(wrap bool) *LabelBuilder {
	b.Widget.SetWrapText(wrap)
	return b
}