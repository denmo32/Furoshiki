package widget

import (
	"errors"
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/utils"
	"image"
	"image/color"
	"log"
	"runtime/debug"

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

	hasBeenLaidOut bool
	minSize        component.Size
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

	b.SetSize(100, 40)
	return b, nil
}

// --- Interface implementations ---

func (b *Button) GetNode() *component.Node                   { return b.Node }
func (b *Button) GetLayoutProperties() *component.LayoutProperties { return b.LayoutProperties }
func (b *Button) Update()                                    {}
func (b *Button) Cleanup()                                   { b.SetParent(nil) }

func (b *Button) Draw(info component.DrawInfo) {
	if !b.IsVisible() || !b.hasBeenLaidOut {
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
	if !b.hasBeenLaidOut {
		b.hasBeenLaidOut = true
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
	if handlers, exists := b.GetEventHandlers()[e.Type]; exists {
		for _, handler := range handlers {
			if e.Handled {
				break
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf(`Recovered from panic in event handler: %v
%s`, r, debug.Stack())
					}
				}()
				if handler(e) == event.StopPropagation {
					e.Handled = true
				}
			}()
		}
	}

	if e != nil && !e.Handled && b.GetParent() != nil {
		if processor, ok := b.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}

// --- AbsolutePositioner Implementation ---
func (b *Button) SetRequestedPosition(x, y int) {
	b.Transform.SetRequestedPosition(x, y)
	b.MarkDirty(true)
}

func (b *Button) GetRequestedPosition() (int, int) {
	return b.Transform.GetRequestedPosition()
}

// --- ButtonBuilder ---
type ButtonBuilder struct {
	button *Button
	errors []error
}

func NewButtonBuilder() *ButtonBuilder {
	button, err := newButton("")
	b := &ButtonBuilder{button: button}
	if err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *ButtonBuilder) Build() (*Button, error) {
	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}
	b.button.MarkDirty(true)
	return b.button, nil
}

func (b *ButtonBuilder) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// --- Builder Methods ---

func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.button.SetText(text)
	return b
}

func (b *ButtonBuilder) WrapText(wrap bool) *ButtonBuilder {
	b.button.SetWrapText(wrap)
	return b
}

func (b *ButtonBuilder) Size(width, height int) *ButtonBuilder {
	b.button.SetSize(width, height)
	return b
}

func (b *ButtonBuilder) MinSize(width, height int) *ButtonBuilder {
	b.button.SetMinSize(width, height)
	return b
}

func (b *ButtonBuilder) Flex(flex int) *ButtonBuilder {
	b.button.SetFlex(flex)
	return b
}

func (b *ButtonBuilder) AbsolutePosition(x, y int) *ButtonBuilder {
	b.button.SetRequestedPosition(x, y)
	return b
}

func (b *ButtonBuilder) AssignTo(target **Button) *ButtonBuilder {
	if target == nil {
		b.errors = append(b.errors, errors.New("AssignTo target cannot be nil"))
		return b
	}
	*target = b.button
	return b
}

func (b *ButtonBuilder) AddOnClick(handler event.EventHandler) *ButtonBuilder {
	b.button.AddEventHandler(event.EventClick, handler)
	return b
}

// --- Style Helper Methods ---

func (b *ButtonBuilder) applyStyle(s style.Style) *ButtonBuilder {
	existingStyle := b.button.GetStyle()
	b.button.SetStyle(style.Merge(existingStyle, s))
	return b
}

func (b *ButtonBuilder) TextColor(c color.Color) *ButtonBuilder {
	return b.applyStyle(style.Style{TextColor: style.PColor(c)})
}

func (b *ButtonBuilder) BackgroundColor(c color.Color) *ButtonBuilder {
	return b.applyStyle(style.Style{Background: style.PColor(c)})
}

func (b *ButtonBuilder) Padding(p int) *ButtonBuilder {
	return b.applyStyle(style.Style{Padding: style.PInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})})
}

func (b *ButtonBuilder) Border(width float32, c color.Color) *ButtonBuilder {
	return b.applyStyle(style.Style{BorderWidth: style.PFloat32(width), BorderColor: style.PColor(c)})
}

func (b *ButtonBuilder) TextAlign(align style.TextAlignType) *ButtonBuilder {
	return b.applyStyle(style.Style{TextAlign: style.PTextAlignType(align)})
}

func (b *ButtonBuilder) VerticalAlign(align style.VerticalAlignType) *ButtonBuilder {
	return b.applyStyle(style.Style{VerticalAlign: style.PVerticalAlignType(align)})
}