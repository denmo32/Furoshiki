package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/theme"

	"github.com/hajimehoshi/ebiten/v2"
)

// Button is a clickable UI element.
type Button struct {
	*component.TextWidget
	component.InteractiveMixin
}

// Compile-time check to ensure Button implements the interactiveTextWidget interface.
var _ interactiveTextWidget = (*Button)(nil)

// NewButton creates a new instance of a Button widget.
func NewButton(text string) *Button {
	button := &Button{}
	button.TextWidget = component.NewTextWidget(text)

	// Initialize the interactive mixin with styles from the theme.
	t := theme.GetCurrent()
	themeStyles := map[component.WidgetState]style.Style{
		component.StateNormal:   t.Button.Normal,
		component.StateHovered:  t.Button.Hovered,
		component.StatePressed:  t.Button.Pressed,
		component.StateDisabled: t.Button.Disabled,
	}
	button.InteractiveMixin.InitStyles(themeStyles)

	// Set the base style and default size.
	button.SetStyle(button.StateStyles[component.StateNormal])
	button.SetSize(100, 40)

	// IMPORTANT: Initialize the LayoutableWidget with the final concrete type.
	button.Init(button)

	return button
}

// Draw renders the Button. It selects a style based on the current state.
func (b *Button) Draw(screen *ebiten.Image) {
	styleToUse := b.GetActiveStyle(b.CurrentState(), b.GetStyle())
	b.TextWidget.DrawWithStyle(screen, styleToUse)
}

// SetStyleForState sets the style for a specific widget state.
func (b *Button) SetStyleForState(state component.WidgetState, s style.Style) {
	b.InteractiveMixin.SetStyleForState(state, s, b.SetStyle)
}

// StyleAllStates applies a style modification to all interaction states.
func (b *Button) StyleAllStates(s style.Style) {
	b.InteractiveMixin.SetAllStyles(s, b.SetStyle)
}

// --- ButtonBuilder ---

// ButtonBuilder is a fluent builder for creating Button widgets.
type ButtonBuilder struct {
	InteractiveTextBuilder[*ButtonBuilder, *Button]
}

// NewButtonBuilder creates a new builder for a Button.
func NewButtonBuilder() *ButtonBuilder {
	button := NewButton("")
	b := &ButtonBuilder{}
	// Initialize the embedded InteractiveTextBuilder.
	b.InteractiveTextBuilder.Init(b, button)
	return b
}

// Build finalizes the construction of the Button.
func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build()
}
