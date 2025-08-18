package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (w *LayoutableWidget) Update() {
	// この基本実装は、具象ウィジェット（Button, Labelなど）で利用されます。
	// Container型は自身のUpdateメソッドでこれをオーバーライドします。
	// ダーティフラグのクリアは、レイアウト計算を行うContainerの責務であるため、
	// ここでは何もしません。
}

func (w *LayoutableWidget) Draw(screen *ebiten.Image) {
	// 非表示のウィジェットは描画しない
	if !w.isVisible {
		return
	}
	// 背景と境界線の描画（フィールドへ直接アクセス）
	DrawStyledBackground(screen, w.x, w.y, w.width, w.height, w.style)
}

func (w *LayoutableWidget) MarkDirty(relayout bool) {
	// すでにダーティで、かつ再レイアウトフラグが既に立っている場合は何もしない
	if w.dirty && (!relayout || w.relayoutDirty) {
		return
	}

	w.dirty = true
	if relayout {
		w.relayoutDirty = true
	}

	// 親が存在し、かつ自身がレイアウト境界でなく、再レイアウトが必要な場合のみ伝播
	if w.parent != nil && !w.relayoutBoundary && relayout {
		// インターフェース経由で直接呼び出す - 型アサーションは不要
		w.parent.MarkDirty(true)
	}
}

func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	if w.relayoutBoundary != isBoundary {
		w.relayoutBoundary = isBoundary
		// 境界設定が変更された場合は、再レイアウトが必要
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) IsDirty() bool {
	return w.dirty
}

func (w *LayoutableWidget) ClearDirty() {
	w.dirty = false
	w.relayoutDirty = false
}

func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.isHovered != hovered {
		w.isHovered = hovered
		// 再描画のみを要求するため、relayoutはfalse
		w.MarkDirty(false)
	}
}

func (w *LayoutableWidget) IsHovered() bool {
	return w.isHovered
}

func (w *LayoutableWidget) SetVisible(visible bool) {
	if w.isVisible != visible {
		w.isVisible = visible
		// 表示状態の変更はレイアウトに影響する可能性があるためtrue
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) IsVisible() bool {
	return w.isVisible
}

// Cleanup は、コンポーネントが不要になったときにリソースを解放するためのメソッドです。
func (w *LayoutableWidget) Cleanup() {
	// イベントハンドラをクリア
	w.eventHandlers = nil
	// 親からの参照を解除
	w.parent = nil
}
