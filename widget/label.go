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
// 以前のTextWidgetを埋め込む代わりに、必要なコンポーネントを直接合成します。
type Label struct {
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Interaction
	*component.Text
	*component.Visibility
	*component.Dirty

	// UPDATE: hasBeenLaidOutフィールドはVisibilityコンポーネントに統合されたため削除されました。
	// hasBeenLaidOut bool
	minSize component.Size
}

// --- インターフェース実装の検証 ---
var _ component.Widget = (*Label)(nil)
var _ component.NodeOwner = (*Label)(nil)
var _ component.AppearanceOwner = (*Label)(nil)
var _ component.InteractionOwner = (*Label)(nil)
var _ component.TextOwner = (*Label)(nil)
var _ component.LayoutPropertiesOwner = (*Label)(nil)
var _ component.VisibilityOwner = (*Label)(nil)
var _ component.DirtyManager = (*Label)(nil)
var _ component.HeightForWider = (*Label)(nil)
var _ component.AbsolutePositioner = (*Label)(nil)
// NOTE: Label is not interactive, but it must implement EventProcessor for event propagation.
var _ component.EventProcessor = (*Label)(nil)
// UPDATE: Buildableインターフェースが削除されたため、実装検証も削除
// var _ component.Buildable = (*Label)(nil)


// newLabelは、新しいコンポーネントベースのLabelを生成します。
func newLabel(text string) (*Label, error) {
	l := &Label{}
	l.Node = component.NewNode(l)
	l.Transform = component.NewTransform()
	l.LayoutProperties = component.NewLayoutProperties()
	l.Appearance = component.NewAppearance(l)
	// NOTE: Labelはインタラクティブではないが、コンポーネントの枠組みとしてInteractionを持つ
	l.Interaction = component.NewInteraction(l)
	l.Text = component.NewText(l, text)
	l.Visibility = component.NewVisibility(l)
	l.Dirty = component.NewDirty()

	t := theme.GetCurrent()
	l.SetStyle(t.Label.Default)

	// NOTE: デフォルトサイズの指定は、具象ウィジェットの責任として残します。
	l.SetSize(100, 30)
	return l, nil
}

// --- インターフェース実装 ---

func (l *Label) GetNode() *component.Node                   { return l.Node }
func (l *Label) GetLayoutProperties() *component.LayoutProperties { return l.LayoutProperties }
func (l *Label) Update()                                    {}
func (l *Label) Cleanup()                                   { l.SetParent(nil) }

func (l *Label) Draw(info component.DrawInfo) {
	// UPDATE: hasBeenLaidOutのチェックをVisibilityコンポーネントのHasBeenLaidOut()に置き換えました。
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

func (l *Label) MarkDirty(relayout bool) {
	l.Dirty.MarkDirty(relayout)
	if relayout && !l.IsLayoutBoundary() {
		if parent := l.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (l *Label) SetPosition(x, y int) {
	// UPDATE: レイアウト済み状態の管理をVisibilityコンポーネントに委譲します。
	if !l.HasBeenLaidOut() {
		l.SetLaidOut(true)
	}
	if posX, posY := l.GetPosition(); posX != x || posY != y {
		l.Transform.SetPosition(x, y)
		l.MarkDirty(false)
	}
}

func (l *Label) SetSize(width, height int) {
	if w, h := l.GetSize(); w != width || h != height {
		l.Transform.SetSize(width, height)
		l.MarkDirty(true)
	}
}

func (l *Label) SetMinSize(width, height int) {
	l.minSize.Width = width
	l.minSize.Height = height
}

func (l *Label) GetMinSize() (int, int) {
	contentMinWidth, contentMinHeight := l.calculateContentMinSize()
	return max(contentMinWidth, l.minSize.Width), max(contentMinHeight, l.minSize.Height)
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

// HandleEvent is a dummy implementation to satisfy the EventProcessor interface.
func (l *Label) HandleEvent(e *event.Event) {
	// Propagate event to parent by default
	if e != nil && !e.Handled && l.GetParent() != nil {
		if processor, ok := l.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// --- AbsolutePositioner and other interface implementations required by Builder ---
func (l *Label) SetRequestedPosition(x, y int) {
	l.Transform.SetRequestedPosition(x, y)
	l.MarkDirty(true)
}

func (l *Label) GetRequestedPosition() (int, int) {
	return l.Transform.GetRequestedPosition()
}

// NOTE: AddEventHandler/RemoveEventHandler は、EventProcessorインターフェースを満たすために必要です。
// Label自体はイベントを発行しませんが、イベントがツリーを伝播できるようにするために実装します。
func (l *Label) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	l.Interaction.AddEventHandler(eventType, handler)
}
func (l *Label) RemoveEventHandler(eventType event.EventType) {
	l.Interaction.RemoveEventHandler(eventType)
}


// --- LabelBuilder ---

// 【提案3対応】LabelBuilderは汎用のcomponent.Builderを埋め込むように変更されました。
// これにより、Size, Flex, Paddingなどの共通メソッドは基底クラスに集約され、
// LabelBuilderはLabel固有のメソッド（Textなど）のみを定義します。
type LabelBuilder struct {
	component.Builder[*LabelBuilder, *Label]
}

// NewLabelBuilderは新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	label, err := newLabel("")
	b := &LabelBuilder{}
	// 基底のビルダーを初期化します。
	b.Init(b, label)
	b.AddError(err)
	return b
}

// Buildは、最終的なLabelを構築して返します。
// 実際のロジックは基底のBuilder.Buildに移譲されます。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}

// --- Label-specific Builder Methods ---

// Text はラベルに表示されるテキストを設定します。
func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.Widget.SetText(text)
	return b
}

// WrapText はラベルのテキストを折り返すかどうかを設定します。
func (b *LabelBuilder) WrapText(wrap bool) *LabelBuilder {
	b.Widget.SetWrapText(wrap)
	return b
}

// NOTE: TextColor, BackgroundColor, Padding, VerticalAlign, TextAlign, Size, Flex, AbsolutePosition, AssignTo
// などの共通メソッドはすべて component.Builder に実装されているため、ここからは完全に削除されました。