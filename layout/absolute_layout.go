package layout

import "furoshiki/component"

// AbsoluteLayout は、子要素をコンテナ内の指定された相対座標に基づいて配置します。
type AbsoluteLayout struct{}

// Layout は AbsoluteLayout のレイアウトロジックを実装します。
// NOTE: Layoutインターフェースの変更に伴い、errorを返すようにシグネチャが更新されました。
func (l *AbsoluteLayout) Layout(container Container) error {
	containerX, containerY := container.GetPosition()
	padding := container.GetPadding()

	for _, child := range container.GetChildren() {
		// 【提案1】型アサーションの追加: 可視状態はWidgetインターフェースではなく
		// InteractiveStateインターフェースが持つため、型アサーションでチェックします。
		// 実装していないウィジェットはデフォルトで可視(true)とみなします。
		isVisible := true
		if is, ok := child.(component.InteractiveState); ok {
			isVisible = is.IsVisible()
		}
		if !isVisible {
			continue
		}

		var requestedX, requestedY int
		// child が component.AbsolutePositioner インターフェースを実装しているかチェックします。
		if pr, ok := child.(component.AbsolutePositioner); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		finalX := containerX + padding.Left + requestedX
		finalY := containerY + padding.Top + requestedY

		// 【提案1】型アサーションの追加: SetPositionはPositionSetterインターフェースが持つため、
		// 型アサーションを行い、実装しているウィジェットのみ位置を設定します。
		if ps, ok := child.(component.PositionSetter); ok {
			ps.SetPosition(finalX, finalY)
		}
	}
	return nil
}