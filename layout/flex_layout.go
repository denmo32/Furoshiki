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
	// AlignContent は、複数行/列になった際の、交差軸方向のラインの揃え位置を設定します。
	// このプロパティは、Wrapがtrueの場合にのみ効果があります。
	AlignContent Alignment
	// Wrap は、アイテムが一行に収らない場合に折り返すかどうかを指定します。
	Wrap bool
	Gap  int
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
	crossAxisSize int // このラインの交差軸方向のサイズ
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
	// VStacksで正しい高さを計算するために、crossSize(幅)とAlignItemsを渡します。
	calculateBaseSizes(items, isRow, crossSize, l.AlignItems)

	if l.Wrap {
		l.layoutMultiLine(items, container, mainSize, crossSize, isRow)
	} else {
		l.layoutSingleLine(items, container, mainSize, crossSize, isRow)
	}
}

// layoutSingleLine は、折り返しなしのレイアウト計算を実行します。
func (l *FlexLayout) layoutSingleLine(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool) {
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
	// シングルラインの場合、最終的なサイズを適用してから配置します。
	applySizes(items, isRow)
	positionItems(items, container, mainSize, crossSize, isRow, l.Justify, l.AlignItems, l.Gap)
}

// layoutMultiLine は、折り返しありのレイアウト計算を実行します。
func (l *FlexLayout) layoutMultiLine(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool) {
	// 1. アイテムを複数のラインに分割
	lines := l.splitIntoLines(items, mainSize)

	// 2. 各ラインのサイズを計算
	var totalCrossAxisSize int
	for _, line := range lines {
		// ライン内のflexアイテムに余剰スペースを分配
		var lineTotalFlex float64
		var lineTotalBaseMainSize int
		for _, item := range line.items {
			if item.flex > 0 {
				lineTotalFlex += float64(item.flex)
			}
			lineTotalBaseMainSize += item.mainSize + item.mainMargin
		}
		distributeRemainingSpace(line.items, mainSize, lineTotalBaseMainSize, lineTotalFlex, l.Gap)

		// ライン内のアイテムの交差軸サイズと、ライン自体の交差軸サイズを計算
		calculateCrossAxisSizes(line.items, crossSize, isRow, l.AlignItems)
		line.crossAxisSize = calculateLineCrossSize(line.items)
		totalCrossAxisSize += line.crossAxisSize
	}

	// 3. ラインを交差軸方向に配置 (AlignContent)
	padding := container.GetPadding()
	crossStart := utils.IfThen(isRow, padding.Top, padding.Left)
	currentCross := crossStart

	freeCrossSpace := crossSize - totalCrossAxisSize
	if len(lines) > 1 {
		freeCrossSpace -= (len(lines) - 1) * l.Gap
	}

	if freeCrossSpace > 0 {
		switch l.AlignContent {
		case AlignCenter:
			currentCross += freeCrossSpace / 2
		case AlignEnd:
			currentCross += freeCrossSpace
			// NOTE: AlignStretch や AlignSpaceBetween などは、より複雑な計算が必要なため現バージョンでは未サポートです。
		}
	}

	// 4. 各ライン内のアイテムを最終配置
	for _, line := range lines {
		// positionItemsをラインごとに呼び出し、ラインの開始位置 (currentCross) を渡してアイテムを配置します。
		positionItems(line.items, container, mainSize, line.crossAxisSize, isRow, l.Justify, l.AlignItems, l.Gap, currentCross)
		currentCross += line.crossAxisSize + l.Gap
	}

	// 5. 計算された最終的なサイズを全ウィジェットに適用
	// マルチラインの場合、すべての配置計算が終わった後にサイズを適用します。
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
		gap := 0
		if len(currentLineItems) > 0 {
			gap = l.Gap
		}

		// アイテムを追加すると主軸サイズを超える場合、現在のラインを確定して新しいラインを開始します。
		if len(currentLineItems) > 0 && currentMainSize+itemMainSize+gap > mainSize {
			lines = append(lines, &flexLine{items: currentLineItems})
			currentLineItems = make([]*flexItemInfo, 0)
			currentMainSize = 0
			gap = 0
		}

		currentLineItems = append(currentLineItems, item)
		currentMainSize += itemMainSize + gap
	}

	// 最後のラインを追加します。
	if len(currentLineItems) > 0 {
		lines = append(lines, &flexLine{items: currentLineItems})
	}

	return lines
}

