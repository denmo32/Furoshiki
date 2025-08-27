package layout

import (
	"furoshiki/component"
	"image"
)

// AbsoluteLayout は、子要素をコンテナ内の指定された相対座標に基づいて配置します。
// ZStackなどで利用されます。
type AbsoluteLayout struct{}

// Measure は、AbsoluteLayoutの要求サイズを計算します。
// 要求サイズは、すべての子要素をその希望位置に配置したときに、
// それらをちょうど内包できる大きさになります。
func (l *AbsoluteLayout) Measure(container Container, availableSize image.Point) image.Point {
	padding := container.GetPadding()
	maxW, maxH := 0, 0

	// パディングを引いた、子が利用できる真のサイズ
	innerAvailableW := availableSize.X - padding.Left - padding.Right
	if innerAvailableW < 0 {
		innerAvailableW = 0
	}
	innerAvailableH := availableSize.Y - padding.Top - padding.Bottom
	if innerAvailableH < 0 {
		innerAvailableH = 0
	}
	childAvailableSize := image.Point{X: innerAvailableW, Y: innerAvailableH}

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		// AbsoluteLayoutでは、子は他の子の影響を受けずにサイズを決定できるため、
		// コンテナの利用可能サイズをそのまま渡して計測させます。
		childDesiredSize := child.Measure(childAvailableSize)

		var requestedX, requestedY int
		if pr, ok := child.(component.AbsolutePositioner); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		// 子が占める右下の座標を計算
		childMaxX := requestedX + childDesiredSize.X
		childMaxY := requestedY + childDesiredSize.Y

		// 全体で必要なサイズを更新
		if childMaxX > maxW {
			maxW = childMaxX
		}
		if childMaxY > maxH {
			maxH = childMaxY
		}
	}

	// パディングを足して、コンテナ全体の要求サイズを決定
	desiredW := maxW + padding.Left + padding.Right
	desiredH := maxH + padding.Top + padding.Bottom

	return image.Point{X: desiredW, Y: desiredH}
}

// Arrange は、AbsoluteLayout内の子要素を最終的な位置に配置します。
func (l *AbsoluteLayout) Arrange(container Container, finalBounds image.Rectangle) error {
	padding := container.GetPadding()

	// パディングを引いた、子が利用できる真の領域
	innerBoundsW := finalBounds.Dx() - padding.Left - padding.Right
	if innerBoundsW < 0 {
		innerBoundsW = 0
	}
	innerBoundsH := finalBounds.Dy() - padding.Top - padding.Bottom
	if innerBoundsH < 0 {
		innerBoundsH = 0
	}
	childAvailableSize := image.Point{X: innerBoundsW, Y: innerBoundsH}

	for _, child := range container.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		// NOTE: Measure/Arrange パターンでは、Arrangeパスで子のサイズを再計算するのが一般的です。
		// 本来はMeasureパスの結果をキャッシュすることで最適化しますが、ここではシンプルに再計測します。
		childDesiredSize := child.Measure(childAvailableSize)

		var requestedX, requestedY int
		if pr, ok := child.(component.AbsolutePositioner); ok {
			requestedX, requestedY = pr.GetRequestedPosition()
		}

		// コンテナの左上座標とパディング、子の希望位置を基に最終的な位置を計算
		finalX := finalBounds.Min.X + padding.Left + requestedX
		finalY := finalBounds.Min.Y + padding.Top + requestedY

		child.SetPosition(finalX, finalY)
		child.SetSize(childDesiredSize.X, childDesiredSize.Y)

		// 子の描画領域を計算し、子のArrangeを再帰的に呼び出す
		childBounds := image.Rect(finalX, finalY, finalX+childDesiredSize.X, finalY+childDesiredSize.Y)
		if err := child.Arrange(childBounds); err != nil {
			// エラーが発生した場合、それを上に伝播させます。
			// 必要であれば、エラーを集約する処理をここに追加することも可能です。
			return err
		}
	}
	return nil
}