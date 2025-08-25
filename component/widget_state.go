package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// 【改善】WidgetStateの定義を、それを利用するロジック（CurrentStateなど）と
// 同じファイルに統合しました。これにより、状態の定義と実装がまとまり、
// コードの関連性が明確になって可読性が向上します。
// （以前は `component/widget_state_types.go` にありました）

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
	// 非表示、または未レイアウトのウィジェットは描画しません。
	if !w.state.isVisible || !w.state.hasBeenLaidOut {
		return
	}
	// 背景と境界線の描画
	DrawStyledBackground(screen, w.position.x, w.position.y, w.size.width, w.size.height, w.style)
}

// MarkDirty はウィジェットの状態が変更されたことをマークします。
// relayoutがtrueの場合、より高いダーティレベル(levelRelayoutDirty)が設定され、
// 親コンテナにも再レイアウトが必要であることが伝播されます。
func (w *LayoutableWidget) MarkDirty(relayout bool) {
	// 要求されたダーティレベルを決定します。
	requestedLevel := levelRedrawDirty
	if relayout {
		requestedLevel = levelRelayoutDirty
	}

	// 現在のダーティレベルが要求されたレベル以上であれば、何もする必要はありません。
	// これにより、levelRelayoutDirtyが設定されているウィジェットにlevelRedrawDirtyを要求しても、
	// ダーティレベルが下がることはありません。
	if w.state.dirtyLevel >= requestedLevel {
		return
	}

	// ダーティレベルを新しいレベルに更新します。
	w.state.dirtyLevel = requestedLevel

	// 親が存在し、かつ自身がレイアウト境界でなく、再レイアウトが必要な場合のみ伝播します。
	if w.hierarchy.parent != nil && !w.layout.relayoutBoundary && relayout {
		// 親コンテナに再レイアウトが必要であることを再帰的に伝播させます。
		w.hierarchy.parent.MarkDirty(true)
	}
}

// SetRelayoutBoundary は、このウィジェットをレイアウト計算の境界とするか設定します。
func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	if w.layout.relayoutBoundary != isBoundary {
		w.layout.relayoutBoundary = isBoundary
		// 境界設定が変更された場合は、再レイアウトが必要です。
		w.MarkDirty(true)
	}
}

// IsDirty はウィジェットが再描画または再レイアウトを必要とするかどうかを返します。
// dirtyLevelがlevelCleanでなければ、何らかの更新が必要であると判断されます。
func (w *LayoutableWidget) IsDirty() bool {
	return w.state.dirtyLevel > levelClean
}

// NeedsRelayout はウィジェットがレイアウトの再計算を必要とするかどうかを返します。
// dirtyLevelが最高のlevelRelayoutDirtyである場合のみtrueを返します。
func (w *LayoutableWidget) NeedsRelayout() bool {
	return w.state.dirtyLevel == levelRelayoutDirty
}

// ClearDirty はダーティレベルをlevelCleanにリセットします。
// レイアウトと描画が完了した後にコンテナから呼び出されます。
func (w *LayoutableWidget) ClearDirty() {
	w.state.dirtyLevel = levelClean
}

// SetHovered はホバー状態を設定し、再描画を要求します。
func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.state.isHovered != hovered {
		w.state.isHovered = hovered
		// ホバー状態の変更は見た目にのみ影響するため、再描画のみを要求します（relayoutはfalse）。
		w.MarkDirty(false)
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
		// 押下状態の変更は見た目にのみ影響するため、再描画のみを要求します（relayoutはfalse）。
		w.MarkDirty(false)
	}
}

// IsPressed はウィジェットが押下状態かどうかを返します。
func (w *LayoutableWidget) IsPressed() bool {
	return w.state.isPressed
}

// CurrentState はウィジェットの現在のインタラクティブな状態を返します。
// このメソッドは InteractiveState インターフェースを実装します。
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
		// 無効状態の変更は見た目に影響するため、再描画を要求します。
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
		// 表示状態の変更はレイアウトに影響するため、再レイアウトを要求します。
		w.MarkDirty(true)
	}
}

// IsVisible はウィジェットが可視状態かどうかを返します。
func (w *LayoutableWidget) IsVisible() bool {
	return w.state.isVisible
}

// HasBeenLaidOut は、ウィジェットが少なくとも一度レイアウト計算されたかを返します。
// このメソッドは InteractiveState インターフェースを実装します。
func (w *LayoutableWidget) HasBeenLaidOut() bool {
	return w.state.hasBeenLaidOut
}

// Cleanup は、コンポーネントが不要になったときにリソースを解放するためのメソッドです。
// UIツリーからウィジェットが削除される際などに呼び出されるべきです。
func (w *LayoutableWidget) Cleanup() {
	// イベントハンドラへの参照をクリアし、ガベージコレクションの対象とします。
	w.eventHandlers = nil
	// 親からの参照を解除します。
	w.hierarchy.parent = nil
}
