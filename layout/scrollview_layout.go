package layout

import "furoshiki/component"

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
	if potentialContentWidth < 0 {
		potentialContentWidth = 0
	}

	var measuredContentHeight int

	// content の種類に応じて、最適な方法で高さを計測します。
	if hw, ok := content.(component.HeightForWider); ok {
		// ケース1: HeightForWiderを実装するウィジェット (例: 折り返し可能なLabel)
		// 幅から直接、必要な高さを計算できるため、最も効率的です。
		measuredContentHeight = hw.GetHeightForWidth(potentialContentWidth)
	} else if c, ok := content.(Container); ok {
		// ケース2: コンテナウィジェット (例: VStack, HStack)
		// 実際にレイアウトを実行させて、子要素の配置から最終的な高さを割り出します。
		const veryLargeHeight = 1_000_000 // 計測用に十分大きな高さを設定
		c.SetSize(potentialContentWidth, veryLargeHeight)
		c.SetPosition(0, 0) // レイアウト計算のために一時的に原点に配置
		c.MarkDirty(true)
		c.Update() // content のレイアウトを強制的に実行

		contentPadding := c.GetPadding()
		maxY := 0
		// コンテナ自身の座標は、この一時的なレイアウト計算では(0,0)に設定されているため、
		// 子のY座標と高さから直接最大Y座標を求められます。
		for _, child := range c.GetChildren() {
			if !child.IsVisible() {
				continue
			}
			_, childY := child.GetPosition()
			_, childH := child.GetSize()
			bottom := childY + childH
			if bottom > maxY {
				maxY = bottom
			}
		}
		measuredContentHeight = maxY + contentPadding.Bottom
	} else {
		// ケース3: 上記以外のウィジェット (HeightForWiderを実装しない単一ウィジェット)
		// コンテンツ自身の最小の高さを必要な高さとみなします。
		_, measuredContentHeight = content.GetMinSize()
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
		// 再計測が必要なのは、幅に依存して高さが変わるウィジェットのみ
		if hw, ok := content.(component.HeightForWider); ok {
			measuredContentHeight = hw.GetHeightForWidth(finalContentWidth)
		} else if c, ok := content.(Container); ok {
			const veryLargeHeight = 1_000_000
			c.SetSize(finalContentWidth, veryLargeHeight)
			c.SetPosition(0, 0)
			c.MarkDirty(true)
			c.Update()
		}
		// NOTE: コンテナ内の子の再計算は c.Update() で行われるため、ここでは不要
		scroller.SetContentHeight(measuredContentHeight)
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