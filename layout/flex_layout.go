package layout

import (
	"furoshiki/component"
	"furoshiki/style"
)

// FlexLayout は、CSS Flexboxにインスパイアされたレイアウトシステムです。
// Direction, Justify, AlignItems, Gapプロパティに基づいて子要素を柔軟に配置します。
type FlexLayout struct {
	Direction  Direction
	Justify    Alignment
	AlignItems Alignment
	Wrap       bool // 現在は未使用ですが、将来的な機能拡張のために残されています
	Gap        int
}

// flexItemInfo は、レイアウト計算中に各子要素の情報を保持するための中間構造体です。
type flexItemInfo struct {
	widget                  component.Widget
	mainSize, crossSize     int // 最終的な主軸・交差軸サイズ
	mainMargin, crossMargin int // 主軸・交差軸のマージン合計
	mainMarginStart         int // 主軸の開始側マージン
	flex                    int // flex値
}

// Layout は FlexLayout のレイアウトロジックを実装します。
// コンテナのサイズと子のプロパティに基づいて、すべての子要素のサイズと位置を計算し、設定します。
func (l *FlexLayout) Layout(container Container) {
	// ステップ 1: 初期設定と可視コンポーネントのフィルタリング
	children := getVisibleChildren(container)
	if len(children) == 0 {
		return
	}

	padding := container.GetPadding()
	containerWidth, containerHeight := container.GetSize()

	// コンテナに描画領域がない場合はレイアウトをスキップします。
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

	// ステップ 2: レイアウト情報の収集と初期サイズ計算
	items, totalFixedMainSize, totalFlex := l.calculateInitialSizes(children, isRow)

	// ステップ 3: 伸縮可能スペースの計算と分配
	l.distributeRemainingSpace(items, mainSize, totalFixedMainSize, totalFlex, isRow)

	// ステップ 4: 交差軸のサイズ計算
	l.calculateCrossAxisSizes(items, crossSize, isRow)

	// ステップ 5: 最終的なサイズをウィジェットに設定
	applySizes(items, isRow)

	// ステップ 6: 位置計算と配置
	l.positionItems(items, container, mainSize, crossSize, isRow)
}

// getVisibleChildren は、コンテナから表示状態の子ウィジェットのみを抽出します。
func getVisibleChildren(container Container) []component.Widget {
	allChildren := container.GetChildren()
	visibleChildren := make([]component.Widget, 0, len(allChildren))
	for _, child := range allChildren {
		if child.IsVisible() {
			visibleChildren = append(visibleChildren, child)
		}
	}
	return visibleChildren
}

// calculateInitialSizes は、各子要素の初期サイズとマージンを計算します。
func (l *FlexLayout) calculateInitialSizes(children []component.Widget, isRow bool) ([]flexItemInfo, int, float64) {
	items := make([]flexItemInfo, len(children))
	var totalFixedMainSize int
	var totalFlex float64

	for i, child := range children {
		s := child.GetStyle()
		minW, minH := child.GetMinSize()
		w, h := child.GetSize()

		var info flexItemInfo
		info.widget = child
		info.flex = child.GetFlex()

		// Marginが設定されていればその値を、なければゼロ値を使用します。
		margin := style.Insets{}
		if s.Margin != nil {
			margin = *s.Margin
		}

		if isRow {
			info.mainMarginStart = margin.Left
			info.mainMargin = margin.Left + margin.Right
			info.crossMargin = margin.Top + margin.Bottom
		} else {
			info.mainMarginStart = margin.Top
			info.mainMargin = margin.Top + margin.Bottom
			info.crossMargin = margin.Left + margin.Right
		}

		if info.flex > 0 {
			totalFlex += float64(info.flex)
			info.mainSize = max(0, ifThen(isRow, minW, minH))
		} else {
			if isRow {
				info.mainSize = max(ifThen(w <= 0, minW, w), minW)
			} else {
				info.mainSize = max(ifThen(h <= 0, minH, h), minH)
			}
		}
		totalFixedMainSize += info.mainSize + info.mainMargin
		items[i] = info
	}
	return items, totalFixedMainSize, totalFlex
}

