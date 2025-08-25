package layout

// ScrollViewLayout は、ScrollViewウィジェットのための専用レイアウトマネージャです。
type ScrollViewLayout struct{}

// Layout は、ScrollViewのレイアウトロジックを実行します。
func (l *ScrollViewLayout) Layout(container Container) {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		// このパスは、ScrollView.Update の修正により通らなくなるはずです。
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
	potentialContentWidth := viewWidth - padding.Left - padding.Right

	const veryLargeHeight = 1_000_000
	content.SetSize(potentialContentWidth, veryLargeHeight)
	content.SetPosition(0, 0) // レイアウト計算のために一時的に原点に配置

	// content のレイアウトを強制的に実行させる
	content.MarkDirty(true)
	content.Update()

	var measuredContentHeight int
	if c, ok := content.(Container); ok {
		contentPadding := c.GetPadding()
		maxY := 0
		_, contentY := c.GetPosition() // contentの現在位置（一時的に0のはず）
		for _, child := range c.GetChildren() {
			if !child.IsVisible() {
				continue
			}
			_, childY := child.GetPosition()
			_, childH := child.GetSize()
			bottom := (childY - contentY) + childH // 子の下端の相対座標
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

	// 幅が変わった場合、再度レイアウトを実行
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