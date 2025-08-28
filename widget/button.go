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

// --- Interface implementation verification ---
var _ component.Widget = (*Button)(nil)
var _ component.NodeOwner = (*Button)(nil)
var _ component.AppearanceOwner = (*Button)(nil)
var _ component.InteractionOwner = (*Button)(nil)
var _ component.TextOwner = (*Button)(nil)
var _ component.LayoutPropertiesOwner = (*Button)(nil)
var _ component.VisibilityOwner = (*Button)(nil)
var _ component.DirtyManager = (*Button)(nil)
var _ component.HeightForWider = (*Button)(nil)
var _ event.EventTarget = (*Button)(nil)
var _ component.EventProcessor = (*Button)(nil)
var _ component.AbsolutePositioner = (*Button)(nil)
// UPDATE: Buildableインターフェースが削除されたため、実装検証も削除
// var _ component.Buildable = (*Button)(nil)

// newButton creates a new component-based Button.
func newButton(text string) (*Button, error) {
	b := &Button{}
	b.Node = component.NewNode(b)
	b.Transform = component.NewTransform()
	b.LayoutProperties = component.NewLayoutProperties()
	b.Appearance = component.NewAppearance(b)
	b.Interaction = component.NewInteraction(b)
	b.Text = component.NewText(b, text)
	b.Visibility = component.NewVisibility(b)
	b.Dirty = component.NewDirty()

	t := theme.GetCurrent()
	b.SetStyle(t.Button.Normal)
	b.SetStyleForState(component.StateHovered, t.Button.Hovered)
	b.SetStyleForState(component.StatePressed, t.Button.Pressed)
	b.SetStyleForState(component.StateDisabled, t.Button.Disabled)

	// NOTE: デフォルトサイズの指定は、具象ウィジェットの責任として残します。
	b.SetSize(100, 40)
	return b, nil
}

// --- Interface implementations ---

func (b *Button) GetNode() *component.Node                   { return b.Node }
func (b *Button) GetLayoutProperties() *component.LayoutProperties { return b.LayoutProperties }
func (b *Button) Update()                                    {}
func (b *Button) Cleanup()                                   { b.SetParent(nil) }

func (b *Button) Draw(info component.DrawInfo) {
	// UPDATE: hasBeenLaidOutのチェックをVisibilityコンポーネントのHasBeenLaidOut()に置き換えました。
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

func (b *Button) MarkDirty(relayout bool) {
	b.Dirty.MarkDirty(relayout)
	if relayout && !b.IsLayoutBoundary() {
		if parent := b.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (b *Button) SetPosition(x, y int) {
	// UPDATE: レイアウト済み状態の管理をVisibilityコンポーネントに委譲します。
	if !b.HasBeenLaidOut() {
		b.SetLaidOut(true)
	}
	if posX, posY := b.GetPosition(); posX != x || posY != y {
		b.Transform.SetPosition(x, y)
		b.MarkDirty(false)
	}
}

func (b *Button) SetSize(width, height int) {
	if w, h := b.GetSize(); w != width || h != height {
		b.Transform.SetSize(width, height)
		b.MarkDirty(true)
	}
}

func (b *Button) SetMinSize(width, height int) {
	b.minSize.Width = width
	b.minSize.Height = height
	b.MarkDirty(true)
}

func (b *Button) GetMinSize() (int, int) {
	contentMinWidth, contentMinHeight := b.calculateContentMinSize()
	return max(contentMinWidth, b.minSize.Width), max(contentMinHeight, b.minSize.Height)
}

func (b *Button) GetHeightForWidth(width int) int {
	if !b.WrapText() {
		_, h := b.calculateContentMinSize()
		return h
	}
	s := b.ReadOnlyStyle()
	if b.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	contentWidth := width - padding.Left - padding.Right
	if contentWidth <= 0 {
		_, h := b.calculateContentMinSize()
		return h
	}
	_, requiredHeight := component.CalculateWrappedText(*s.Font, b.Text.Text(), contentWidth)
	return requiredHeight + padding.Top + padding.Bottom
}

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
	contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom
	if b.WrapText() {
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
	// UPDATE: イベントハンドラの安全な実行をInteractionコンポーネントに委譲
	b.Interaction.TriggerHandlers(e)

	// 親ウィジェットへのイベント伝播
	if e != nil && !e.Handled && b.GetParent() != nil {
		if processor, ok := b.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// --- AbsolutePositioner and other interface implementations required by Builder ---
// NOTE: 以下のメソッドは、Builderが型アサーションで動的にチェックするインターフェースを
// 満たすために実装されています。

func (b *Button) SetRequestedPosition(x, y int) {
	b.Transform.SetRequestedPosition(x, y)
	b.MarkDirty(true)
}

func (b *Button) GetRequestedPosition() (int, int) {
	return b.Transform.GetRequestedPosition()
}

// --- ButtonBuilder ---

// 【提案3対応】ButtonBuilderは汎用のcomponent.Builderを埋め込むように変更されました。
// これにより、Size, Flex, Padding, AssignToなどの共通メソッドは基底クラスに集約され、
// ButtonBuilderはButton固有のメソッド（Text, AddOnClickなど）のみを定義します。
type ButtonBuilder struct {
	component.Builder[*ButtonBuilder, *Button]
}

// NewButtonBuilderは新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	button, err := newButton("")
	b := &ButtonBuilder{}
	// 基底のビルダーを初期化します。
	b.Init(b, button)
	b.AddError(err)
	return b
}

// Build は、最終的なButtonを構築して返します。
// 実際のロジックは基底のBuilder.Buildに移譲されます。
func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build()
}

// --- Button-specific Builder Methods ---

// Text はボタンに表示されるテキストを設定します。
func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.Widget.SetText(text)
	return b
}

// WrapText はボタンのテキストを折り返すかどうかを設定します。
func (b *ButtonBuilder) WrapText(wrap bool) *ButtonBuilder {
	b.Widget.SetWrapText(wrap)
	return b
}

// NOTE: Size, Flex, AbsolutePosition, AssignTo, BackgroundColor, Padding, Border, AddOnClick
// などの共通メソッドはすべて component.Builder に実装されているため、ここからは完全に削除されました。