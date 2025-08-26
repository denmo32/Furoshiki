package widget

import (
	"furoshiki/component"
	"furoshiki/style"
)

// interactiveTextWidget defines the behavior for widgets using the InteractiveTextBuilder.
type interactiveTextWidget interface {
	textWidget // From widget/builder.go
	SetStyleForState(state component.WidgetState, s style.Style)
	StyleAllStates(s style.Style)
}

// InteractiveTextBuilder provides a reusable builder for interactive text widgets.
type InteractiveTextBuilder[T any, W interactiveTextWidget] struct {
	Builder[T, W] // Embeds the text widget builder
}

// SetStyleForState sets the style for a specific widget state.
func (b *InteractiveTextBuilder[T, W]) SetStyleForState(state component.WidgetState, s style.Style) T {
	b.Widget.SetStyleForState(state, s)
	return b.Self
}

// StyleAllStates applies a style modification to all interaction states.
func (b *InteractiveTextBuilder[T, W]) StyleAllStates(s style.Style) T {
	b.Widget.StyleAllStates(s)
	return b.Self
}

// HoverStyle sets the style for the Hovered state.
func (b *InteractiveTextBuilder[T, W]) HoverStyle(s style.Style) T {
	return b.SetStyleForState(component.StateHovered, s)
}

// PressedStyle sets the style for the Pressed state.
func (b *InteractiveTextBuilder[T, W]) PressedStyle(s style.Style) T {
	return b.SetStyleForState(component.StatePressed, s)
}

// DisabledStyle sets the style for the Disabled state.
func (b *InteractiveTextBuilder[T, W]) DisabledStyle(s style.Style) T {
	return b.SetStyleForState(component.StateDisabled, s)
}

// Style overrides the base Style method to set the style for the Normal state.
func (b *InteractiveTextBuilder[T, W]) Style(s style.Style) T {
	return b.SetStyleForState(component.StateNormal, s)
}
