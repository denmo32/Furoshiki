package core

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

func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更は必ずしも再レイアウトを必要としないかもしれないが、
	// Paddingなどが変わる可能性があるため、安全策としてtrueにする
	w.MarkDirty(true)
}

func (w *LayoutableWidget) GetStyle() *style.Style {
	return &w.style
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
