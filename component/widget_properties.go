package component

import "furoshiki/style"

// SetPosition はウィジェットの絶対座標を設定します。
// 座標が変更された場合、再描画を要求します（レイアウトの再計算は不要）。
func (w *LayoutableWidget) SetPosition(x, y int) {
	if w.x != x || w.y != y {
		w.x = x
		w.y = y
		w.MarkDirty(false) // 位置変更は再描画のみが必要
	}
}

func (w *LayoutableWidget) GetPosition() (x, y int) {
	return w.x, w.y
}

// SetSize はウィジェットのサイズを設定します。
// サイズが変更された場合、親コンテナにレイアウトの再計算を要求します。
func (w *LayoutableWidget) SetSize(width, height int) {
	// サイズが負の値の場合は処理しません。
	if width < 0 || height < 0 {
		return
	}

	if w.width != width || w.height != height {
		w.width = width
		w.height = height
		w.MarkDirty(true) // サイズ変更は再レイアウトが必要
	}
}

func (w *LayoutableWidget) GetSize() (width, height int) {
	return w.width, w.height
}

// SetMinSize はウィジェットの最小サイズを設定します。
// 最小サイズが変更された場合、親コンテナにレイアウトの再計算を要求します。
func (w *LayoutableWidget) SetMinSize(width, height int) {
	// 最小サイズが負の値の場合は処理しません。
	if width < 0 || height < 0 {
		return
	}

	if w.minWidth != width || w.minHeight != height {
		w.minWidth = width
		w.minHeight = height
		w.MarkDirty(true) // 最小サイズ変更は再レイアウトが必要
	}
}

func (w *LayoutableWidget) GetMinSize() (width, height int) {
	return w.minWidth, w.minHeight
}

// SetRequestedPosition は、レイアウトに対する希望の相対位置を設定します。
// このメソッドは `AbsoluteLayout` のような特定のレイアウトと協調するために存在します。
// `FlexLayout` など他のレイアウトでは効果がありません。
func (w *LayoutableWidget) SetRequestedPosition(x, y int) {
	if w.requestedX != x || w.requestedY != y {
		w.requestedX = x
		w.requestedY = y
		// 希望位置の変更は再レイアウトをトリガーすべきです。
		w.MarkDirty(true)
	}
}

// GetRequestedPosition は、レイアウトに対する希望の相対位置を取得します。
func (w *LayoutableWidget) GetRequestedPosition() (int, int) {
	return w.requestedX, w.requestedY
}


// SetStyle はウィジェットのスタイルを設定します。
// スタイルの変更はレイアウトに影響する可能性があるため、再レイアウトを要求します。
func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更はパディングやフォントサイズに影響し、レイアウトが変わる可能性があるため、
	// 安全策として再レイアウトを要求します。
	w.MarkDirty(true)
}

// GetStyle はウィジェットの現在のスタイルの安全なコピーを返します。
// 内部の値もコピーされるため、この戻り値を変更しても元のウィジェットには影響しません。
func (w *LayoutableWidget) GetStyle() style.Style {
	return w.style.DeepCopy()
}

// SetFlex はFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0 // flex値は0以上である必要があります。
	}
	if w.flex != flex {
		w.flex = flex
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) GetFlex() int {
	return w.flex
}

func (w *LayoutableWidget) SetParent(parent Container) {
	w.parent = parent
}

func (w *LayoutableWidget) GetParent() Container {
	return w.parent
}