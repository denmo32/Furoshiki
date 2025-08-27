package layout

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/utils"
	"image"
)

// FlexLayout は、CSS Flexboxにインスパイアされたレイアウトシステムです。
type FlexLayout struct {
	Direction    Direction
	Justify      Alignment
	AlignItems   Alignment
	AlignContent Alignment
	Wrap         bool
	Gap          int
}

// flexItemInfo は、レイアウト計算中に各子要素の情報を保持するための中間構造体です。
type flexItemInfo struct {
	widget                  component.Widget
	mainSize, crossSize     int
	mainMargin, crossMargin int
	mainMarginStart         int
	flex                    int
}

// flexLine は、折り返しレイアウト時に一行（または一列）を表現する内部構造体です。
type flexLine struct {
	items         []*flexItemInfo
	mainAxisSize  int // このラインの主軸方向のサイズ
	crossAxisSize int // このラインの交差軸方向のサイズ
}

// Measure は FlexLayout の要求サイズを計算します。
func (l *FlexLayout) Measure(container Container, availableSize image.Point) image.Point {
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return image.Point{}
	}

	padding := container.GetPadding()
	isRow := l.Direction == DirectionRow

	// 子が利用できるスペース
	availableW := availableSize.X - padding.Left - padding.Right
	availableH := availableSize.Y - padding.Top - padding.Bottom
	if availableW < 0 {
		availableW = 0
	}
	if availableH < 0 {
		availableH = 0
	}

	mainAvailable, crossAvailable := availableW, availableH
	if !isRow {
		mainAvailable, crossAvailable = availableH, availableW
	}

	items := collectItemInfo(children, isRow)
	// NOTE: Measureパスでは、子のサイズを純粋に計測したい。
	// そのため、利用可能な交差軸サイズ(crossAvailable)を渡して、HeightForWiderなどが正しく動作するようにする。
	calculateBaseSizes(items, isRow, crossAvailable, l.AlignItems)

	var desiredMain, desiredCross int

	if l.Wrap {
		lines := l.splitIntoLines(items, mainAvailable)
		maxMain := 0
		totalCross := 0
		for i, line := range lines {
			lineMainSize := 0
			for _, item := range line.items {
				lineMainSize += item.mainSize + item.mainMargin
			}
			if len(line.items) > 1 {
				lineMainSize += (len(line.items) - 1) * l.Gap
			}
			if lineMainSize > maxMain {
				maxMain = lineMainSize
			}
			line.crossAxisSize = calculateLineCrossSize(line.items)
			totalCross += line.crossAxisSize
			if i > 0 {
				totalCross += l.Gap
			}
		}
		desiredMain = maxMain
		desiredCross = totalCross
	} else {
		totalMain := 0
		maxCross := 0
		for _, item := range items {
			totalMain += item.mainSize + item.mainMargin
			if item.crossSize+item.crossMargin > maxCross {
				maxCross = item.crossSize + item.crossMargin
			}
		}
		if len(items) > 1 {
			totalMain += (len(items) - 1) * l.Gap
		}
		desiredMain = totalMain
		desiredCross = maxCross
	}

	var desiredW, desiredH int
	if isRow {
		desiredW, desiredH = desiredMain, desiredCross
	} else {
		desiredW, desiredH = desiredCross, desiredMain
	}

	return image.Point{
		X: desiredW + padding.Left + padding.Right,
		Y: desiredH + padding.Top + padding.Bottom,
	}
}

