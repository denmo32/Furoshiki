package layout

import (
	"furoshiki/component"
	"furoshiki/style"
)

type Layout interface {
	Layout(container Container)
}

type Container interface {
	GetSize() (width, height int)
	GetPosition() (x, y int)
	GetChildren() []component.Widget
	GetStyle() *style.Style
}

type LayoutType int

const (
	LayoutAbsolute LayoutType = iota
	LayoutFlex
)

type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)

type AbsoluteLayout struct{}

func (l *AbsoluteLayout) Layout(container Container) {
	containerX, containerY := container.GetPosition()
	for _, child := range container.GetChildren() {
		// 非表示のコンポーネントはレイアウト処理をスキップ
		if !child.IsVisible() {
			continue
		}
		childX, childY := child.GetPosition()
		child.SetPosition(containerX+childX, containerY+childY)
	}
}

// --- FlexLayout (改修版) ---
type FlexLayout struct {
	Direction  Direction
	Justify    Alignment
	AlignItems Alignment
	Gap        int
	Wrap       bool
}

// max は2つのintの大きい方を返します。
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// flexItemInfo は、レイアウト計算中に各子要素の情報を保持するための中間構造体です。
// これにより、ループの回数を減らし、計算を効率化します。
type flexItemInfo struct {
	widget                  component.Widget
	mainSize, crossSize     int // 最終的な主軸・交差軸サイズ
	mainMargin, crossMargin int // 主軸・交差軸のマージン合計
	mainMarginStart         int // 主軸の開始側マージン
}

func (l *FlexLayout) Layout(container Container) {
	// --- ステップ 1: 初期設定と可視コンポーネントのフィルタリング ---
	allChildren := container.GetChildren()
	children := make([]component.Widget, 0, len(allChildren))
	for _, child := range allChildren {
		if child.IsVisible() {
			children = append(children, child)
		}
	}

	if len(children) == 0 {
		return
	}

	containerStyle := container.GetStyle()
	containerWidth, containerHeight := container.GetSize()
	containerX, containerY := container.GetPosition()

	availableWidth := containerWidth - containerStyle.Padding.Left - containerStyle.Padding.Right
	availableHeight := containerHeight - containerStyle.Padding.Top - containerStyle.Padding.Bottom

	isRow := l.Direction == DirectionRow
	mainSize, crossSize := availableWidth, availableHeight
	if !isRow {
		mainSize, crossSize = availableHeight, availableWidth
	}

	// --- ステップ 2: レイアウト情報の収集と初期サイズ計算 (1回目のループ) ---
	items := make([]flexItemInfo, len(children))
	var totalFixedMainSize int
	var totalFlex float64

	for i, child := range children {
		style := child.GetStyle()
		minW, minH := child.GetMinSize()
		w, h := child.GetSize()

		var info flexItemInfo
		info.widget = child

		// 主軸と交差軸のマージンを計算
		if isRow {
			info.mainMarginStart = style.Margin.Left
			info.mainMargin = style.Margin.Left + style.Margin.Right
			info.crossMargin = style.Margin.Top + style.Margin.Bottom
		} else {
			info.mainMarginStart = style.Margin.Top
			info.mainMargin = style.Margin.Top + style.Margin.Bottom
			info.crossMargin = style.Margin.Left + style.Margin.Right
		}

		// 主軸の初期サイズを決定 (flexでない要素)
		if child.GetFlex() > 0 {
			totalFlex += float64(child.GetFlex())
			// flexアイテムでも最小サイズは固定サイズとして確保
			if isRow {
				info.mainSize = max(0, minW)
			} else {
				info.mainSize = max(0, minH)
			}
		} else {
			if isRow {
				info.mainSize = max(w, minW)
			} else {
				info.mainSize = max(h, minH)
			}
		}
		totalFixedMainSize += info.mainSize + info.mainMargin
		items[i] = info
	}

	// --- ステップ 3: 伸縮可能スペースの計算と分配 ---
	totalGap := 0
	if len(children) > 1 {
		totalGap = (len(children) - 1) * l.Gap
	}
	remainingSpace := mainSize - totalFixedMainSize - totalGap

	if totalFlex > 0 && remainingSpace > 0 {
		sizePerFlex := float64(remainingSpace) / totalFlex
		for i := range items {
			if items[i].widget.GetFlex() > 0 {
				flexSize := int(sizePerFlex * float64(items[i].widget.GetFlex()))
				items[i].mainSize += flexSize
			}
		}
	}

	// --- ステップ 4: 交差軸のサイズ計算 (AlignStretch) ---
	if l.AlignItems == AlignStretch {
		for i := range items {
			minW, minH := items[i].widget.GetMinSize()
			if isRow {
				items[i].crossSize = max(minH, crossSize-items[i].crossMargin)
			} else {
				items[i].crossSize = max(minW, crossSize-items[i].crossMargin)
			}
		}
	} else {
		// Stretchでない場合は、ウィジェット本来のサイズを維持
		for i := range items {
			w, h := items[i].widget.GetSize()
			minW, minH := items[i].widget.GetMinSize()
			if isRow {
				items[i].crossSize = max(h, minH)
			} else {
				items[i].crossSize = max(w, minW)
			}
		}
	}

	// --- ステップ 5: 最終的なサイズをウィジェットに設定 ---
	for _, item := range items {
		if isRow {
			item.widget.SetSize(item.mainSize, item.crossSize)
		} else {
			item.widget.SetSize(item.crossSize, item.mainSize)
		}
	}

	// --- ステップ 6: 位置計算と配置 (2回目のループ) ---
	var currentTotalMainSize int
	for _, item := range items {
		currentTotalMainSize += item.mainSize + item.mainMargin
	}
	currentTotalMainSize += totalGap

	freeSpace := mainSize - currentTotalMainSize
	mainOffset := 0
	switch l.Justify {
	case AlignCenter:
		mainOffset = freeSpace / 2
	case AlignEnd:
		mainOffset = freeSpace
	}
	if mainOffset < 0 {
		mainOffset = 0
	}

	mainStart := containerStyle.Padding.Left
	crossStart := containerStyle.Padding.Top
	if !isRow {
		mainStart = containerStyle.Padding.Top
		crossStart = containerStyle.Padding.Left
	}
	currentMain := mainStart + mainOffset

	for _, item := range items {
		currentMain += item.mainMarginStart

		crossOffset := 0
		switch l.AlignItems {
		case AlignCenter:
			crossOffset = (crossSize - item.crossSize - item.crossMargin) / 2
		case AlignEnd:
			crossOffset = crossSize - item.crossSize - item.crossMargin + item.mainMarginStart
		}

		finalCrossPos := crossStart + crossOffset

		if isRow {
			item.widget.SetPosition(containerX+currentMain, containerY+finalCrossPos)
		} else {
			item.widget.SetPosition(containerX+finalCrossPos, containerY+currentMain)
		}

		currentMain += item.mainSize + (item.mainMargin - item.mainMarginStart) + l.Gap
	}
}
