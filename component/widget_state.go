package component

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	// 非表示のウィジェットは描画しません。
	if !w.isVisible {
		return
	}
	// 背景と境界線の描画
	DrawStyledBackground(screen, w.x, w.y, w.width, w.height, w.style)
}

// MarkDirty はウィジェットの状態が変更されたことをマークします。
// relayoutがtrueの場合、親コンテナにも再レイアウトが必要であることを伝播させます。
func (w *LayoutableWidget) MarkDirty(relayout bool) {
	// すでにダーティで、かつ再レイアウト要求が既に立っている場合は何もしません。
	if w.dirty && (!relayout || w.relayoutDirty) {
		return
	}

	w.dirty = true
	if relayout {
		w.relayoutDirty = true
	}

	// 親が存在し、かつ自身がレイアウト境界でなく、再レイアウトが必要な場合のみ伝播します。
	if w.parent != nil && !w.relayoutBoundary && relayout {
		// 親コンテナに再レイアウトが必要であることを再帰的に伝播させます。
		w.parent.MarkDirty(true)
	}
}

// SetRelayoutBoundary は、このウィジェットをレイアウト計算の境界とするか設定します。
func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	if w.relayoutBoundary != isBoundary {
		w.relayoutBoundary = isBoundary
		// 境界設定が変更された場合は、再レイアウトが必要です。
		w.MarkDirty(true)
	}
}

// IsDirty はウィジェットが再描画または再レイアウトを必要とするかどうかを返します。
// レイアウトコンテナはこのフラグを見て、レイアウトの再計算を行うかを判断します。
func (w *LayoutableWidget) IsDirty() bool {
	return w.dirty
}

// ClearDirty はダーティフラグと再レイアウトダーティフラグの両方をクリアします。
func (w *LayoutableWidget) ClearDirty() {
	w.dirty = false
	w.relayoutDirty = false
}

// SetHovered はホバー状態を設定し、再描画を要求します。
func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.isHovered != hovered {
		w.isHovered = hovered
		// ホバー状態の変更は見た目にのみ影響するため、再描画のみを要求します（relayoutはfalse）。
		w.MarkDirty(false)
	}
}

func (w *LayoutableWidget) IsHovered() bool {
	return w.isHovered
}

// SetDisabled はウィジェットの有効・無効状態を設定します。
func (w *LayoutableWidget) SetDisabled(disabled bool) {
	if w.isDisabled != disabled {
		w.isDisabled = disabled
		// 無効状態の変更は見た目に影響するため、再描画を要求します。
		w.MarkDirty(false)
	}
}

// IsDisabled はウィジェットが無効状態であるかを返します。
func (w *LayoutableWidget) IsDisabled() bool {
	return w.isDisabled
}

// SetVisible はウィジェットの可視性を設定します。
func (w *LayoutableWidget) SetVisible(visible bool) {
	if w.isVisible != visible {
		w.isVisible = visible
		// 表示状態の変更はレイアウトに影響するため、再レイアウトを要求します。
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) IsVisible() bool {
	return w.isVisible
}

// Cleanup は、コンポーネントが不要になったときにリソースを解放するためのメソッドです。
// UIツリーからウィジェットが削除される際などに呼び出されるべきです。
func (w *LayoutableWidget) Cleanup() {
	// イベントハンドラへの参照をクリアし、ガベージコレクションの対象とします。
	w.eventHandlers = nil
	// 親からの参照を解除します。
	w.parent = nil
}