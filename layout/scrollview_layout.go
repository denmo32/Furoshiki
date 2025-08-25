package layout

// 【修正】widgetへの依存を完全に削除
// import "furoshiki/component" // このファイルでは不要になった
// import "furoshiki/widget"

// ScrollViewLayout は、ScrollViewウィジェットのための専用レイアウトマネージャです。
type ScrollViewLayout struct{}

// Layout は、ScrollViewのレイアウトロジックを実行します。
func (l *ScrollViewLayout) Layout(container Container) {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		return
	}
	content := scroller.GetContentContainer()
	vScrollBar := scroller.GetVScrollBar()
	if content == nil {
		vScrollBar.SetVisible(false)
		return
	}

	viewX, viewY := scroller.GetPosition()
	viewWidth, viewHeight := scroller.GetSize()
	padding := scroller.GetPadding()
	scrollBarWidth, _ := vScrollBar.GetSize()
	contentAreaHeight := viewHeight - padding.Top - padding.Bottom
	if contentAreaHeight < 0 {
		contentAreaHeight = 0
	}

	// --- 1. 計測(Measure)パス ---
	potentialContentWidth := viewWidth - padding.Left - padding.Right

	const veryLargeHeight = 1_000_000
	content.SetSize(potentialContentWidth, veryLargeHeight)
	content.SetPosition(0, 0)

	// 【重要】子のレイアウトマネージャを直接呼ばず、子のUpdateを呼ぶことで
	// 子のレイアウトを自律的に実行させる
	content.Update()

	var measuredContentHeight int
	if c, ok := content.(Container); ok {
		contentPadding := c.GetPadding()
		for _, child := range c.GetChildren() {
			if !child.IsVisible() {
				continue
			}
			_, y := child.GetPosition()
			_, h := child.GetSize()
			bottom := y + h
			if bottom > measuredContentHeight {
				measuredContentHeight = bottom
			}
		}
		measuredContentHeight += contentPadding.Bottom
	} else {
		_, measuredContentHeight = content.GetSize()
	}
	scroller.SetContentHeight(measuredContentHeight)

	// --- 2. 配置(Arrange)パス ---
	isVScrollNeeded := measuredContentHeight > contentAreaHeight
	vScrollBar.SetVisible(isVScrollNeeded)

	finalContentWidth := viewWidth - padding.Left - padding.Right
	if isVScrollNeeded {
		finalContentWidth -= scrollBarWidth
	}

	if finalContentWidth != potentialContentWidth {
		content.SetSize(finalContentWidth, veryLargeHeight)
		content.SetPosition(0, 0)
		content.Update() // 幅が変わったので再度レイアウト
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

	content.SetSize(finalContentWidth, measuredContentHeight)
	content.SetPosition(viewX+padding.Left, viewY+padding.Top-int(currentScrollY))

	if isVScrollNeeded {
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