// Arrange は FlexLayout のレイアウトロジックを実装します。
func (l *FlexLayout) Arrange(container Container, finalBounds image.Rectangle) error {
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return nil
	}

	padding := container.GetPadding()
	isRow := l.Direction == DirectionRow

	availableW := finalBounds.Dx() - padding.Left - padding.Right
	availableH := finalBounds.Dy() - padding.Top - padding.Bottom
	if availableW < 0 {
		availableW = 0
	}
	if availableH < 0 {
		availableH = 0
	}

	mainSize, crossSize := availableW, availableH
	if !isRow {
		mainSize, crossSize = availableH, availableW
	}

	items := collectItemInfo(children, isRow)
	calculateBaseSizes(items, isRow, crossSize, l.AlignItems)

	if l.Wrap {
		l.layoutMultiLine(items, container, mainSize, crossSize, isRow, finalBounds)
	} else {
		l.layoutSingleLine(items, container, mainSize, crossSize, isRow, finalBounds)
	}

	// すべての子のArrangeを再帰的に呼び出す
	for _, item := range items {
		x, y := item.widget.GetPosition()
		w, h := item.widget.GetSize()
		childBounds := image.Rect(x, y, x+w, y+h)
		if err := item.widget.Arrange(childBounds); err != nil {
			return err
		}
	}

	return nil
}

// layoutSingleLine は、折り返しなしのレイアウト計算と配置を実行します。
func (l *FlexLayout) layoutSingleLine(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool, finalBounds image.Rectangle) {
	var totalFlex float64
	var totalBaseMainSize int
	for _, item := range items {
		if item.flex > 0 {
			totalFlex += float64(item.flex)
		}
		totalBaseMainSize += item.mainSize + item.mainMargin
	}

	distributeRemainingSpace(items, mainSize, totalBaseMainSize, totalFlex, l.Gap)
	calculateCrossAxisSizes(items, crossSize, isRow, l.AlignItems)
	applySizes(items, isRow)
	positionItems(items, container, mainSize, crossSize, isRow, l.Justify, l.AlignItems, l.Gap, finalBounds)
}

// layoutMultiLine は、折り返しありのレイアウト計算と配置を実行します。
// [修正] 折り返しレイアウトのロジックを全面的に修正しました。
func (l *FlexLayout) layoutMultiLine(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool, finalBounds image.Rectangle) {
	lines := l.splitIntoLines(items, mainSize)
	if len(lines) == 0 {
		return
	}

	var totalCrossAxisSize int
	// 最初のループ: 各行のサイズを計算します
	for _, line := range lines {
		var lineTotalFlex float64
		var lineTotalBaseMainSize int
		for _, item := range line.items {
			if item.flex > 0 {
				lineTotalFlex += float64(item.flex)
			}
			lineTotalBaseMainSize += item.mainSize + item.mainMargin
		}
		// 1. 主軸方向のサイズを確定させます
		distributeRemainingSpace(line.items, mainSize, lineTotalBaseMainSize, lineTotalFlex, l.Gap)

		// 2. 確定した主軸サイズを元に、各アイテムの交差軸サイズを再計算します。
		//    (テキスト折り返しなどで高さが変わるウィジェットに対応するため)
		for _, item := range line.items {
			measuredSize := item.widget.Measure(image.Point{
				X: utils.IfThen(isRow, item.mainSize, crossSize-item.crossMargin),
				Y: utils.IfThen(isRow, crossSize-item.crossMargin, item.mainSize),
			})
			if isRow {
				item.crossSize = measuredSize.Y
			} else {
				item.crossSize = measuredSize.X
			}
		}

		// 3. アイテムの最終的な交差軸サイズに基づいて、行自体のサイズを計算します
		line.crossAxisSize = calculateLineCrossSize(line.items)
		totalCrossAxisSize += line.crossAxisSize
	}

	if len(lines) > 1 {
		totalCrossAxisSize += (len(lines) - 1) * l.Gap
	}

	// 4. AlignContent に基づいて、各行の配置開始位置を計算します
	padding := container.GetPadding()
	crossStart := utils.IfThen(isRow, finalBounds.Min.Y+padding.Top, finalBounds.Min.X+padding.Left)
	currentCross := float64(crossStart)

	freeCrossSpace := crossSize - totalCrossAxisSize
	if freeCrossSpace > 0 {
		switch l.AlignContent {
		case AlignCenter:
			currentCross += float64(freeCrossSpace) / 2.0
		case AlignEnd:
			currentCross += float64(freeCrossSpace)
		}
	}

	// 2番目のループ: 計算した行サイズに基づいてアイテムを配置します
	for _, line := range lines {
		// AlignItems == AlignStretch の場合、アイテムを行の高さまで引き伸ばします
		if l.AlignItems == AlignStretch {
			for _, item := range line.items {
				item.crossSize = line.crossAxisSize - item.crossMargin
				if item.crossSize < 0 {
					item.crossSize = 0
				}
			}
		}

		// positionItems は、確定したサイズと AlignItems に基づいてオフセットを計算し、位置を設定します
		positionItems(line.items, container, mainSize, line.crossAxisSize, isRow, l.Justify, l.AlignItems, l.Gap, finalBounds, int(currentCross))
		currentCross += float64(line.crossAxisSize + l.Gap)
	}

	// 5. すべてのアイテムに最終的なサイズを適用します
	applySizes(items, isRow)
}

