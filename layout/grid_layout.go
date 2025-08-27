package layout

import (
	"image"
	"math"
)

// GridLayout は、子要素を格子状（グリッド）に配置するレイアウトです。
type GridLayout struct {
	Columns       int
	Rows          int
	HorizontalGap int
	VerticalGap   int
}

// calculateGridDimensions は、子要素の数と設定に基づいてグリッドの列数と行数を計算します。
func (l *GridLayout) calculateGridDimensions(childCount int) (cols, rows int) {
	cols = l.Columns
	if cols < 1 {
		cols = 1
	}

	rows = l.Rows
	if rows <= 0 {
		// 行数が指定されていない場合、列数と子の数から計算します。
		if childCount == 0 {
			rows = 0
		} else {
			rows = int(math.Ceil(float64(childCount) / float64(cols)))
		}
	}
	return cols, rows
}

// Measure は、GridLayoutの要求サイズを計算します。
// すべての子を同じサイズ（最も大きい子のサイズ）と仮定して、グリッド全体のサイズを決定します。
func (l *GridLayout) Measure(container Container, availableSize image.Point) image.Point {
	children := getVisibleChildren(container)
	childCount := len(children)
	if childCount == 0 {
		return image.Point{}
	}

	// すべての子を計測し、最大の要求幅と高さを求めます。
	maxChildW, maxChildH := 0, 0
	for _, child := range children {
		// グリッドレイアウトでは、子のサイズは最終的にグリッドセルによって決定されるため、
		// 計測時点では利用可能な最大サイズを渡して、子の理想的なサイズを把握します。
		childDesiredSize := child.Measure(availableSize)
		if childDesiredSize.X > maxChildW {
			maxChildW = childDesiredSize.X
		}
		if childDesiredSize.Y > maxChildH {
			maxChildH = childDesiredSize.Y
		}
	}

	columns, rows := l.calculateGridDimensions(childCount)
	if columns == 0 || rows == 0 {
		return image.Point{}
	}

	padding := container.GetPadding()

	// 最大の子サイズを基に、グリッド全体の要求サイズを計算します。
	totalHorizontalGap := (columns - 1) * l.HorizontalGap
	totalVerticalGap := (rows - 1) * l.VerticalGap

	desiredW := maxChildW*columns + totalHorizontalGap + padding.Left + padding.Right
	desiredH := maxChildH*rows + totalVerticalGap + padding.Top + padding.Bottom

	return image.Point{X: desiredW, Y: desiredH}
}

// Arrange は、GridLayout内の子要素を最終的な位置に配置します。
func (l *GridLayout) Arrange(container Container, finalBounds image.Rectangle) error {
	children := getVisibleChildren(container)
	childCount := len(children)
	if childCount == 0 {
		return nil
	}

	columns, rows := l.calculateGridDimensions(childCount)
	if columns == 0 || rows == 0 {
		return nil
	}

	padding := container.GetPadding()

	// コンテナの最終的な描画領域から、パディングとギャップを引いて、
	// 各セルのサイズを計算します。
	availableWidth := finalBounds.Dx() - padding.Left - padding.Right
	availableHeight := finalBounds.Dy() - padding.Top - padding.Bottom

	totalHorizontalGap := (columns - 1) * l.HorizontalGap
	totalVerticalGap := (rows - 1) * l.VerticalGap

	cellWidth := (availableWidth - totalHorizontalGap) / columns
	cellHeight := (availableHeight - totalVerticalGap) / rows

	if cellWidth < 0 {
		cellWidth = 0
	}
	if cellHeight < 0 {
		cellHeight = 0
	}

	for i, child := range children {
		row := i / columns
		col := i % columns

		// セルの位置を計算
		cellX := finalBounds.Min.X + padding.Left + col*(cellWidth+l.HorizontalGap)
		cellY := finalBounds.Min.Y + padding.Top + row*(cellHeight+l.VerticalGap)

		childBounds := image.Rect(cellX, cellY, cellX+cellWidth, cellY+cellHeight)

		// 子のサイズと位置を設定
		child.SetPosition(childBounds.Min.X, childBounds.Min.Y)
		child.SetSize(childBounds.Dx(), childBounds.Dy())

		// 子のArrangeを再帰的に呼び出し
		if err := child.Arrange(childBounds); err != nil {
			return err
		}
	}
	return nil
}