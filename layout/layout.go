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
	GetChildren() []component.Component
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

func (l *FlexLayout) Layout(container Container) {
	allChildren := container.GetChildren()
	// レイアウト対象となるのは、表示されている子要素のみ
	var children []component.Component
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

	// コンテナのPaddingを考慮した、子要素が利用可能な領域を計算
	availableWidth := containerWidth - containerStyle.Padding.Left - containerStyle.Padding.Right
	availableHeight := containerHeight - containerStyle.Padding.Top - containerStyle.Padding.Bottom

	isRow := l.Direction == DirectionRow
	mainSize, crossSize := availableWidth, availableHeight
	if !isRow {
		mainSize, crossSize = availableHeight, availableWidth
	}

	// --- 1. サイズ計算フェーズ ---
	// このフェーズでは、各子要素の最終的なサイズを決定します。
	// Flexboxの計算は、主軸（Main Axis）と交差軸（Cross Axis）の2つの軸を基準に行われます。
	// DirectionがRowなら主軸は水平方向、Columnなら垂直方向です。

	// --- ステップ 1.1: レイアウト情報の収集 ---
	// まず、すべての子要素をループし、以下の情報を収集します。
	// - Flex値を持つ子要素（伸縮する要素）のリスト
	// - Flex値の合計 (スペースを分配する際の比率として使用)
	// - 固定サイズ要素の主軸方向の合計サイズ (Flexアイテムの最小サイズも含む)
	var totalFixedMainSize int
	var totalFlex float64
	var flexibleChildren []component.Component

	for _, child := range children {
		childStyle := child.GetStyle()
		minW, minH := child.GetMinSize()
		childMarginMain := 0
		if isRow {
			childMarginMain = childStyle.Margin.Left + childStyle.Margin.Right
		} else {
			childMarginMain = childStyle.Margin.Top + childStyle.Margin.Bottom
		}

		if child.GetFlex() > 0 {
			totalFlex += float64(child.GetFlex())
			flexibleChildren = append(flexibleChildren, child)
			// flexアイテムでも、その最小サイズとマージンは固定サイズとして事前に確保する
			if isRow {
				totalFixedMainSize += max(0, minW) + childMarginMain
			} else {
				totalFixedMainSize += max(0, minH) + childMarginMain
			}
		} else {
			// 固定サイズアイテムのサイズ（最小サイズを考慮）とマージンを合計
			if isRow {
				w, _ := child.GetSize()
				totalFixedMainSize += max(w, minW) + childMarginMain
			} else {
				_, h := child.GetSize()
				totalFixedMainSize += max(h, minH) + childMarginMain
			}
		}
	}

	// --- ステップ 1.2: 伸縮可能スペースの計算 ---
	// コンテナの利用可能サイズから、固定サイズの要素とGapが占めるスペースを引いて、
	// 伸縮可能な要素（Flexアイテム）に分配できる残りのスペースを計算します。
	totalGap := 0
	if len(children) > 1 {
		totalGap = (len(children) - 1) * l.Gap
	}
	remainingSpace := mainSize - totalFixedMainSize - totalGap // Gapの合計も考慮に入れる

	sizePerFlex := 0.0
	if totalFlex > 0 && remainingSpace > 0 {
		sizePerFlex = float64(remainingSpace) / totalFlex
	}

	// --- ステップ 1.3: Flexアイテムの主軸サイズを決定 ---
	// 計算された残りのスペースを、各FlexアイテムのFlex値の比率に応じて分配します。
	// 最終的なサイズは「最小サイズ + 分配されたスペース」となり、最小サイズが保証されます。
	// これにより、コンテナが縮んでもアイテムが必要最低限のサイズを維持できます。
	for _, child := range flexibleChildren {
		flexSize := int(sizePerFlex * float64(child.GetFlex()))
		minW, minH := child.GetMinSize()
		if isRow {
			_, h := child.GetSize()
			// 最終的な幅 = 最小幅 + 追加のflexスペース
			finalWidth := max(minW, minW+flexSize)
			child.SetSize(finalWidth, h)
		} else {
			w, _ := child.GetSize()
			// 最終的な高さ = 最小高 + 追加のflexスペース
			finalHeight := max(minH, minH+flexSize)
			child.SetSize(w, finalHeight)
		}
	}

	// --- ステップ 1.4: 固定サイズアイテムの最小サイズを適用 ---
	// Flex値を持たないアイテムにも、設定された最小サイズが適用されるようにします。
	// ユーザーが指定したサイズが最小サイズより小さい場合、最小サイズが優先されます。
	for _, child := range children {
		if child.GetFlex() == 0 {
			w, h := child.GetSize()
			minW, minH := child.GetMinSize()
			child.SetSize(max(w, minW), max(h, minH))
		}
	}

	// --- ステップ 1.5: 交差軸のサイズを決定 (AlignStretch) ---
	// AlignItemsプロパティがAlignStretchに設定されている場合、
	// 子要素はコンテナの交差軸いっぱいに引き伸ばされます。
	// ここでもマージンを考慮し、最小サイズが保証されます。
	if l.AlignItems == AlignStretch {
		for _, child := range children {
			childStyle := child.GetStyle()
			minW, minH := child.GetMinSize()

			if isRow {
				w, _ := child.GetSize()
				childCrossMargin := childStyle.Margin.Top + childStyle.Margin.Bottom
				finalHeight := max(minH, crossSize-childCrossMargin)
				child.SetSize(w, finalHeight)
			} else {
				_, h := child.GetSize()
				childCrossMargin := childStyle.Margin.Left + childStyle.Margin.Right
				finalWidth := max(minW, crossSize-childCrossMargin)
				child.SetSize(finalWidth, h)
			}
		}
	}

	// --- 2. 位置計算フェーズ ---
	// このフェーズでは、ステップ1で確定したサイズを持つ各子要素を、
	// Justify (主軸方向の揃え) と AlignItems (交差軸方向の揃え) の設定に従って
	// コンテナ内に実際に配置していきます。

	// --- ステップ 2.1: 主軸方向の開始オフセットを計算 ---
	// まず、配置後の全要素の合計サイズを再計算します。
	var currentTotalMainSize int
	for _, child := range children {
		childStyle := child.GetStyle()
		if isRow {
			w, _ := child.GetSize()
			currentTotalMainSize += w + childStyle.Margin.Left + childStyle.Margin.Right
		} else {
			_, h := child.GetSize()
			currentTotalMainSize += h + childStyle.Margin.Top + childStyle.Margin.Bottom
		}
	}
	currentTotalMainSize += totalGap

	// Justify設定に応じて、要素群全体の開始位置（オフセット）を決定します。
	freeSpace := mainSize - currentTotalMainSize
	mainOffset := 0
	switch l.Justify {
	case AlignStart:
		mainOffset = 0
	case AlignCenter:
		mainOffset = freeSpace / 2
	case AlignEnd:
		mainOffset = freeSpace
	}
	if mainOffset < 0 { // コンテンツがはみ出す場合は、先頭から配置
		mainOffset = 0
	}

	// --- ステップ 2.2: 各要素の位置を決定し配置 ---
	// コンテナのPaddingと計算したオフセットを開始点として、各子要素を順番に配置します。
	mainStart := containerStyle.Padding.Left
	crossStart := containerStyle.Padding.Top
	if !isRow {
		mainStart = containerStyle.Padding.Top
		crossStart = containerStyle.Padding.Left
	}
	currentMain := mainStart + mainOffset

	for _, child := range children {
		childStyle := child.GetStyle()
		var childMainSize, childCrossSize int
		var childMarginMainStart, childMarginMainEnd, childMarginCrossStart, childMarginCrossEnd int

		if isRow {
			childMainSize, childCrossSize = child.GetSize()
			childMarginMainStart, childMarginMainEnd = childStyle.Margin.Left, childStyle.Margin.Right
			childMarginCrossStart, childMarginCrossEnd = childStyle.Margin.Top, childStyle.Margin.Bottom
		} else {
			childCrossSize, childMainSize = child.GetSize()
			childMarginMainStart, childMarginMainEnd = childStyle.Margin.Top, childStyle.Margin.Bottom
			childMarginCrossStart, childMarginCrossEnd = childStyle.Margin.Left, childStyle.Margin.Right
		}

		// 主軸方向の位置を決定 (マージンを考慮)
		currentMain += childMarginMainStart

		// 交差軸方向の位置をAlignItems設定に応じて決定
		crossOffset := 0
		switch l.AlignItems {
		case AlignStart, AlignStretch:
			crossOffset = 0
		case AlignCenter:
			crossOffset = (crossSize - childCrossSize - childMarginCrossStart - childMarginCrossEnd) / 2
		case AlignEnd:
			crossOffset = crossSize - childCrossSize - childMarginCrossEnd
		}
		finalCrossPos := crossStart + crossOffset + childMarginCrossStart

		// 最終的な座標を設定
		if isRow {
			child.SetPosition(containerX+currentMain, containerY+finalCrossPos)
		} else {
			child.SetPosition(containerX+finalCrossPos, containerY+currentMain)
		}

		// 次の要素の開始位置を計算 (要素サイズ + マージン + Gap)
		currentMain += childMainSize + childMarginMainEnd + l.Gap
	}
}
