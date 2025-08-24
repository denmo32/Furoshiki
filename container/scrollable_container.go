package container

import (
	"furoshiki/component"
	"furoshiki/event"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollableContainer は縦方向スクロール機能を持つコンテナです
type ScrollableContainer struct {
	*Container

	// スクロール関連のプロパティ
	contentHeight int  // コンテンツ全体の高さ
	scrollY       int  // Y方向のスクロール位置
	canScrollY    bool // スクロール可能かどうか
}

// NewScrollableContainer は新しいScrollableContainerを作成します
func NewScrollableContainer() *ScrollableContainer {
	// NewContainerBuilderを使ってコンテナを作成
	builder := NewContainerBuilder()
	container, _ := builder.Build()

	sc := &ScrollableContainer{
		Container: container,
		scrollY:   0,
	}

	// コンテナのクリッピングを有効にする
	sc.SetClipsChildren(true)

	return sc
}

// AddChild は子ウィジェットを追加し、コンテンツサイズを更新します
func (sc *ScrollableContainer) AddChild(child component.Widget) {
	sc.Container.AddChild(child)
	sc.updateContentHeight()
}

// RemoveChild は子ウィジェットを削除し、コンテンツサイズを更新します
func (sc *ScrollableContainer) RemoveChild(child component.Widget) {
	sc.Container.RemoveChild(child)
	sc.updateContentHeight()
}

// updateContentHeight はコンテンツ全体の高さを計算します
func (sc *ScrollableContainer) updateContentHeight() {
	maxHeight := 0

	for _, child := range sc.GetChildren() {
		if child.IsVisible() {
			_, y := child.GetPosition()
			_, height := child.GetSize()

			// 子要素の下端を計算
			bottom := y + height
			if bottom > maxHeight {
				maxHeight = bottom
			}
		}
	}

	sc.contentHeight = maxHeight

	// スクロール可能かどうかを判定
	_, containerHeight := sc.GetSize()
	sc.canScrollY = sc.contentHeight > containerHeight

	// スクロール位置を有効範囲内に収める
	sc.constrainScrollPosition()
}

// constrainScrollPosition はスクロール位置を有効範囲内に収めます
func (sc *ScrollableContainer) constrainScrollPosition() {
	if !sc.canScrollY {
		sc.scrollY = 0
		return
	}

	_, containerHeight := sc.GetSize()
	maxScrollY := sc.contentHeight - containerHeight

	if sc.scrollY < 0 {
		sc.scrollY = 0
	} else if sc.scrollY > maxScrollY {
		sc.scrollY = maxScrollY
	}
}

// SetScrollPosition はスクロール位置を設定します
func (sc *ScrollableContainer) SetScrollPosition(y int) {
	oldY := sc.scrollY
	sc.scrollY = y
	sc.constrainScrollPosition()

	// スクロール位置が変更された場合のみ再描画
	if oldY != sc.scrollY {
		sc.MarkDirty(false)
	}
}

// GetScrollPosition は現在のスクロール位置を返します
func (sc *ScrollableContainer) GetScrollPosition() int {
	return sc.scrollY
}

// GetContentHeight はコンテンツ全体の高さを返します
func (sc *ScrollableContainer) GetContentHeight() int {
	return sc.contentHeight
}

// CanScroll はスクロール可能かどうかを返します
func (sc *ScrollableContainer) CanScroll() bool {
	return sc.canScrollY
}

// Draw はスクロール位置を考慮して子要素を描画します
func (sc *ScrollableContainer) Draw(screen *ebiten.Image) {
	if !sc.IsVisible() {
		return
	}

	// コンテナ自身の背景を描画
	x, y := sc.GetPosition()
	width, height := sc.GetSize()
	component.DrawStyledBackground(screen, x, y, width, height, sc.GetStyle())

	// クリッピングを有効にして子要素を描画
	sc.drawChildrenWithClipping(screen)
}

// drawChildrenWithClipping はクリッピングありで子要素を描画します
func (sc *ScrollableContainer) drawChildrenWithClipping(screen *ebiten.Image) {
	containerX, containerY := sc.GetPosition()
	containerWidth, containerHeight := sc.GetSize()

	// 描画領域がない場合は何もしない
	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	// オフスクリーン画像の準備
	if sc.offscreenImage == nil ||
		sc.offscreenImage.Bounds().Dx() != containerWidth ||
		sc.offscreenImage.Bounds().Dy() != containerHeight {
		if sc.offscreenImage != nil {
			sc.offscreenImage.Deallocate()
		}
		sc.offscreenImage = ebiten.NewImage(containerWidth, containerHeight)
	}
	sc.offscreenImage.Clear()

	// コンテナ自身の背景をオフスクリーン画像に描画
	component.DrawStyledBackground(sc.offscreenImage, 0, 0, containerWidth, containerHeight, sc.GetStyle())

	// 子要素をスクロール位置を考慮してオフスクリーン画像に描画
	for _, child := range sc.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		// 元の位置を保存
		originalX, originalY := child.GetPosition()

		// スクロールを考慮した相対位置を計算
		relativeX := originalX - containerX
		relativeY := originalY - containerY - sc.scrollY

		// 一時的に位置を設定して描画
		child.SetPosition(relativeX, relativeY)
		child.Draw(sc.offscreenImage)

		// 位置を元に戻す
		child.SetPosition(originalX, originalY)
	}

	// 完成したオフスクリーン画像をスクリーンに描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(containerX), float64(containerY))
	screen.DrawImage(sc.offscreenImage, opts)
}

// HandleEvent はマウスホイールイベントを処理します
func (sc *ScrollableContainer) HandleEvent(e *event.Event) {
	// 親クラスのイベント処理を呼び出す
	sc.Container.HandleEvent(e)

	// イベントが既に処理済みの場合は何もしない
	if e.Handled {
		return
	}

	// マウスホイールイベントを処理
	if e.Type == event.MouseScroll && sc.canScrollY {
		// マウスホイールのスクロール量を取得
		scrollAmount := int(e.ScrollY * 30) // スクロール速度を調整

		// 新しいスクロール位置を計算
		newScrollY := sc.scrollY + scrollAmount

		// スクロール位置を設定
		sc.SetScrollPosition(newScrollY)

		// イベントを処理済みとしてマーク
		e.Handled = true
	}
}

// Update はコンテナの状態を更新します
func (sc *ScrollableContainer) Update() {
	sc.Container.Update()

	// コンテナのサイズが変更された場合にコンテンツサイズを再計算
	if sc.NeedsRelayout() {
		sc.updateContentHeight()
	}
}
