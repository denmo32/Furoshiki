package container

import (
	"furoshiki/component"
	"furoshiki/event"
)

// ScrollableContainer は縦方向スクロール機能を持つコンテナです
type ScrollableContainer struct {
	*Container

	// スクロール関連のプロパティ
	contentHeight int  // コンテンツ全体の高さ
	scrollY       int  // Y方向のスクロール位置
	canScrollY    bool // スクロール可能かどうか
}

// 【改善】コンパイル時にScrollerインターフェースの実装を検証します。
var _ Scroller = (*ScrollableContainer)(nil)

// NewScrollableContainer は新しいScrollableContainerを作成します
func NewScrollableContainer() *ScrollableContainer {
	// 【改善】ビルダーを使わずに直接構築することで、`self`参照を正しく設定します。
	sc := &ScrollableContainer{
		scrollY: 0,
	}

	// 埋め込むContainerのインスタンスを生成します。
	container := &Container{
		children: make([]component.Widget, 0),
	}
	// LayoutableWidgetを初期化し、`self`として`sc` (*ScrollableContainer) 自身を渡します。
	// これにより、イベントやヒットテストが正しく ScrollableContainer を参照するようになります。
	container.LayoutableWidget = component.NewLayoutableWidget()
	container.Init(sc)
	sc.Container = container

	// コンテナのクリッピングを有効にする
	sc.SetClipsChildren(true)

	return sc
}

// 【新規追加】 GetScrollOffset はScrollerインターフェースを実装します。
// 埋め込まれたContainerの描画ロジックは、このメソッドを呼び出して
// 子要素を描画する際のY座標オフセット（スクロール量）を取得します。
func (sc *ScrollableContainer) GetScrollOffset() (x, y int) {
	// このコンテナは垂直スクロールのみなので、Xオフセットは0です。
	// Yオフセットは、現在のスクロール位置を反転させた値です。
	return 0, -sc.scrollY
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

	// GetChildrenは埋め込まれたContainerのメソッドを呼び出します
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

// 【改善】Drawメソッドは不要になりました。
// 埋め込まれたContainerのDrawメソッドが自動的に呼び出されます。
// その際、Scrollerインターフェースを通じてGetScrollOffsetが呼び出され、
// スクロールが適用されたクリッピング描画が実行されます。
/*
func (sc *ScrollableContainer) Draw(screen *ebiten.Image) {
	// ... (このメソッド全体が不要になる) ...
}
*/

// 【改善】drawChildrenWithClippingメソッドも不要になりました。
// このロジックは汎用化され、Container.drawWithClippingに統合されました。
/*
func (sc *ScrollableContainer) drawChildrenWithClipping(screen *ebiten.Image) {
    // ... (このメソッド全体が不要になる) ...
}
*/

// HandleEvent はマウスホイールイベントを処理します
func (sc *ScrollableContainer) HandleEvent(e *event.Event) {
	// 親クラス(Container)のイベント処理を呼び出し、イベントバブリングを機能させます。
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
	// 埋め込まれたContainerのUpdateを呼び出します。
	// これによりレイアウト計算や子の更新が実行されます。
	sc.Container.Update()

	// コンテナのサイズが変更された場合にコンテンツサイズを再計算
	if sc.NeedsRelayout() {
		sc.updateContentHeight()
	}
}