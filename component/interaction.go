package component

import "furoshiki/event"

// InteractionOwnerは、Interactionコンポーネントを所有するオブジェクトが
// 満たすべきインターフェースを定義します。
// これにより、Interactionは自身のオーナーのダーティ状態を更新できます。
type InteractionOwner interface {
	DirtyManager
}

// Interactionは、ウィジェットの対話的な状態（ホバー、押下など）と
// イベントハンドラを管理します。
// これにより、イベント処理に関連するロジックとデータをカプセル化します。
type Interaction struct {
	owner         InteractionOwner
	isHovered     bool
	isPressed     bool
	isDisabled    bool
	eventHandlers map[event.EventType][]event.EventHandler
}

// NewInteractionは、新しいInteractionコンポーネントを生成します。
func NewInteraction(owner InteractionOwner) *Interaction {
	return &Interaction{
		owner:         owner,
		eventHandlers: make(map[event.EventType][]event.EventHandler),
	}
}

// AddEventHandlerは、指定されたイベントタイプにイベントハンドラを追加します。
func (i *Interaction) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	i.eventHandlers[eventType] = append(i.eventHandlers[eventType], handler)
}

// RemoveEventHandlerは、指定されたイベントタイプのすべてのハンドラを削除します。
func (i *Interaction) RemoveEventHandler(eventType event.EventType) {
	delete(i.eventHandlers, eventType)
}

// GetEventHandlersは、イベントハンドラのマップを返します。
func (i *Interaction) GetEventHandlers() map[event.EventType][]event.EventHandler {
	return i.eventHandlers
}

// SetHoveredはホバー状態を設定し、必要であれば再描画を要求します。
func (i *Interaction) SetHovered(hovered bool) {
	if i.isHovered != hovered {
		i.isHovered = hovered
		i.owner.MarkDirty(false) // スタイル変更のみなので再描画を要求
	}
}

// IsHoveredは現在のホバー状態を返します。
func (i *Interaction) IsHovered() bool {
	return i.isHovered
}

// SetPressedは押下状態を設定し、必要であれば再描画を要求します。
func (i *Interaction) SetPressed(pressed bool) {
	if i.isPressed != pressed {
		i.isPressed = pressed
		i.owner.MarkDirty(false)
	}
}

// IsPressedは現在の押下状態を返します。
func (i *Interaction) IsPressed() bool {
	return i.isPressed
}

// SetDisabledは無効状態を設定し、必要であれば再描画を要求します。
func (i *Interaction) SetDisabled(disabled bool) {
	if i.isDisabled != disabled {
		i.isDisabled = disabled
		i.owner.MarkDirty(false)
	}
}

// IsDisabledは現在の無効状態を返します。
func (i *Interaction) IsDisabled() bool {
	return i.isDisabled
}

// CurrentStateは、現在のフラグに基づいてウィジェットの総合的な状態を返します。
func (i *Interaction) CurrentState() WidgetState {
	if i.isDisabled {
		return StateDisabled
	}
	if i.isPressed {
		return StatePressed
	}
	if i.isHovered {
		return StateHovered
	}
	return StateNormal
}
