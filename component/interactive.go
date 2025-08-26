package component

import "furoshiki/style"

// InteractiveMixin provides style management for stateful widgets.
// It is intended to be embedded in a concrete widget.
type InteractiveMixin struct {
	StateStyles map[WidgetState]style.Style
}

// InitStyles initializes the mixin with styles from a theme.
func (im *InteractiveMixin) InitStyles(styles map[WidgetState]style.Style) {
	im.StateStyles = make(map[WidgetState]style.Style)
	for state, s := range styles {
		im.StateStyles[state] = s.DeepCopy()
	}
}

// GetActiveStyle determines the style to use based on the widget's current state.
func (im *InteractiveMixin) GetActiveStyle(currentState WidgetState, defaultStyle style.Style) style.Style {
	if style, ok := im.StateStyles[currentState]; ok {
		return style
	}
	return defaultStyle
}

// SetStyleForState merges a new style for a specific state.
// The normalSetter function is a callback to update the widget's base style when the normal state style changes.
func (im *InteractiveMixin) SetStyleForState(state WidgetState, s style.Style, normalSetter func(style.Style)) {
	baseStyle := im.StateStyles[state]
	im.StateStyles[state] = style.Merge(baseStyle, s)
	if state == StateNormal {
		normalSetter(im.StateStyles[StateNormal])
	}
}

// SetAllStyles merges a new style across all states.
func (im *InteractiveMixin) SetAllStyles(s style.Style, normalSetter func(style.Style)) {
	for state, baseStyle := range im.StateStyles {
		im.StateStyles[state] = style.Merge(baseStyle, s)
	}
	normalSetter(im.StateStyles[StateNormal])
}