package component

import (
	"furoshiki/style"
)

// max は2つのintのうち大きい方を返します。
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SetPosition はウィジェットの絶対座標を設定します。
// このメソッドはレイアウトシステムによって呼び出されるため、ここで初めてウィジェットが
// 「レイアウト済み」であるとマークします。
func (w *LayoutableWidget) SetPosition(x, y int) {
	if !w.state.hasBeenLaidOut {
		w.state.hasBeenLaidOut = true
	}

	if w.position.x != x || w.position.y != y {
		w.position.x = x
		w.position.y = y
		w.MarkDirty(false) // 位置変更は再描画のみが必要
	}
}

// GetPosition はウィジェットの絶対座標を返します。
func (w *LayoutableWidget) GetPosition() (x, y int) {
	return w.position.x, w.position.y
}

// SetSize はウィジェットのサイズを設定します。
func (w *LayoutableWidget) SetSize(width, height int) {
	if width < 0 || height < 0 {
		return
	}

	if w.size.width != width || w.size.height != height {
		w.size.width = width
		w.size.height = height
		w.MarkDirty(true) // サイズ変更は再レイアウトが必要
	}
}

// GetSize はウィジェットのサイズを返します。
func (w *LayoutableWidget) GetSize() (width, height int) {
	return w.size.width, w.size.height
}

// SetMinSize はウィジェットの最小サイズを設定します。
func (w *LayoutableWidget) SetMinSize(width, height int) {
	if width < 0 || height < 0 {
		return
	}

	if w.minSize.width != width || w.minSize.height != height {
		w.minSize.width = width
		w.minSize.height = height
		w.MarkDirty(true) // 最小サイズ変更は再レイアウトが必要
	}
}

// GetMinSize はウィジェットの最小サイズを返します。
// ユーザーが明示的に設定した最小サイズと、ウィジェットのコンテンツ（テキストなど）から
// 計算される最小サイズのうち、大きい方を返します。
func (w *LayoutableWidget) GetMinSize() (width, height int) {
	userMinWidth, userMinHeight := w.minSize.width, w.minSize.height

	if w.contentMinSizeFunc != nil {
		contentMinWidth, contentMinHeight := w.contentMinSizeFunc()
		finalMinWidth := max(contentMinWidth, userMinWidth)
		finalMinHeight := max(contentMinHeight, userMinHeight)
		return finalMinWidth, finalMinHeight
	}

	return userMinWidth, userMinHeight
}

// SetRequestedPosition は、レイアウトに対する希望の相対位置を設定します。
// このメソッドは、親コンテナが `AbsoluteLayout` (主に `ui.ZStack` で作成) を
// 使用している場合にのみ有効です。
func (w *LayoutableWidget) SetRequestedPosition(x, y int) {
	if w.requestedPos.x != x || w.requestedPos.y != y {
		w.requestedPos.x = x
		w.requestedPos.y = y
		w.MarkDirty(true)
	}
}

// GetRequestedPosition は、レイアウトに対する希望の相対位置を返します。
func (w *LayoutableWidget) GetRequestedPosition() (int, int) {
	return w.requestedPos.x, w.requestedPos.y
}

// SetStyle はウィジェットのスタイルを設定します。
func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更はパディングやフォントサイズに影響し、レイアウトが変わる可能性があるため、
	// 安全策として再レイアウトを要求します。
	w.MarkDirty(true)
}

// GetStyle はウィジェットの現在のスタイルの安全なコピーを返します。
// 注意: このメソッドは呼び出されるたびにスタイルのディープコピーを生成します。
// パフォーマンスが重要な場面で頻繁に呼び出すとオーバーヘッドになる可能性があります。
// しかし、これにより外部からの意図しないスタイルの変更を防ぎ、安全性を保証します。
func (w *LayoutableWidget) GetStyle() style.Style {
	return w.style.DeepCopy()
}

// SetFlex はFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
	}
	if w.layout.flex != flex {
		w.layout.flex = flex
		w.MarkDirty(true)
	}
}

// GetFlex はウィジェットの伸縮係数を返します。
func (w *LayoutableWidget) GetFlex() int {
	return w.layout.flex
}

// SetLayoutBoundary は、このウィジェットをレイアウト計算の境界とするか設定します。
// メソッド名を SetRelayoutBoundary から SetLayoutBoundary に変更して直感性を向上させました。
func (w *LayoutableWidget) SetLayoutBoundary(isBoundary bool) {
	if w.layout.relayoutBoundary != isBoundary {
		w.layout.relayoutBoundary = isBoundary
		w.MarkDirty(true)
	}
}



// SetParent はウィジェットの親コンテナを設定します。
func (w *LayoutableWidget) SetParent(parent Container) {
	w.hierarchy.parent = parent
}

// GetParent はウィジェットの親コンテナを返します。
func (w *LayoutableWidget) GetParent() Container {
	return w.hierarchy.parent
}

// SetLayoutData はウィジェットにレイアウト固有のデータを設定します。
func (w *LayoutableWidget) SetLayoutData(data any) {
	w.layout.layoutData = data
	w.MarkDirty(true)
}

// GetLayoutData はウィジェットからレイアウト固有のデータを取得します。
func (w *LayoutableWidget) GetLayoutData() any {
	return w.layout.layoutData
}