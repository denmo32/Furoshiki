package layout

import (
	"furoshiki/component"
	"reflect"
)

// getVisibleChildren は、コンテナから表示状態の子ウィジェットのみを抽出し、
// レイアウト計算の対象となるスライスを返します。
// この関数は複数のレイアウト実装で共通して使用されます。
func getVisibleChildren(container Container) []component.Widget {
	allChildren := container.GetChildren()
	visibleChildren := make([]component.Widget, 0, len(allChildren))
	for _, child := range allChildren {
		// ui.Builderが*container.Containerをラップしている場合があるため、
		// 型名チェックは安全のために残しても良いが、本質的な問題ではない。
		// ここではシンプルにIsVisible()のみをチェックする。
		if child.IsVisible() {
			// さらに、型名がBuilderであるような内部的なウィジェットを除外するロジックは
			// ここでは不要かもしれないため、一旦コメントアウトまたは削除。
			// 実際のウィジェットのみがツリーに含まれるべき。
			visibleChildren = append(visibleChildren, child)
		}
	}
	return visibleChildren
}

// 【追加】リフレクションのインポートが必要な場合
func getWidgetTypeName(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}