// splitIntoLines は、アイテムをコンテナの主軸サイズに基づいて複数のラインに分割します。
func (l *FlexLayout) splitIntoLines(items []*flexItemInfo, mainSize int) []*flexLine {
	var lines []*flexLine
	if len(items) == 0 {
		return lines
	}

	currentLineItems := make([]*flexItemInfo, 0)
	currentMainSize := 0

	for _, item := range items {
		itemMainSize := item.mainSize + item.mainMargin
		gap := utils.IfThen(len(currentLineItems) > 0, l.Gap, 0)

		if len(currentLineItems) > 0 && currentMainSize+itemMainSize+gap > mainSize {
			lines = append(lines, &flexLine{items: currentLineItems})
			currentLineItems = make([]*flexItemInfo, 0)
			currentMainSize = 0
			gap = 0
		}

		currentLineItems = append(currentLineItems, item)
		currentMainSize += itemMainSize + gap
	}

	if len(currentLineItems) > 0 {
		lines = append(lines, &flexLine{items: currentLineItems})
	}

	return lines
}

// calculateLineCrossSize は、ライン内のアイテムに基づいてライン自体の交差軸サイズを決定します。
func calculateLineCrossSize(lineItems []*flexItemInfo) int {
	maxCrossSize := 0
	for _, item := range lineItems {
		size := item.crossSize + item.crossMargin
		if size > maxCrossSize {
			maxCrossSize = size
		}
	}
	return maxCrossSize
}

// collectItemInfo は、子ウィジェットからレイアウト計算に必要な情報を収集します。
func collectItemInfo(children []component.Widget, isRow bool) []*flexItemInfo {
	items := make([]*flexItemInfo, len(children))
	for i, child := range children {
		s := child.ReadOnlyStyle()
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

		items[i] = &flexItemInfo{
			widget:          child,
			flex:            child.GetFlex(),
			mainMargin:      mainMargin,
			crossMargin:     crossMargin,
			mainMarginStart: mainMarginStart,
		}
	}
	return items
}

// calculateBaseSizes は、各アイテムの基本サイズを決定します。
func calculateBaseSizes(items []*flexItemInfo, isRow bool, crossAvailable int, alignItems Alignment) {
	for _, item := range items {
		minW, minH := item.widget.GetMinSize()

		var mainMin int
		if isRow {
			mainMin = minW
		} else {
			mainMin = minH
		}

		if item.flex > 0 {
			item.mainSize = mainMin
		} else {
			// flexでないアイテムは、自身のMeasure結果を基本サイズとする
			childAvailable := image.Point{X: utils.IfThen(isRow, 0, crossAvailable-item.crossMargin), Y: utils.IfThen(isRow, crossAvailable-item.crossMargin, 0)}
			if childAvailable.X < 0 {
				childAvailable.X = 0
			}
			if childAvailable.Y < 0 {
				childAvailable.Y = 0
			}
			measuredSize := item.widget.Measure(childAvailable)
			if isRow {
				item.mainSize = measuredSize.X
			} else {
				item.mainSize = measuredSize.Y
			}
		}

		// crossSizeの基本もMeasure結果から取得
		measuredSize := item.widget.Measure(image.Point{X: utils.IfThen(isRow, item.mainSize, crossAvailable-item.crossMargin), Y: utils.IfThen(isRow, crossAvailable-item.crossMargin, item.mainSize)})
		if isRow {
			item.crossSize = measuredSize.Y
		} else {
			item.crossSize = measuredSize.X
		}
	}
}

