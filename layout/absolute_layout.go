package layout

// positionRequester は、AbsoluteLayoutが子の相対位置を取得するために使用するインターフェースです。
type positionRequester interface {
	GetRequestedPosition() (int, int)
}

// AbsoluteLayout は、子要素をコンテナ内の指定された相対座標に基づいて配置します。
type AbsoluteLayout struct{}

// Layout は AbsoluteLayout のレイアウトロジックを実装します。
func (l *AbsoluteLayout) Layout(container Container) {
	containerX, containerY := container.GetPosition()
	padding := container.GetPadding()

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		var requestedX, requestedY int
		if pr, ok := child.(positionRequester); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		finalX := containerX + padding.Left + requestedX
		finalY := containerY + padding.Top + requestedY
		child.SetPosition(finalX, finalY)
	}
}
