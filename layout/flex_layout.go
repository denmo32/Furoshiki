package layout

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/utils"
)

// FlexLayout は、CSS Flexboxにインスパイアされたレイアウトシステムです。
type FlexLayout struct {
	Direction  Direction
	Justify    Alignment
	AlignItems Alignment
	Wrap       bool
	Gap        int
}

// flexItemInfo は、レイアウト計算中に各子要素の情報を保持するための中間構造体です。
type flexItemInfo struct {
	widget                  component.Widget
	mainSize, crossSize     int
	mainMargin, crossMargin int
	mainMarginStart         int
	flex                    int
}

// Layout は FlexLayout のレイアウトロジックを実装します。
func (l *FlexLayout) Layout(container Container) {
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return
	}

	padding := container.GetPadding()
	containerWidth, containerHeight := container.GetSize()
	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	availableWidth := max(0, containerWidth-padding.Left-padding.Right)
	availableHeight := max(0, containerHeight-padding.Top-padding.Bottom)

	isRow := l.Direction == DirectionRow
	mainSize, crossSize := availableWidth, availableHeight
	if !isRow {
		mainSize, crossSize = availableHeight, availableWidth
	}

	items := collectItemInfo(children, isRow)
	totalFlex, totalBaseMainSize := calculateBaseSizes(items, isRow)
	distributeRemainingSpace(items, mainSize, totalBaseMainSize, totalFlex, l.Gap)
	calculateCrossAxisSizes(items, crossSize, isRow, l.AlignItems)
	applySizes(items, isRow)
	positionItems(items, container, mainSize, crossSize, isRow, l.Justify, l.AlignItems, l.Gap)
}

// collectItemInfo は、子ウィジェットからレイアウト計算に必要な情報を収集します。
func collectItemInfo(children []component.Widget, isRow bool) []flexItemInfo {
	items := make([]flexItemInfo, len(children))
	for i, child := range children {
		s := child.GetStyle()
		margin := style.Insets{}
		if s.Margin != nil {
			margin = *s.Margin
		}

		var mainMargin, crossMargin, mainMarginStart int
		if isRow {
			mainMarginStart = margin.Left
			mainMargin = margin.Left + margin.Right
			crossMargin = margin.Top + margin.Bottom
		} else {
			mainMarginStart = margin.Top
			mainMargin = margin.Top + margin.Bottom
			crossMargin = margin.Left + margin.Right
		}

		items[i] = flexItemInfo{
			widget:          child,
			flex:            child.GetFlex(),
			mainMargin:      mainMargin,
			crossMargin:     crossMargin,
			mainMarginStart: mainMarginStart,
		}
	}
	return items
}

// calculateBaseSizes は、各アイテムの基本サイズを決定し、合計flex値と合計基本サイズを返します。
func calculateBaseSizes(items []flexItemInfo, isRow bool) (float64, int) {
	var totalFlex float64
	var totalBaseMainSize int
	for i := range items {
		item := &items[i]
		if item.flex > 0 {
			totalFlex += float64(item.flex)
		}

		w, h := item.widget.GetSize()
		minW, minH := item.widget.GetMinSize()

		if isRow {
			if item.flex > 0 {
				item.mainSize = minW
			} else {
				item.mainSize = max(utils.IfThen(w <= 0, minW, w), minW)
			}
		} else {
			if item.flex > 0 {
				item.mainSize = minH
			} else {
				item.mainSize = max(utils.IfThen(h <= 0, minH, h), minH)
			}
		}
		totalBaseMainSize += item.mainSize + item.mainMargin
	}
	return totalFlex, totalBaseMainSize
}

// distributeRemainingSpace は、残りの空間をflexアイテムに分配します。
func distributeRemainingSpace(items []flexItemInfo, mainSize, totalBaseMainSize int, totalFlex float64, gap int) {
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * gap
	}

	remainingSpace := mainSize - totalBaseMainSize - totalGap

	if totalFlex > 0 && remainingSpace > 0 {
		sizePerFlex := float64(remainingSpace) / totalFlex
		for i := range items {
			if items[i].flex > 0 {
				items[i].mainSize += int(sizePerFlex * float64(items[i].flex))
			}
		}
	}
}

// calculateCrossAxisSizes は、交差軸のサイズを計算します。
func calculateCrossAxisSizes(items []flexItemInfo, crossSize int, isRow bool, alignItems Alignment) {
	for i := range items {
		item := &items[i]

		// デフォルトでは、子は親の利用可能な交差軸スペース全体を占有します (AlignStretch)。
		item.crossSize = crossSize - item.crossMargin

		// AlignStretchでない場合、子は自身のコンテンツに合わせたサイズになることができます。
		if alignItems != AlignStretch {
			w, h := item.widget.GetSize()
			minW, minH := item.widget.GetMinSize()

			var intrinsicCrossSize int
			if isRow {
				intrinsicCrossSize = max(utils.IfThen(h <= 0, minH, h), minH)
			} else {
				intrinsicCrossSize = max(utils.IfThen(w <= 0, minW, w), minW)
			}

			if intrinsicCrossSize > 0 && intrinsicCrossSize < item.crossSize {
				item.crossSize = intrinsicCrossSize
			}
		}
	}
}

// applySizes は、計算されたサイズを各ウィジェットに設定します。
func applySizes(items []flexItemInfo, isRow bool) {
	for _, item := range items {
		if isRow {
			item.widget.SetSize(item.mainSize, item.crossSize)
		} else {
			item.widget.SetSize(item.crossSize, item.mainSize)
		}
	}
}

// positionItems は、主軸と交差軸の揃え位置に基づいて各ウィジェットを配置します。
func positionItems(items []flexItemInfo, container Container, mainSize, crossSize int, isRow bool, justify, alignItems Alignment, gap int) {
	var currentTotalMainSize int
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * gap
	}

	for _, item := range items {
		currentTotalMainSize += item.mainSize + item.mainMargin
	}
	currentTotalMainSize += totalGap

	freeSpace := mainSize - currentTotalMainSize
	mainOffset := 0
	if freeSpace > 0 {
		switch justify {
		case AlignCenter:
			mainOffset = freeSpace / 2
		case AlignEnd:
			mainOffset = freeSpace
		}
	}

	padding := container.GetPadding()
	containerX, containerY := container.GetPosition()

	mainStart := utils.IfThen(isRow, padding.Left, padding.Top)
	crossStart := utils.IfThen(isRow, padding.Top, padding.Left)
	currentMain := mainStart + mainOffset

	for _, item := range items {
		currentMain += item.mainMarginStart

		crossOffset := 0
		availableCrossSpace := crossSize - item.crossSize - item.crossMargin
		if availableCrossSpace > 0 {
			switch alignItems {
			case AlignCenter:
				crossOffset = availableCrossSpace / 2
			case AlignEnd:
				crossOffset = availableCrossSpace
			}
		}

		finalCrossPos := crossStart + crossOffset

		if isRow {
			item.widget.SetPosition(containerX+currentMain, containerY+finalCrossPos)
		} else {
			item.widget.SetPosition(containerX+finalCrossPos, containerY+currentMain)
		}

		currentMain += item.mainSize + (item.mainMargin - item.mainMarginStart) + gap
	}
}
