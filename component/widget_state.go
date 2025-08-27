package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// WidgetState は、ウィジェットが取りうるインタラクティブな状態を定義します。
type WidgetState int

const (
	// StateNormal は、ユーザー入力がないデフォルトの状態です。
	StateNormal WidgetState = iota
	// StateHovered は、マウスカーソルがウィジェット上にある状態です。
	StateHovered
	// StatePressed は、ウィジェットがクリックまたはタップされている最中の状態です。
	StatePressed
	// StateDisabled は、ウィジェットが無効化され、ユーザー入力を受け付けない状態です。
	StateDisabled
)

// Update はウィジェットの状態を更新します。
// この基本実装は、具象ウィジェット（Button, Labelなど）で利用されます。
// Container型は自身のUpdateメソッドでこれをオーバーライドして、子の更新やレイアウト処理を行います。
func (w *LayoutableWidget) Update() {
	// ダーティフラグのクリアは、レイアウト計算を行う親コンテナの責務であるため、
	// この基本実装では何もしません。
}

// Draw はウィジェットの背景と境界線を描画します。
// テキストなどを持つウィジェットは、このメソッドをオーバーライドして追加の描画処理を行います。
func (w *LayoutableWidget) Draw(screen *ebiten.Image) {
	if !w.state.isVisible || !w.state.hasBeenLaidOut {
		return
	}
	// NOTE: [FIX] 削除された w.style の代わりに、StyleManagerから現在の状態に合ったスタイルを取得します。
	styleToUse := w.StyleManager.GetStyleForState(w.CurrentState())
	DrawStyledBackground(screen, w.position.x, w.position.y, w.size.width, w.size.height, styleToUse)
}

// MarkDirty はウィジェットの状態が変更されたことをマークします。
// relayoutがtrueの場合、より高いダーティレベル(levelRelayoutDirty)が設定され、
// 親コンテナにも再レイアウトが必要であることが伝播されます。
func (w *LayoutableWidget) MarkDirty(relayout bool) {
	requestedLevel := levelRedrawDirty
	if relayout {
		requestedLevel = levelRelayoutDirty
	}

	// 現在のダーティレベルが要求されたレベルより低い場合のみ更新します。
	// これにより、levelRelayoutDirtyが設定されているウィジェットにlevelRedrawDirtyを要求しても、
	// ダーティレベルが意図せず下がることを防ぎます。
	if w.state.dirtyLevel >= requestedLevel {
		return
	}

	w.state.dirtyLevel = requestedLevel

	// 親が存在し、かつ自身がレイアウト境界でなく、再レイアウトが必要な場合のみ伝播します。
	if w.hierarchy.parent != nil && !w.layout.relayoutBoundary && relayout {
		w.hierarchy.parent.MarkDirty(true)
	}
}

// IsDirty はウィジェットが再描画または再レイアウトを必要とするかどうかを返します。
func (w *LayoutableWidget) IsDirty() bool {
	return w.state.dirtyLevel > levelClean
}

// NeedsRelayout はウィジェットがレイアウトの再計算を必要とするかどうかを返します。
func (w *LayoutableWidget) NeedsRelayout() bool {
	return w.state.dirtyLevel == levelRelayoutDirty
}

// ClearDirty はダーティレベルをlevelCleanにリセットします。
func (w *LayoutableWidget) ClearDirty() {
	w.state.dirtyLevel = levelClean
}

// SetHovered はホバー状態を設定し、再描画を要求します。
func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.state.isHovered != hovered {
		w.state.isHovered = hovered
		w.MarkDirty(false) // ホバー状態の変更は再描画のみ
	}
}

// IsHovered はウィジェットがホバー状態かどうかを返します。
func (w *LayoutableWidget) IsHovered() bool {
	return w.state.isHovered
}

// SetPressed は押下状態を設定し、再描画を要求します。
func (w *LayoutableWidget) SetPressed(pressed bool) {
	if w.state.isPressed != pressed {
		w.state.isPressed = pressed
		w.MarkDirty(false) // 押下状態の変更は再描画のみ
	}
}

// IsPressed はウィジェットが押下状態かどうかを返します。
func (w *LayoutableWidget) IsPressed() bool {
	return w.state.isPressed
}

// CurrentState はウィジェットの現在のインタラクティブな状態を返します。
func (w *LayoutableWidget) CurrentState() WidgetState {
	if w.state.isDisabled {
		return StateDisabled
	}
	if w.state.isPressed {
		return StatePressed
	}
	if w.state.isHovered {
		return StateHovered
	}
	return StateNormal
}

// SetDisabled はウィジェットの有効・無効状態を設定します。
func (w *LayoutableWidget) SetDisabled(disabled bool) {
	if w.state.isDisabled != disabled {
		w.state.isDisabled = disabled
		w.MarkDirty(false)
	}
}

// IsDisabled はウィジェットが無効状態かどうかを返します。
func (w *LayoutableWidget) IsDisabled() bool {
	return w.state.isDisabled
}

// SetVisible はウィジェットの可視性を設定します。
func (w *LayoutableWidget) SetVisible(visible bool) {
	if w.state.isVisible != visible {
		w.state.isVisible = visible
		w.MarkDirty(true) // 表示状態の変更はレイアウトに影響
	}
}

// IsVisible はウィジェットが可視状態かどうかを返します。
func (w *LayoutableWidget) IsVisible() bool {
	return w.state.isVisible
}

// HasBeenLaidOut は、ウィジェットが少なくとも一度レイアウト計算されたかを返します。
func (w *LayoutableWidget) HasBeenLaidOut() bool {
	return w.state.hasBeenLaidOut
}

// Cleanup は、コンポーネントが不要になったときにリソースを解放するためのメソッドです。
func (w *LayoutableWidget) Cleanup() {
	w.eventHandlers = nil
	w.hierarchy.parent = nil
}