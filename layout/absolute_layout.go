package layout

// [追加] AbsoluteLayoutが子の相対位置を取得するために使用するインターフェース。
// これにより、具体的なウィジェット型への依存を避けつつ、必要な機能にアクセスできます。
type positionRequester interface {
	GetRequestedPosition() (int, int)
}

// AbsoluteLayout は、子要素をコンテナ内の指定された相対座標に基づいて配置します。
// 各子要素のビルダーで設定された Position(x, y) が、コンテナの左上からのオフセットとして使用されます。
// [改善] レイアウト計算のたびに位置がずれていくバグを修正しました。
type AbsoluteLayout struct{}

// Layout は AbsoluteLayout のレイアウトロジックを実装します。
// 親コンテナの位置を基準に、子の要求された相対位置から絶対位置を計算し、設定します。
func (l *AbsoluteLayout) Layout(container Container) {
	containerX, containerY := container.GetPosition()
	padding := container.GetPadding()

	for _, child := range container.GetChildren() {
		// 非表示のコンポーネントはレイアウト処理をスキップ
		if !child.IsVisible() {
			continue
		}

		// 子が要求された相対位置を持っているか、型アサーションで確認します。
		var requestedX, requestedY int
		if pr, ok := child.(positionRequester); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		// 子要素の最終的な絶対座標を、コンテナの絶対座標、パディング、子の要求相対座標から計算します。
		finalX := containerX + padding.Left + requestedX
		finalY := containerY + padding.Top + requestedY
		child.SetPosition(finalX, finalY)
	}
}