// distributeRemainingSpace は、残りの空間をflexアイテムに分配します。
func distributeRemainingSpace(items []*flexItemInfo, mainSize, totalBaseMainSize int, totalFlex float64, gap int) {
	totalGap := utils.IfThen(len(items) > 1, (len(items)-1)*gap, 0)
	remainingSpace := mainSize - totalBaseMainSize - totalGap

	if totalFlex > 0 && remainingSpace > 0 {
		sizePerFlex := float64(remainingSpace) / totalFlex
		for _, item := range items {
			if item.flex > 0 {
				item.mainSize += int(sizePerFlex * float64(item.flex))
			}
		}
	}
}

// calculateCrossAxisSizes は、交差軸のサイズを計算します。
func calculateCrossAxisSizes(items []*flexItemInfo, crossSize int, isRow bool, alignItems Alignment) {
	for _, item := range items {
		finalCrossSize := crossSize - item.crossMargin
		if alignItems != AlignStretch {
			measuredSize := item.widget.Measure(image.Point{X: utils.IfThen(isRow, item.mainSize, 0), Y: utils.IfThen(isRow, 0, item.mainSize)})
			intrinsicCrossSize := utils.IfThen(isRow, measuredSize.Y, measuredSize.X)
			if intrinsicCrossSize < finalCrossSize {
				finalCrossSize = intrinsicCrossSize
			}
		}
		if finalCrossSize < 0 {
			finalCrossSize = 0
		}
		item.crossSize = finalCrossSize
	}
}

// applySizes は、計算されたサイズを各ウィジェットに設定します。
func applySizes(items []*flexItemInfo, isRow bool) {
	for _, item := range items {
		if isRow {
			item.widget.SetSize(item.mainSize, item.crossSize)
		} else {
			item.widget.SetSize(item.crossSize, item.mainSize)
		}
	}
}

// positionItems は、主軸と交差軸の揃え位置に基づいて各ウィジェットを配置します。
func positionItems(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool, justify, alignItems Alignment, gap int, finalBounds image.Rectangle, crossOffsetOverride ...int) {
	var currentTotalMainSize int
	totalGap := utils.IfThen(len(items) > 1, (len(items)-1)*gap, 0)

	for _, item := range items {
		currentTotalMainSize += item.mainSize + item.mainMargin
	}
	currentTotalMainSize += totalGap

	freeSpace := mainSize - currentTotalMainSize
	mainOffset := 0.0
	spacing := 0.0

	if freeSpace > 0 {
		switch justify {
		case AlignCenter:
			mainOffset = float64(freeSpace) / 2.0
		case AlignEnd:
			mainOffset = float64(freeSpace)
		}
	}

	padding := container.GetPadding()
	mainStart := utils.IfThen(isRow, finalBounds.Min.X+padding.Left, finalBounds.Min.Y+padding.Top)
	crossStart := utils.IfThen(isRow, finalBounds.Min.Y+padding.Top, finalBounds.Min.X+padding.Left)
	if len(crossOffsetOverride) > 0 {
		crossStart = crossOffsetOverride[0]
	}

	currentMain := float64(mainStart) + mainOffset

	for _, item := range items {
		currentMain += float64(item.mainMarginStart)

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
			item.widget.SetPosition(int(currentMain), finalCrossPos)
		} else {
			item.widget.SetPosition(finalCrossPos, int(currentMain))
		}

		currentMain += float64(item.mainSize + (item.mainMargin - item.mainMarginStart) + gap)
		currentMain += spacing
	}
}