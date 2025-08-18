package layout

// AbsoluteLayout は、子要素をコンテナ内の絶対座標に基づいて配置します。
// 各子要素の (x, y) 座標がそのままコンテナの左上からのオフセットとして使用されます。
type AbsoluteLayout struct{}

// Layout は AbsoluteLayout のレイアウトロジックを実装します。
func (l *AbsoluteLayout) Layout(container Container) {
	containerX, containerY := container.GetPosition()
	for _, child := range container.GetChildren() {
		// 非表示のコンポーネントはレイアウト処理をスキップ
		if !child.IsVisible() {
			continue
		}
		// 子要素のローカル座標をコンテナのグローバル座標に変換して設定
		childX, childY := child.GetPosition()
		child.SetPosition(containerX+childX, containerY+childY)
	}
}
