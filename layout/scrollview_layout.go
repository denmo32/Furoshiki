package layout

// ScrollViewLayout は、ScrollViewウィジェットのための専用レイアウトマネージャです。
type ScrollViewLayout struct{}

// Layout は、ScrollViewのレイアウトロジックを実行します。
// 2パスレイアウト（計測→配置）のアプローチを取ります。
func (l *ScrollViewLayout) Layout(container Container) {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		return
	}
	content := scroller.GetContentContainer()
	vScrollBar := scroller.GetVScrollBar()

	if content == nil {
		if vScrollBar != nil {
			vScrollBar.SetVisible(false)
		}
		return
	}

	viewX, viewY := scroller.GetPosition()
	viewWidth, viewHeight := scroller.GetSize()
	padding := scroller.GetPadding()

	var scrollBarWidth int
	if vScrollBar != nil {
		scrollBarWidth, _ = vScrollBar.GetSize()
	}

	contentAreaHeight := viewHeight - padding.Top - padding.Bottom
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}

	// --- 1. 計測(Measure)パス ---
	// スクロールバーがないと仮定した幅で、コンテンツが本来必要とする高さを計測します。
	potentialContentWidth := viewWidth - padding.Left - padding.Right
	const veryLargeHeight = 1_000_000 // 計測用に十分大きな高さを設定
	content.SetSize(potentialContentWidth, veryLargeHeight)
	content.SetPosition(0, 0) // レイアウト計算のために一時的に原点に配置
	content.MarkDirty(true)
	content.Update() // content のレイアウトを強制的に実行

	var measuredContentHeight int
	if c, ok := content.(Container); ok {
		// コンテンツがコンテナの場合、子の位置から最大高さを計算
		contentPadding := c.GetPadding()
		maxY := 0
		_, contentY := c.GetPosition()
		for _, child := range c.GetChildren() {
			if !child.IsVisible() {
				continue
			}
			_, childY := child.GetPosition()
			_, childH := child.GetSize()
			bottom := (childY - contentY) + childH
			if bottom > maxY {
				maxY = bottom
			}
		}
		measuredContentHeight = maxY + contentPadding.Bottom
	} else {
		_, measuredContentHeight = content.GetSize()
	}
	scroller.SetContentHeight(measuredContentHeight)

	// --- 2. 配置(Arrange)パス ---
	// 計測した高さに基づき、スクロールバーの要否を決定し、最終的な配置を計算します。
	isVScrollNeeded := measuredContentHeight > contentAreaHeight
	if vScrollBar != nil {
		vScrollBar.SetVisible(isVScrollNeeded)
	}

	finalContentWidth := viewWidth - padding.Left - padding.Right
	if isVScrollNeeded {
		finalContentWidth -= scrollBarWidth
	}
	if finalContentWidth < 0 {
		finalContentWidth = 0
	}

	// スクロールバーの有無で幅が変わった場合、再度レイアウトを実行
	if finalContentWidth != potentialContentWidth {
		content.SetSize(finalContentWidth, veryLargeHeight)
		content.SetPosition(0, 0)
		content.MarkDirty(true)
		content.Update()
	}

	maxScrollY := 0.0
	if measuredContentHeight > contentAreaHeight {
		maxScrollY = float64(measuredContentHeight - contentAreaHeight)
	}
	currentScrollY := scroller.GetScrollY()
	if currentScrollY > maxScrollY {
		currentScrollY = maxScrollY
	}
	if currentScrollY < 0 {
		currentScrollY = 0
	}
	scroller.SetScrollY(currentScrollY)

	// コンテンツの最終的な位置とサイズを設定
	content.SetSize(finalContentWidth, measuredContentHeight)
	content.SetPosition(viewX+padding.Left, viewY+padding.Top-int(currentScrollY))

	if isVScrollNeeded && vScrollBar != nil {
		vScrollBar.SetPosition(viewX+viewWidth-padding.Right-scrollBarWidth, viewY+padding.Top)
		vScrollBar.SetSize(scrollBarWidth, contentAreaHeight)

		contentRatio := float64(contentAreaHeight) / float64(measuredContentHeight)
		scrollRatio := 0.0
		if maxScrollY > 0 {
			scrollRatio = currentScrollY / maxScrollY
		}
		vScrollBar.SetRatios(contentRatio, scrollRatio)
	}
}
