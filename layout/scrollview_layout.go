package layout

import "furoshiki/component"

// ScrollViewLayout は、ScrollViewウィジェットのための専用レイアウトマネージャです。
type ScrollViewLayout struct{}

// Layout は、ScrollViewのレイアウトロジックを実行します。
// 2パスレイアウト（計測→配置）のアプローチを取ります。
// NOTE: Layoutインターフェースの変更に伴い、errorを返すようにシグネチャが更新されました。
func (l *ScrollViewLayout) Layout(container Container) error {
	scroller, ok := container.(ScrollViewer)
	if !ok {
		// 【提案1】インターフェースがスリム化したため、containerがScrollViewerを
		// 実装していないケースも考慮し、エラーを返さず単に処理を中断します。
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

	viewX, viewY := scroller.GetPosition()
	viewWidth, viewHeight := scroller.GetSize()
	padding := scroller.GetPadding()

	var scrollBarWidth int
	if vScrollBar != nil {
		// NOTE: vScrollBarはScrollBarWidgetインターフェースです。
		// GetSizeはSizeSetterが持ち、ScrollBarWidgetインターフェース定義に
		// SizeSetterが含まれているため、直接呼び出し可能です。
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
		// 描画せずにコンテンツの本来の高さを知るため、十分な高さを与えて
		// 一時的にレイアウトを計算させ、子の配置から最大高さを割り出す。
		// NOTE: このマジックナンバー的なアプローチは、コンテンツの高さが子の配置に依存し、
		//       かつそれを事前に計算する簡単な方法がない場合に有効な手法です。
		const layoutMeasureHeight = 1_000_000 // 計測用に十分大きな高さを設定
		c.SetSize(potentialContentWidth, layoutMeasureHeight)
		c.SetPosition(0, 0) // レイアウト計算のために一時的に原点に配置
		c.MarkDirty(true)
		c.Update() // content のレイアウトを強制的に実行

		contentPadding := c.GetPadding()
		maxY := 0
		// コンテナ自身の座標は、この一時的なレイアウト計算では(0,0)に設定されているため、
		// 子のY座標と高さから直接最大Y座標を求められます。
		for _, child := range c.GetChildren() {
			// UPDATE: 型アサーションをヘルパー関数に置き換え
			if !IsWidgetVisible(child) {
				continue
			}

			// UPDATE: コンパイラエラー(undefined: GetPosition)を回避するため、
			// この箇所ではヘルパー関数を使わず、明示的な型アサーションに戻します。
			// 他のレイアウトファイルではヘルパーが機能するため、このファイル特有の問題の可能性があります。
			var childY, childH int
			if ps, okPos := child.(component.PositionSetter); okPos {
				_, childY = ps.GetPosition()
			}
			if ss, okSize := child.(component.SizeSetter); okSize {
				_, childH = ss.GetSize()
			}

			bottom := childY + childH
			if bottom > maxY {
				maxY = bottom
			}
		}
		measuredContentHeight = maxY + contentPadding.Bottom
	} else {
		// ケース3: 上記以外のウィジェット (HeightForWiderを実装しない単一ウィジェット)
		// コンテンツ自身の最小の高さを必要な高さとみなします。
		// UPDATE: 型アサーションをヘルパー関数に置き換え
		_, measuredContentHeight = GetMinSize(content)
	}
	scroller.SetContentHeight(measuredContentHeight)

	// --- 2. 配置(Arrange)パス ---
	// 計測した高さに基づき、スクロールバーの要否を決定し、最終的な配置を計算します。
	isVScrollNeeded := measuredContentHeight > contentAreaHeight
	if vScrollBar != nil {
		// NOTE: SetVisibleはInteractiveStateが持ち、ScrollBarWidgetに
		// 含まれているため、直接呼び出し可能です。
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
			// NOTE: 幅が変わったため、コンテナの高さも再計測します。
			const layoutMeasureHeight = 1_000_000
			c.SetSize(finalContentWidth, layoutMeasureHeight)
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
	// UPDATE: 型アサーションをヘルパー関数に置き換え
	SetSize(content, finalContentWidth, measuredContentHeight)
	SetPosition(content, viewX+padding.Left, viewY+padding.Top-int(currentScrollY))

	if isVScrollNeeded && vScrollBar != nil {
		// NOTE: インターフェース定義にPositionSetterとSizeSetterを
		// 含めたので直接呼び出し可能です。
		vScrollBar.SetPosition(viewX+viewWidth-padding.Right-scrollBarWidth, viewY+padding.Top)
		vScrollBar.SetSize(scrollBarWidth, contentAreaHeight)

		contentRatio := float64(contentAreaHeight) / float64(measuredContentHeight)
		scrollRatio := 0.0
		if maxScrollY > 0 {
			scrollRatio = currentScrollY / maxScrollY
		}
		vScrollBar.SetRatios(contentRatio, scrollRatio)
	}
	return nil
}