package layout

import (
	"furoshiki/component"
	"math"
)

// GridLayout は、子要素を格子状（グリッド）に配置するレイアウトです。
type GridLayout struct {
	Columns       int
	Rows          int
	HorizontalGap int
	VerticalGap   int
}

// Layout は GridLayout のレイアウトロジックを実装します。
// NOTE: Layoutインターフェースの変更に伴い、errorを返すようにシグネチャが更新されました。
func (l *GridLayout) Layout(container Container) error {
	children := getVisibleChildren(container)
	childCount := len(children)
	if childCount == 0 {
		return nil
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
		return nil
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

		// 【提案1】型アサーションの追加: 位置とサイズの設定はそれぞれ
		// PositionSetterとSizeSetterインターフェースが持つため、型アサーションを行います。
		if ps, ok := child.(component.PositionSetter); ok {
			ps.SetPosition(cellX, cellY)
		}
		if ss, ok := child.(component.SizeSetter); ok {
			ss.SetSize(cellWidth, cellHeight)
		}
	}
	return nil
}