package layout

import "math"

// GridLayout は、子要素を格子状（グリッド）に配置するレイアウトです。
// 子要素は追加された順に、左から右、上から下へとセルに配置されます。
// すべてのセルは同じサイズになります。
type GridLayout struct {
	// Columns はグリッドの列数を指定します。1以上の値を設定する必要があります。
	Columns int
	// Rows はグリッドの行数を指定します。0または負の値を指定した場合、
	// 子要素の数と列数から自動的に計算されます。
	Rows int
	// HorizontalGap は、グリッドの列（セル）間の水平方向の間隔です。
	HorizontalGap int
	// VerticalGap は、グリッドの行（セル）間の垂直方向の間隔です。
	VerticalGap int
}

// Layout は GridLayout のレイアウトロジックを実装します。
// コンテナの利用可能なスペースを列数と行数で均等に分割し、各子要素をセルに配置します。
func (l *GridLayout) Layout(container Container) {
	// [改良] 共通化された getVisibleChildren を使用します。
	children := getVisibleChildren(container)
	childCount := len(children)
	if childCount == 0 {
		return
	}

	// 列数が1未満の場合は、デフォルトで1列に設定します。
	columns := l.Columns
	if columns < 1 {
		columns = 1
	}

	// 行数を計算します。
	// l.Rowsが指定されていればそれを使用し、そうでなければ子の数と列数から計算します。
	rows := l.Rows
	if rows <= 0 {
		// math.Ceilを使って、すべての子を収めるのに必要な行数を計算します。
		// 例: 7個の子を3列で配置する場合 -> ceil(7.0 / 3.0) = 3行
		rows = int(math.Ceil(float64(childCount) / float64(columns)))
	}
	if rows == 0 {
		return
	}

	padding := container.GetPadding()
	containerX, containerY := container.GetPosition()
	containerWidth, containerHeight := container.GetSize()

	// パディングを引いた、実際に子を配置できる領域のサイズを計算します。
	availableWidth := containerWidth - padding.Left - padding.Right
	availableHeight := containerHeight - padding.Top - padding.Bottom

	// すべてのギャップの合計サイズを計算します。
	totalHorizontalGap := (columns - 1) * l.HorizontalGap
	totalVerticalGap := (rows - 1) * l.VerticalGap

	// ギャップを除いた、セル自体が占めることができる合計サイズを計算します。
	totalCellWidth := availableWidth - totalHorizontalGap
	totalCellHeight := availableHeight - totalVerticalGap

	// 各セルの幅と高さを均等に計算します。
	cellWidth := totalCellWidth / columns
	cellHeight := totalCellHeight / rows

	// すべての子をループし、位置とサイズを設定します。
	for i, child := range children {
		// 非表示の子はスキップ (getVisibleChildrenでフィルタリング済みですが念のため)
		if !child.IsVisible() {
			continue
		}

		// インデックスから現在の行と列を計算します。
		row := i / columns
		col := i % columns

		// セルの左上のX座標を計算します。
		// コンテナの開始位置 + 左パディング + (列インデックス * (セル幅 + ギャップ))
		cellX := containerX + padding.Left + col*(cellWidth+l.HorizontalGap)

		// セルの左上のY座標を計算します。
		// コンテナの開始位置 + 上パディング + (行インデックス * (セル高さ + ギャップ))
		cellY := containerY + padding.Top + row*(cellHeight+l.VerticalGap)

		// 子ウィジェットに計算した位置とサイズを設定します。
		// これにより、子はセルのサイズいっぱいに引き伸ばされます。
		child.SetPosition(cellX, cellY)
		child.SetSize(cellWidth, cellHeight)
	}
}

// [削除] getVisibleChildren は共通の layout/utils.go に移動しました。