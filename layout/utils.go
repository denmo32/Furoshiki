package layout

import (
	"furoshiki/component"
)

// getVisibleChildren は、コンテナから表示状態の子ウィジェットのみを抽出し、
// レイアウト計算の対象となるスライスを返します。
func getVisibleChildren(container Container) []component.Widget {
	allChildren := container.GetChildren()
	visibleChildren := make([]component.Widget, 0, len(allChildren))
	for _, child := range allChildren {
		// 【提案1】型アサーションの追加: IsVisibleはInteractiveStateインターフェースが持つため、
		// 型アサーションでチェックします。実装していないウィジェットはデフォルトで可視(true)とみなします。
		isVisible := true
		if is, ok := child.(component.InteractiveState); ok {
			isVisible = is.IsVisible()
		}

		if isVisible {
			visibleChildren = append(visibleChildren, child)
		}
	}
	return visibleChildren
}