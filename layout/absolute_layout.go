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
		if !child.IsVisible() {
			continue
		}

		var requestedX, requestedY int
		// child が component.AbsolutePositioner インターフェースを実装しているかチェックします。
		if pr, ok := child.(component.AbsolutePositioner); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		finalX := containerX + padding.Left + requestedX
		finalY := containerY + padding.Top + requestedY
		child.SetPosition(finalX, finalY)
	}
	return nil
}