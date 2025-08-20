package layout

import "furoshiki/component"

// getVisibleChildren は、コンテナから表示状態の子ウィジェットのみを抽出し、
// レイアウト計算の対象となるスライスを返します。
// この関数は複数のレイアウト実装で共通して使用されます。
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