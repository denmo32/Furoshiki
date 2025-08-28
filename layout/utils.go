package layout

import (
	"furoshiki/component"
	"furoshiki/style"
)

// UPDATE: レイアウト計算用のヘルパー関数群を追加

// IsWidgetVisible は、ウィジェットが可視状態であるかを安全にチェックします。
// InteractiveStateインターフェースを実装していないウィジェットは、デフォルトで可視(true)とみなします。
func IsWidgetVisible(w component.Widget) bool {
	if is, ok := w.(component.InteractiveState); ok {
		return is.IsVisible()
	}
	return true
}

// getVisibleChildren は、コンテナから表示状態の子ウィジェットのみを抽出し、
// レイアウト計算の対象となるスライスを返します。
func getVisibleChildren(container Container) []component.Widget {
	allChildren := container.GetChildren()
	visibleChildren := make([]component.Widget, 0, len(allChildren))
	for _, child := range allChildren {
		// UPDATE: 型アサーションをヘルパー関数呼び出しに置き換え
		if IsWidgetVisible(child) {
			visibleChildren = append(visibleChildren, child)
		}
	}
	return visibleChildren
}

// GetFlex はウィジェットからFlex値を取得します。実装していない場合は0を返します。
func GetFlex(w component.Widget) int {
	if lpo, ok := w.(component.LayoutPropertiesOwner); ok {
		return lpo.GetLayoutProperties().GetFlex()
	}
	return 0
}

// GetMargin はウィジェットからマージンを取得します。設定されていない場合はゼロ値を返します。
func GetMargin(w component.Widget) style.Insets {
	if sgs, ok := w.(component.StyleGetterSetter); ok {
		// 読み取り専用なのでコピーを避ける
		style := sgs.ReadOnlyStyle()
		if style.Margin != nil {
			return *style.Margin
		}
	}
	return style.Insets{}
}

// GetSize はウィジェットからサイズを取得します。実装していない場合は(0, 0)を返します。
func GetSize(w component.Widget) (int, int) {
	if ss, ok := w.(component.SizeSetter); ok {
		return ss.GetSize()
	}
	return 0, 0
}

// GetMinSize はウィジェットから最小サイズを取得します。実装していない場合は(0, 0)を返します。
func GetMinSize(w component.Widget) (int, int) {
	if mss, ok := w.(component.MinSizeSetter); ok {
		return mss.GetMinSize()
	}
	return 0, 0
}

// GetRequestedPosition はウィジェットから希望配置位置を取得します。実装していない場合は(0, 0)を返します。
func GetRequestedPosition(w component.Widget) (int, int) {
	if ap, ok := w.(component.AbsolutePositioner); ok {
		return ap.GetRequestedPosition()
	}
	return 0, 0
}

// GetLayoutData はウィジェットからレイアウト固有データを取得します。実装していない場合はnilを返します。
func GetLayoutData(w component.Widget) any {
	if lpo, ok := w.(component.LayoutPropertiesOwner); ok {
		return lpo.GetLayoutProperties().GetLayoutData()
	}
	return nil
}

// SetPosition はウィジェットに位置を設定します。
func SetPosition(w component.Widget, x, y int) {
	if ps, ok := w.(component.PositionSetter); ok {
		ps.SetPosition(x, y)
	}
}

// SetSize はウィジェットにサイズを設定します。
func SetSize(w component.Widget, width, height int) {
	if ss, ok := w.(component.SizeSetter); ok {
		ss.SetSize(width, height)
	}
}