// distributeRemainingSpace は、flex値に基づいて余剰スペースを分配または不足分を縮小します。
func (l *FlexLayout) distributeRemainingSpace(items []flexItemInfo, mainSize, totalFixedMainSize int, totalFlex float64, isRow bool) {
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * l.Gap
	}
	remainingSpace := mainSize - totalFixedMainSize - totalGap

	if remainingSpace < 0 {
		// スペースが不足している場合、flex値を持たない（固定サイズの）アイテムを、それぞれの主軸サイズに比例して縮小させます。
		// ただし、各アイテムは自身の最小サイズより小さくなることはありません。
		// これはCSS Flexboxの`flex-shrink`の挙動を簡略化したものです。
		if totalFixedMainSize > 0 {
			// scaleは、利用可能なスペースを要求されたスペースで割った縮小率です。
			scale := float64(mainSize-totalGap) / float64(totalFixedMainSize)
			if scale > 0 { // scaleが0以下だとすべてのサイズが0になってしまうためチェック
				for i := range items {
					// flexアイテムは後で調整されるため、ここでは固定サイズアイテムのみを対象とします。
					if items[i].flex == 0 {
						// 縮小後の期待サイズを計算します。
						newSize := int(float64(items[i].mainSize+items[i].mainMargin) * scale)

						// ウィジェットの最小サイズを取得します。
						minW, minH := items[i].widget.GetMinSize()
						minMainSize := ifThen(isRow, minW, minH)

						// 最小サイズを下回らないようにサイズを決定します。
						items[i].mainSize = max(minMainSize, newSize-items[i].mainMargin)
					}
				}
			}
		}
	} else if totalFlex > 0 && remainingSpace > 0 {
		// スペースに余剰がある場合、flexアイテムにflex値に応じて分配します。
		sizePerFlex := float64(remainingSpace) / totalFlex
		for i := range items {
			if items[i].flex > 0 {
				flexSize := int(sizePerFlex * float64(items[i].flex))
				items[i].mainSize += flexSize
			}
		}
	}
}

// calculateCrossAxisSizes は、交差軸のサイズを計算します（AlignStretch対応）。
func (l *FlexLayout) calculateCrossAxisSizes(items []flexItemInfo, crossSize int, isRow bool) {
	if l.AlignItems == AlignStretch {
		// AlignStretchの場合、すべての子をコンテナの交差軸サイズいっぱいに引き伸ばします。
		for i := range items {
			minW, minH := items[i].widget.GetMinSize()
			// ただし、子の最小サイズは下回らないようにします。
			items[i].crossSize = max(ifThen(isRow, minH, minW), crossSize-items[i].crossMargin)
		}
	} else {
		// それ以外の場合は、子の本来のサイズまたは最小サイズを使用します。
		for i := range items {
			w, h := items[i].widget.GetSize()
			minW, minH := items[i].widget.GetMinSize()
			if isRow {
				items[i].crossSize = max(ifThen(h <= 0, minH, h), minH)
			} else {
				items[i].crossSize = max(ifThen(w <= 0, minW, w), minW)
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
func (l *FlexLayout) positionItems(items []flexItemInfo, container Container, mainSize, crossSize int, isRow bool) {
	var currentTotalMainSize int
	totalGap := 0
	if len(items) > 1 {
		totalGap = (len(items) - 1) * l.Gap
	}

	for _, item := range items {
		currentTotalMainSize += item.mainSize + item.mainMargin
	}
	currentTotalMainSize += totalGap

	// Justifyプロパティに基づいて、主軸方向の開始オフセットを計算します。
	freeSpace := mainSize - currentTotalMainSize
	mainOffset := 0
	switch l.Justify {
	case AlignCenter:
		mainOffset = freeSpace / 2
	case AlignEnd:
		mainOffset = freeSpace
	}
	mainOffset = max(0, mainOffset)

	padding := container.GetPadding()
	containerX, containerY := container.GetPosition()

	mainStart := ifThen(isRow, padding.Left, padding.Top)
	crossStart := ifThen(isRow, padding.Top, padding.Left)
	currentMain := mainStart + mainOffset

	for _, item := range items {
		currentMain += item.mainMarginStart

		// AlignItemsプロパティに基づいて、交差軸方向のオフセットを計算します。
		crossOffset := 0
		switch l.AlignItems {
		case AlignCenter:
			crossOffset = (crossSize - item.crossSize - item.crossMargin) / 2
		case AlignEnd:
			crossOffset = crossSize - item.crossSize - item.crossMargin
		}

		finalCrossPos := crossStart + crossOffset

		// 最終的な絶対座標を設定します。
		if isRow {
			item.widget.SetPosition(containerX+currentMain, containerY+finalCrossPos)
		} else {
			item.widget.SetPosition(containerX+finalCrossPos, containerY+currentMain)
		}

		currentMain += item.mainSize + (item.mainMargin - item.mainMarginStart) + l.Gap
	}
}

// --- Helpers ---

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ifThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}