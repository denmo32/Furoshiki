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

// GetActiveStyle は、ウィジェットの現在の状態に基づいて適用すべきスタイルを決定します。
// 内部のStateStylesマップからスタイルを検索し、見つからない場合はNormal状態のスタイルをフォールバックとして返します。
func (im *InteractiveMixin) GetActiveStyle(currentState WidgetState) style.Style {
	if style, ok := im.StateStyles[currentState]; ok {
		return style
	}
	// フォールバックとして、必ず存在するはずのNormalスタイルを返す。
	// これにより、呼び出し側がデフォルトスタイルを渡す必要がなくなり、不要なコピーを防げます。
	return im.StateStyles[StateNormal]
}

// SetStyleForState は指定された状態のスタイルを、既存のスタイルにマージします。
// このメソッドは InteractiveMixin 内部のスタイルマップのみを更新し、
// ウィジェットの基本スタイルへの反映は呼び出し元の責務です。
func (im *InteractiveMixin) SetStyleForState(state WidgetState, s style.Style) {
	baseStyle := im.StateStyles[state]
	mergedStyle := style.Merge(baseStyle, s)
	im.StateStyles[state] = mergedStyle
}

// SetAllStyles はすべての状態に新しいスタイルをマージします。
// このメソッドは InteractiveMixin 内部のスタイルマップのみを更新し、
// ウィジェットの基本スタイルへの反映は呼び出し元の責務です。
func (im *InteractiveMixin) SetAllStyles(s style.Style) {
	for state, baseStyle := range im.StateStyles {
		im.StateStyles[state] = style.Merge(baseStyle, s)
	}
}