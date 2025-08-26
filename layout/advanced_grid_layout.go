package layout

import (
	"math"
)

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

// Layout は AdvancedGridLayout のレイアウトロジックを実装します。
func (l *AdvancedGridLayout) Layout(container Container) {
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return
	}

	padding := container.GetPadding()
	containerX, containerY := container.GetPosition()
	containerWidth, containerHeight := container.GetSize()

	availableWidth := containerWidth - padding.Left - padding.Right
	availableHeight := containerHeight - padding.Top - padding.Bottom

	numCols := len(l.ColumnDefinitions)
	numRows := len(l.RowDefinitions)
	if numCols == 0 || numRows == 0 {
		return
	}

	// 1. ギャップを考慮に入れた利用可能スペースを計算
	totalHorizontalGap := max(0, (numCols-1)*l.HorizontalGap)
	totalVerticalGap := max(0, (numRows-1)*l.VerticalGap)
	netWidth := availableWidth - totalHorizontalGap
	netHeight := availableHeight - totalVerticalGap

	// 2. 各トラックのサイズを計算
	colWidths := calculateTrackSizes(l.ColumnDefinitions, netWidth)
	rowHeights := calculateTrackSizes(l.RowDefinitions, netHeight)

	// 3. 各トラックの開始位置を計算
	colPositions := calculateTrackPositions(colWidths, l.HorizontalGap, containerX+padding.Left)
	rowPositions := calculateTrackPositions(rowHeights, l.VerticalGap, containerY+padding.Top)

	// 4. 子要素を配置
	for _, child := range children {
		data, ok := child.GetLayoutData().(GridPlacementData)
		if !ok {
			// 配置情報がないウィジェットはレイアウト対象外とします。
			continue
		}

		// 範囲チェックを行い、グリッドの範囲内に収める
		startCol := clamp(data.Col, 0, numCols-1)
		startRow := clamp(data.Row, 0, numRows-1)
		endCol := clamp(data.Col+data.ColSpan, startCol+1, numCols)
		endRow := clamp(data.Row+data.RowSpan, startRow+1, numRows)

		// 位置とサイズを計算
		x := colPositions[startCol]
		y := rowPositions[startRow]
		width := colPositions[endCol-1] + colWidths[endCol-1] - x
		height := rowPositions[endRow-1] + rowHeights[endRow-1] - y

		child.SetPosition(x, y)
		child.SetSize(width, height)
	}
}

// calculateTrackSizes は、定義に基づいて各トラック（列または行）の最終的なサイズを計算します。
func calculateTrackSizes(definitions []TrackDefinition, availableSpace int) []int {
	sizes := make([]int, len(definitions))
	var totalWeightedValue float64
	remainingSpace := float64(availableSpace)

	// 固定サイズのトラックを先に計算し、残りのスペースから引きます。
	for i, def := range definitions {
		if def.Sizing == TrackSizingFixed {
			size := int(def.Value)
			sizes[i] = size
			remainingSpace -= float64(size)
		} else {
			totalWeightedValue += def.Value
		}
	}

	remainingSpace = math.Max(0, remainingSpace)

	// 重み付けされたトラックを、残りのスペースと重みの比率で計算します。
	if totalWeightedValue > 0 {
		for i, def := range definitions {
			if def.Sizing == TrackSizingWeighted {
				sizes[i] = int(remainingSpace * (def.Value / totalWeightedValue))
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

// clamp は値を指定された最小値と最大値の範囲内に収めます。
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}