package layout

// AbsoluteLayout は、子要素をコンテナ内の指定された相対座標に基づいて配置します。
type AbsoluteLayout struct{}

// Layout は AbsoluteLayout のレイアウトロジックを実装します。
// NOTE: Layoutインターフェースの変更に伴い、errorを返すようにシグネチャが更新されました。
func (l *AbsoluteLayout) Layout(container Container) error {
	containerX, containerY := container.GetPosition()
	padding := container.GetPadding()

	for _, child := range container.GetChildren() {
		// UPDATE: 型アサーションをIsWidgetVisibleヘルパー関数に置き換え
		if !IsWidgetVisible(child) {
			continue
		}

		// UPDATE: 型アサーションをGetRequestedPositionヘルパー関数に置き換え
		requestedX, requestedY := GetRequestedPosition(child)

		finalX := containerX + padding.Left + requestedX
		finalY := containerY + padding.Top + requestedY

		// UPDATE: 型アサーションをSetPositionヘルパー関数に置き換え
		SetPosition(child, finalX, finalY)
	}
	return nil
}
