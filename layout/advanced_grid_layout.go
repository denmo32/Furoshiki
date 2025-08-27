package layout

import (
	"furoshiki/component"
	"furoshiki/utils"
	"image"
	"math"
)

var _ component.Widget // Dummy var to force import usage

// --- AdvancedGridLayout ---

// TrackSizing は、グリッドの列または行のサイズ決定方法を定義します。
type TrackSizing int

const (
	// TrackSizingFixed は、トラックサイズをピクセル単位で固定します。
	TrackSizingFixed TrackSizing = iota
	// TrackSizingWeighted は、利用可能な残りのスペースを重みに応じて分配します。
	TrackSizingWeighted
)

// TrackDefinition は、単一の列または行のサイズ定義を保持します。
type TrackDefinition struct {
	Sizing TrackSizing
	Value  float64
}

// GridPlacementData は、AdvancedGridLayout内のウィジェットの配置情報を定義します。
// この構造体のインスタンスは、ウィジェットの `layoutData` フィールドに格納されます。
type GridPlacementData struct {
	Row, Col         int
	RowSpan, ColSpan int
}

// AdvancedGridLayout は、子要素をセル結合や可変サイズ指定が可能なグリッドに配置します。
type AdvancedGridLayout struct {
	ColumnDefinitions []TrackDefinition
	RowDefinitions    []TrackDefinition
	HorizontalGap     int
	VerticalGap       int
}

// Measure は、AdvancedGridLayoutの要求サイズを計算します。
func (l *AdvancedGridLayout) Measure(container Container, availableSize image.Point) image.Point {
	numCols := len(l.ColumnDefinitions)
	numRows := len(l.RowDefinitions)
	if numCols == 0 || numRows == 0 {
		return image.Point{}
	}

	children := getVisibleChildren(container)
	minColWidths := make([]int, numCols)
	minRowHeights := make([]int, numRows)

	// 1. セル結合(span > 1)のないウィジェットを計測し、各トラックの最小サイズを決定
	for _, child := range children {
		data, ok := child.GetLayoutData().(GridPlacementData)
		if !ok {
			continue
		}
		if data.ColSpan == 1 {
			// 利用可能サイズ0で計測し、ウィジェットの純粋な最小サイズを取得
			minSize := child.Measure(image.Point{0, 0})
			if minSize.X > minColWidths[data.Col] {
				minColWidths[data.Col] = minSize.X
			}
		}
		if data.RowSpan == 1 {
			minSize := child.Measure(image.Point{0, 0})
			if minSize.Y > minRowHeights[data.Row] {
				minRowHeights[data.Row] = minSize.Y
			}
		}
	}
	// TODO: セル結合を持つウィジェットの最小サイズを考慮に入れるロジックを追加

	// 2. トラック定義と最小サイズから、全体の要求サイズを計算
	colSizes := calculateTrackSizes(l.ColumnDefinitions, 0, minColWidths)
	rowSizes := calculateTrackSizes(l.RowDefinitions, 0, minRowHeights)

	totalWidth := 0
	for _, w := range colSizes {
		totalWidth += w
	}
	totalHeight := 0
	for _, h := range rowSizes {
		totalHeight += h
	}

	padding := container.GetPadding()
	totalHorizontalGap := max(0, (numCols-1)*l.HorizontalGap)
	totalVerticalGap := max(0, (numRows-1)*l.VerticalGap)

	return image.Point{
		X: totalWidth + totalHorizontalGap + padding.Left + padding.Right,
		Y: totalHeight + totalVerticalGap + padding.Top + padding.Bottom,
	}
}

// Arrange は、AdvancedGridLayout内の子要素を最終的な位置に配置します。
func (l *AdvancedGridLayout) Arrange(container Container, finalBounds image.Rectangle) error {
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return nil
	}

	padding := container.GetPadding()
	availableWidth := finalBounds.Dx() - padding.Left - padding.Right
	availableHeight := finalBounds.Dy() - padding.Top - padding.Bottom

	numCols := len(l.ColumnDefinitions)
	numRows := len(l.RowDefinitions)
	if numCols == 0 || numRows == 0 {
		return nil
	}

	totalHorizontalGap := max(0, (numCols-1)*l.HorizontalGap)
	totalVerticalGap := max(0, (numRows-1)*l.VerticalGap)
	netWidth := availableWidth - totalHorizontalGap
	netHeight := availableHeight - totalVerticalGap

	colWidths := calculateTrackSizes(l.ColumnDefinitions, netWidth, nil)
	rowHeights := calculateTrackSizes(l.RowDefinitions, netHeight, nil)

	colPositions := calculateTrackPositions(colWidths, l.HorizontalGap, finalBounds.Min.X+padding.Left)
	rowPositions := calculateTrackPositions(rowHeights, l.VerticalGap, finalBounds.Min.Y+padding.Top)

	for _, child := range children {
		data, ok := child.GetLayoutData().(GridPlacementData)
		if !ok {
			continue
		}

		startCol := utils.Clamp(data.Col, 0, numCols-1)
		startRow := utils.Clamp(data.Row, 0, numRows-1)
		endCol := utils.Clamp(data.Col+data.ColSpan, startCol+1, numCols)
		endRow := utils.Clamp(data.Row+data.RowSpan, startRow+1, numRows)

		x := colPositions[startCol]
		y := rowPositions[startRow]
		width := (colPositions[endCol-1] + colWidths[endCol-1]) - x
		height := (rowPositions[endRow-1] + rowHeights[endRow-1]) - y

		childBounds := image.Rect(x, y, x+width, y+height)
		child.SetPosition(childBounds.Min.X, childBounds.Min.Y)
		child.SetSize(childBounds.Dx(), childBounds.Dy())

		if err := child.Arrange(childBounds); err != nil {
			return err
		}
	}
	return nil
}

// calculateTrackSizes は、定義と利用可能スペース、最小サイズに基づいて各トラックのサイズを計算します。
func calculateTrackSizes(definitions []TrackDefinition, availableSpace int, minSizes []int) []int {
	sizes := make([]int, len(definitions))
	var totalWeightedValue float64
	remainingSpace := float64(availableSpace)

	// 1. 固定サイズと最小サイズを適用
	for i, def := range definitions {
		minSize := 0
		if minSizes != nil && i < len(minSizes) {
			minSize = minSizes[i]
		}

		if def.Sizing == TrackSizingFixed {
			sizes[i] = max(int(def.Value), minSize)
		} else {
			sizes[i] = minSize
			totalWeightedValue += def.Value
		}
		remainingSpace -= float64(sizes[i])
	}

	remainingSpace = math.Max(0, remainingSpace)

	// 2. 重み付けされたトラックに残りスペースを分配
	if totalWeightedValue > 0 {
		for i, def := range definitions {
			if def.Sizing == TrackSizingWeighted {
				sizes[i] += int(remainingSpace * (def.Value / totalWeightedValue))
			}
		}
	}

	return sizes
}

// calculateTrackPositions は、各トラックの開始座標を計算します。
func calculateTrackPositions(sizes []int, gap int, startOffset int) []int {
	positions := make([]int, len(sizes))
	if len(sizes) == 0 {
		return positions
	}
	currentPos := startOffset
	for i, size := range sizes {
		positions[i] = currentPos
		currentPos += size + gap
	}
	return positions
}