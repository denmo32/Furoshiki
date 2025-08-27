package layout

import (
	"furoshiki/component"
	"image"
)

var _ component.Widget // Dummy var to force import usage

// ScrollViewLayout は、ScrollViewウィジェットのための専用レイアウトマネージャです。
type ScrollViewLayout struct{}

// measureResult は、Measureパスの結果をArrangeパスに渡すための中間構造体です。
type measureResult struct {
	contentSize     image.Point
	isVScrollNeeded bool
	scrollBarWidth  int
}

// measureInternal は、MeasureとArrangeの両方から呼び出される共通の計測ロジックです。
// これにより、ロジックの重複を防ぎ、一貫性を保ちます。
func (l *ScrollViewLayout) measureInternal(scroller ScrollViewer, availableSize image.Point) measureResult {
	content := scroller.GetContentContainer()
	vScrollBar := scroller.GetVScrollBar()
	if content == nil {
		return measureResult{}
	}

	padding := scroller.GetPadding()
	viewWidth := availableSize.X
	viewHeight := availableSize.Y

	// 1. スクロールバーがないと仮定した幅でコンテンツを計測
	contentAvailableWidth := viewWidth - padding.Left - padding.Right
	if contentAvailableWidth < 0 {
		contentAvailableWidth = 0
	}
	// 高さは無制限として渡し、コンテンツが要求する真の高さを取得
	contentDesiredSize := content.Measure(image.Point{X: contentAvailableWidth, Y: 0})

	// 2. スクロールバーの要否を判断
	contentAreaHeight := viewHeight - padding.Top - padding.Bottom
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}
	isVScrollNeeded := contentDesiredSize.Y > contentAreaHeight

	// 3. スクロールバーが必要な場合、幅を考慮して再計測
	scrollBarWidth := 0
	if isVScrollNeeded {
		if vScrollBar != nil {
			// スクロールバー自体の要求サイズを取得
			sbSize := vScrollBar.Measure(image.Point{X: 0, Y: contentAreaHeight})
			scrollBarWidth = sbSize.X
		}
		// スクロールバーの幅を差し引いて再度コンテンツの利用可能幅を計算
		newContentAvailableWidth := viewWidth - padding.Left - padding.Right - scrollBarWidth
		if newContentAvailableWidth < 0 {
			newContentAvailableWidth = 0
		}
		// 幅が変わった場合のみ、コンテンツを再計測
		if newContentAvailableWidth != contentAvailableWidth {
			contentDesiredSize = content.Measure(image.Point{X: newContentAvailableWidth, Y: 0})
		}
	}

	return measureResult{
		contentSize:     contentDesiredSize,
		isVScrollNeeded: isVScrollNeeded,
		scrollBarWidth:  scrollBarWidth,
	}
}

// Measure は、ScrollViewの要求サイズを計算します。
func (l *ScrollViewLayout) Measure(container Container, availableSize image.Point) image.Point {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		// 型アサーションが失敗した場合、コンテナの最小サイズを返す
		minW, minH := container.GetMinSize()
		return image.Point{X: minW, Y: minH}
	}

	// ScrollViewは通常、親から与えられたスペースを埋めるように動作するため、
	// availableSizeをそのまま自身の要求サイズとして返します。
	// コンテンツのサイズによってScrollView自体のサイズが変わるわけではありません。
	// ただし、もし親がサイズを指定しない場合（availableSizeが0の場合）のために、
	// コンテンツの最小幅を考慮した最小サイズを返すこともできます。
	res := l.measureInternal(scroller, availableSize)
	padding := scroller.GetPadding()

	desiredW := res.contentSize.X + padding.Left + padding.Right
	if res.isVScrollNeeded {
		desiredW += res.scrollBarWidth
	}
	// 高さは利用可能な高さに依存するため、ここではavailableSize.Yを尊重する
	return image.Point{X: desiredW, Y: availableSize.Y}
}

// Arrange は、ScrollView内のコンテンツとスクロールバーを配置します。
func (l *ScrollViewLayout) Arrange(container Container, finalBounds image.Rectangle) error {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		return nil
	}
	content := scroller.GetContentContainer()
	vScrollBar := scroller.GetVScrollBar()

	if content == nil {
		if vScrollBar != nil {
			vScrollBar.SetVisible(false)
		}
		return nil
	}

	// NOTE: Measureパスの結果をキャッシュするのが理想的だが、ここでは再計算する
	res := l.measureInternal(scroller, finalBounds.Size())
	scroller.SetContentHeight(res.contentSize.Y)

	if vScrollBar != nil {
		vScrollBar.SetVisible(res.isVScrollNeeded)
	}

	padding := scroller.GetPadding()
	viewHeight := finalBounds.Dy()
	contentAreaHeight := viewHeight - padding.Top - padding.Bottom
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}

	// スクロール位置を正規化
	maxScrollY := 0.0
	if res.contentSize.Y > contentAreaHeight {
		maxScrollY = float64(res.contentSize.Y - contentAreaHeight)
	}
	currentScrollY := scroller.GetScrollY()
	if currentScrollY > maxScrollY {
		currentScrollY = maxScrollY
	}
	if currentScrollY < 0 {
		currentScrollY = 0
	}
	scroller.SetScrollY(currentScrollY)

	// コンテンツの配置
	contentX := finalBounds.Min.X + padding.Left
	contentY := finalBounds.Min.Y + padding.Top - int(currentScrollY)
	contentBounds := image.Rect(contentX, contentY, contentX+res.contentSize.X, contentY+res.contentSize.Y)
	content.SetPosition(contentBounds.Min.X, contentBounds.Min.Y)
	content.SetSize(contentBounds.Dx(), contentBounds.Dy())
	if err := content.Arrange(contentBounds); err != nil {
		return err
	}

	// スクロールバーの配置
	if res.isVScrollNeeded && vScrollBar != nil {
		sbX := finalBounds.Max.X - padding.Right - res.scrollBarWidth
		sbY := finalBounds.Min.Y + padding.Top
		sbBounds := image.Rect(sbX, sbY, sbX+res.scrollBarWidth, sbY+contentAreaHeight)
		vScrollBar.SetPosition(sbBounds.Min.X, sbBounds.Min.Y)
		vScrollBar.SetSize(sbBounds.Dx(), sbBounds.Dy())

		contentRatio := float64(contentAreaHeight) / float64(res.contentSize.Y)
		scrollRatio := 0.0
		if maxScrollY > 0 {
			scrollRatio = currentScrollY / maxScrollY
		}
		vScrollBar.SetRatios(contentRatio, scrollRatio)

		if err := vScrollBar.Arrange(sbBounds); err != nil {
			return err
		}
	}

	return nil
}