// calculateLineCrossSize は、ライン内のアイテムに基づいてライン自体の交差軸サイズを決定します。
// ラインのサイズは、その中の最も大きいアイテムのサイズに合わせられます。
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
// ポインタのスライスを返すように変更しました。
func collectItemInfo(children []component.Widget, isRow bool) []*flexItemInfo {
	items := make([]*flexItemInfo, len(children))
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
// VStacks (`isRow == false`) のために、crossSize と alignItems を受け取るように修正されました。
func calculateBaseSizes(items []*flexItemInfo, isRow bool, crossSize int, alignItems Alignment) {
	for _, item := range items {
		w, h := item.widget.GetSize()
		minW, minH := item.widget.GetMinSize()

		if isRow { // HStack のロジックは変更なし
			if item.flex > 0 {
				item.mainSize = minW
			} else {
				item.mainSize = max(utils.IfThen(w <= 0, minW, w), minW)
			}
		} else { // VStack のための新しいロジック
			// mainSize は高さであり、幅(crossSize)に依存する可能性があるため、先に幅を決定します。
			itemWidth := crossSize - item.crossMargin // 利用可能な最大幅から開始
			if alignItems != AlignStretch {
				// stretchでない場合、アイテムは自身の本来の幅を使います。
				intrinsicWidth := max(utils.IfThen(w <= 0, minW, w), minW)
				if intrinsicWidth < itemWidth {
					itemWidth = intrinsicWidth
				}
			}
			if itemWidth < 0 {
				itemWidth = 0
			}

			// 確定した幅を使って、正しい基本の高さを計算します。
			if hw, ok := item.widget.(component.HeightForWider); ok {
				item.mainSize = hw.GetHeightForWidth(itemWidth)
			} else {
				// 折り返しをサポートしないウィジェットのフォールバック
				item.mainSize = max(utils.IfThen(h <= 0, minH, h), minH)
			}
		}
	}
}

// distributeRemainingSpace は、残りの空間をflexアイテムに分配します。
// ポインタのスライスを受け取るように変更しました。
func distributeRemainingSpace(items []*flexItemInfo, mainSize, totalBaseMainSize int, totalFlex float64, gap int) {
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * gap
	}

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
// ポインタのスライスを受け取るように変更しました。
func calculateCrossAxisSizes(items []*flexItemInfo, crossSize int, isRow bool, alignItems Alignment) {
	for _, item := range items {
		// デフォルトでは、子は親の利用可能な交差軸スペース全体を占有します (AlignStretch)。
		item.crossSize = crossSize - item.crossMargin

		// AlignStretchでない場合、子は自身のコンテンツに合わせたサイズになることができます。
		if alignItems != AlignStretch {
			var intrinsicCrossSize int
			w, h := item.widget.GetSize()
			minW, minH := item.widget.GetMinSize()

			if isRow { // HStack の交差軸(高さ)を計算
				if hw, ok := item.widget.(component.HeightForWider); ok {
					// アイテムの幅(mainSize)は既に確定しているので、それに基づき正しい高さを計算
					intrinsicCrossSize = hw.GetHeightForWidth(item.mainSize)
				} else {
					intrinsicCrossSize = max(utils.IfThen(h <= 0, minH, h), minH)
				}
			} else { // VStack の交差軸(幅)を計算
				// 幅はテキストの折り返しに依存しないので、単純に本来の幅を計算
				intrinsicCrossSize = max(utils.IfThen(w <= 0, minW, w), minW)
			}

			if intrinsicCrossSize > 0 && intrinsicCrossSize < item.crossSize {
				item.crossSize = intrinsicCrossSize
			}
		}

		// サイズが負の値にならないように保証します。
		if item.crossSize < 0 {
			item.crossSize = 0
		}
	}
}

// applySizes は、計算されたサイズを各ウィジェットに設定します。
// ポインタのスライスを受け取るように変更しました。
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
// crossOffsetOverride をオプションの引数として追加し、マルチラインレイアウト時に
// ラインの開始位置を指定できるように変更しました。
// ポインタのスライスを受け取るように変更しました。
func positionItems(items []*flexItemInfo, container Container, mainSize, crossSize int, isRow bool, justify, alignItems Alignment, gap int, crossOffsetOverride ...int) {
	var currentTotalMainSize int
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * gap
	}

	// 主軸方向のアイテムの合計サイズを計算
	for _, item := range items {
		currentTotalMainSize += item.mainSize + item.mainMargin
	}
	currentTotalMainSize += totalGap

	// Justifyプロパティに基づいて主軸方向の開始オフセットを計算
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

	// 交差軸の開始位置を決定。マルチラインの場合はオーバーライド値を使用します。
	crossStart := utils.IfThen(isRow, padding.Top, padding.Left)
	if len(crossOffsetOverride) > 0 {
		crossStart = crossOffsetOverride[0]
	}

	currentMain := mainStart + mainOffset

	for _, item := range items {
		currentMain += item.mainMarginStart

		// AlignItemsプロパティに基づいて交差軸方向のオフセットを計算
		crossOffset := 0
		availableCrossSpace := crossSize - item.crossSize - item.crossMargin
		if availableCrossSpace > 0 {
			switch alignItems {
			case AlignCenter:
				crossOffset = availableCrossSpace / 2
			case AlignEnd:
				crossOffset = availableCrossSpace
				// AlignStretch の場合、item.crossSize が既に crossSize - item.crossMargin になっているため offset は 0
			}
		}

		finalCrossPos := crossStart + crossOffset

		// 最終的な座標を設定
		if isRow {
			item.widget.SetPosition(containerX+currentMain, containerY+finalCrossPos)
		} else {
			item.widget.SetPosition(containerX+finalCrossPos, containerY+currentMain)
		}

		currentMain += item.mainSize + (item.mainMargin - item.mainMarginStart) + gap
	}
}