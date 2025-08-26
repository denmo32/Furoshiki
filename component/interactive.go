

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

// SetStyleForState は指定された状態のスタイルをマージします。
// normalStyleUpdated が true で返された場合、呼び出し元はウィジェットの
// ベーススタイルを更新する必要があります。
func (im *InteractiveMixin) SetStyleForState(state WidgetState, s style.Style) (newNormalStyle style.Style, normalStyleUpdated bool) {
	baseStyle := im.StateStyles[state]
	mergedStyle := style.Merge(baseStyle, s)
	im.StateStyles[state] = mergedStyle
	if state == StateNormal {
		// Normal状態が更新されたことを、更新後のスタイルと共に呼び出し元に通知します。
		return mergedStyle, true
	}
	return style.Style{}, false
}

// SetAllStyles はすべての状態に新しいスタイルをマージし、更新後のNormalスタイルを返します。
// 呼び出し元は、この戻り値を使ってウィジェットのベーススタイルを更新する必要があります。
func (im *InteractiveMixin) SetAllStyles(s style.Style) style.Style {
	for state, baseStyle := range im.StateStyles {
		im.StateStyles[state] = style.Merge(baseStyle, s)
	}
	// 更新後のNormalスタイルを返します。
	return im.StateStyles[StateNormal]
}