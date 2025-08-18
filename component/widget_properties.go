package component

import "furoshiki/style"

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

func (w *LayoutableWidget) SetSize(width, height int) {
	// サイズが負の値の場合は処理しない
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

func (w *LayoutableWidget) SetMinSize(width, height int) {
	// 最小サイズが負の値の場合は処理しない
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

// [追加] AbsoluteLayoutのために、要求された相対位置を設定・取得するメソッドを追加します。
// これらはWidgetインターフェースには含まれず、特定のレイアウト(AbsoluteLayout)と
// ウィジェットが協調するために使用されます。

// SetRequestedPosition は、レイアウトに対する希望の相対位置を設定します。
func (w *LayoutableWidget) SetRequestedPosition(x, y int) {
	if w.requestedX != x || w.requestedY != y {
		w.requestedX = x
		w.requestedY = y
		// 希望位置の変更は再レイアウトをトリガーすべき
		w.MarkDirty(true)
	}
}

// GetRequestedPosition は、レイアウトに対する希望の相対位置を取得します。
func (w *LayoutableWidget) GetRequestedPosition() (int, int) {
	return w.requestedX, w.requestedY
}


func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更はパディングやマージンに影響し、レイアウトが変わる可能性があるため、
	// 安全策として再レイアウトを要求します。
	w.MarkDirty(true)
}

// GetStyle はウィジェットの現在のスタイルを返します。
// スタイルのコピーを返すため、この戻り値を変更してもウィジェットには影響しません。
// スタイルを変更するには SetStyle を使用してください。
func (w *LayoutableWidget) GetStyle() style.Style {
	return w.style
}

func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
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