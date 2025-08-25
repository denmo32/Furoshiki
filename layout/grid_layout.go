package layout

import "math"

// GridLayout は、子要素を格子状（グリッド）に配置するレイアウトです。
type GridLayout struct {
	Columns       int
	Rows          int
	HorizontalGap int
	VerticalGap   int
}

// Layout は GridLayout のレイアウトロジックを実装します。
func (l *GridLayout) Layout(container Container) {
	children := getVisibleChildren(container)
	childCount := len(children)
	if childCount == 0 {
		return
	}

	columns := l.Columns
	if columns < 1 {
		columns = 1
	}

	rows := l.Rows
	if rows <= 0 {
		rows = int(math.Ceil(float64(childCount) / float64(columns)))
	}
	if rows == 0 {
		return
	}

	padding := container.GetPadding()
	containerX, containerY := container.GetPosition()
	containerWidth, containerHeight := container.GetSize()

	availableWidth := containerWidth - padding.Left - padding.Right
	availableHeight := containerHeight - padding.Top - padding.Bottom

	totalHorizontalGap := (columns - 1) * l.HorizontalGap
	totalVerticalGap := (rows - 1) * l.VerticalGap

	cellWidth := (availableWidth - totalHorizontalGap) / columns
	cellHeight := (availableHeight - totalVerticalGap) / rows

	for i, child := range children {
		row := i / columns
		col := i % columns

		cellX := containerX + padding.Left + col*(cellWidth+l.HorizontalGap)
		cellY := containerY + padding.Top + row*(cellHeight+l.VerticalGap)

		child.SetPosition(cellX, cellY)
		child.SetSize(cellWidth, cellHeight)
	}
}
