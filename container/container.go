package container

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/layout"
	"log"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
)

// 【新規追加】 Scroller は、Containerがクリッピング描画時にスクロールオフセットを
// 取得するために使用するインターフェースです。
// ScrollableContainerのような、スクロール機能を持つウィジェットがこのインターフェースを実装します。
type Scroller interface {
	GetScrollOffset() (x, y int)
}

// Containerは子Widgetを保持し、レイアウトを管理するコンポーネントです。
// component.Containerインターフェースを実装します。
type Container struct {
	*component.LayoutableWidget
	children []component.Widget
	layout   layout.Layout
	warned   bool // サイズ警告を一度だけ出すためのフラグ

	// [新規追加] クリッピング関連フィールド
	clipsChildren  bool          // 子要素をクリッピングするかどうか
	offscreenImage *ebiten.Image // クリッピング描画用のオフスクリーンバッファ
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)

// [新規追加]
// SetClipsChildren はコンテナのクリッピング動作を設定します。
func (c *Container) SetClipsChildren(clips bool) {
	if c.clipsChildren != clips {
		c.clipsChildren = clips
		// クリッピング設定の変更は再描画を要求します
		c.MarkDirty(false)
	}
}

// GetLayout はコンテナが使用しているレイアウトを返します。
func (c *Container) GetLayout() layout.Layout {
	return c.layout
}

// SetLayout はコンテナが使用するレイアウトを設定し、再レイアウトを要求します。
func (c *Container) SetLayout(layout layout.Layout) {
	c.layout = layout
	c.MarkDirty(true)
}

// Updateはコンテナと子要素の状態を更新します。
// component.LayoutableWidgetのUpdateをオーバーライドして、レイアウト計算と子の更新を追加します。
// このメソッドはUIツリーのルートから毎フレーム再帰的に呼び出されます。
func (c *Container) Update() {
	if !c.IsVisible() {
		return
	}

	// ルートコンテナのサイズに関する警告を一度だけチェックします。
	c.checkSizeWarning()

	// このコンテナ、またはその子孫のいずれかがダーティな場合、処理を行います。
	if c.IsDirty() {
		// NeedsRelayout() は、ダーティレベルが「レイアウト再計算」を要求しているかチェックします。
		// これにより、ホバー状態の変更など、再描画のみを要求するダーティ状態では
		// 無駄なレイアウト計算が実行されないようになります。
		if c.NeedsRelayout() {
			if c.layout != nil {
				// レイアウト計算は複雑なため、予期せぬ状況でパニックする可能性があります。
				// このdeferは、特定のレイアウト実装のバグがアプリケーション全体をクラッシュさせるのを防ぎます。
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic during layout calculation: %v\n%s", r, debug.Stack())
					}
				}()
				c.layout.Layout(c)
			}
		}
		// レイアウト計算の有無にかかわらず、このフレームで処理されたダーティフラグはクリアします。
		// これにより、次のフレームで不要な処理が走るのを防ぎます。
		c.ClearDirty()
	}

	// すべての子ウィジェットのUpdateを再帰的に呼び出します。
	// これにより、UIツリー内のダーティなコンポーネントがすべて更新されることが保証されます。
	for _, child := range c.children {
		child.Update()
	}
}

// checkSizeWarning はコンテナのサイズに関する警告を出力します。
func (c *Container) checkSizeWarning() {
	// すでに警告を出している場合は何もしません。
	if c.warned {
		return
	}

	// Flexが設定されておらず、かつサイズが両方0のルートコンテナは描画されない可能性が高いため警告します。
	if c.GetFlex() == 0 {
		width, height := c.GetSize()
		if width == 0 && height == 0 {
			// 親がいない（ルートコンテナ）の場合のみ警告します。
			if c.GetParent() == nil {
				fmt.Printf("Warning: Root container created with no flex and zero size (width=0, height=0). It may not be visible.\n")
				c.warned = true // 警告を出したことを記録し、再表示を防ぎます。
			}
		}
	}
}

// Drawはコンテナ自身と、そのすべての子を描画します。
// component.LayoutableWidgetのDrawをオーバーライドします。
// [ロジック更新] クリッピングが有効な場合、オフスクリーンレンダリングを行います。
func (c *Container) Draw(screen *ebiten.Image) {
	// コンテナが非表示の場合、自身も子も描画しません。
	if !c.IsVisible() {
		return
	}

	// --- [ロジック更新] クリッピングが有効な場合の処理 ---
	if c.clipsChildren {
		c.drawWithClipping(screen)
	} else {
		// --- 従来の描画処理（クリッピングなし） ---
		c.drawWithoutClipping(screen)
	}
}

// drawWithoutClipping は、クリッピングを行わずにコンテナと子を描画します。
func (c *Container) drawWithoutClipping(screen *ebiten.Image) {
	// まずコンテナ自身の背景などを描画ヘルパーで描画します。
	x, y := c.GetPosition()
	width, height := c.GetSize()
	component.DrawStyledBackground(screen, x, y, width, height, c.GetStyle())

	// 次にすべての子ウィジェットを描画します。
	for _, child := range c.children {
		child.Draw(screen)
	}
}

// 【改善】drawWithClipping は、オフスクリーンバッファを利用してクリッピングを行いながら描画します。
// Scrollerインターフェースをチェックすることで、ScrollableContainerのようなウィジェットの
// スクロールオフセットを汎用的に扱うことができます。
func (c *Container) drawWithClipping(screen *ebiten.Image) {
	containerX, containerY := c.GetPosition()
	containerWidth, containerHeight := c.GetSize()

	// 描画領域がない場合は何もしません
	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	// 1. オフスクリーン画像の準備
	// サイズが変更されていたら、新しいサイズの画像を再生成します。
	if c.offscreenImage == nil || c.offscreenImage.Bounds().Dx() != containerWidth || c.offscreenImage.Bounds().Dy() != containerHeight {
		if c.offscreenImage != nil {
			c.offscreenImage.Deallocate()
		}
		c.offscreenImage = ebiten.NewImage(containerWidth, containerHeight)
	}
	c.offscreenImage.Clear() // 前のフレームの描画をクリア

	// 2. コンテナ自身の背景をオフスクリーン画像に描画
	// オフスクリーン画像への描画なので、描画位置は(0,0)から開始します。
	component.DrawStyledBackground(c.offscreenImage, 0, 0, containerWidth, containerHeight, c.GetStyle())

	// 【改善点】コンテナ自身がScrollerインターフェースを実装しているかチェックします。
	// これにより、具象型に依存することなくスクロール機能に対応できます。
	var scrollOffsetX, scrollOffsetY int
	if scroller, ok := any(c).(Scroller); ok {
		scrollOffsetX, scrollOffsetY = scroller.GetScrollOffset()
	}

	// 3. 子要素をオフスクリーン画像に描画
	for _, child := range c.children {
		// 子の描画メソッドは、ウィジェットの絶対座標に基づいて描画を行います。
		// オフスクリーン画像はコンテナの左上を(0,0)としているため、
		// 子の描画座標を「コンテナ基準の相対座標」に一時的に変換する必要があります。
		originalX, originalY := child.GetPosition()

		// 【改善点】スクロールオフセットを考慮して相対座標を計算します。
		// Scrollerでない場合はオフセットが(0,0)なので、通常のクリッピングとして機能します。
		relativeX := originalX - containerX + scrollOffsetX
		relativeY := originalY - containerY + scrollOffsetY

		// 座標を一時的に設定して描画
		child.SetPosition(relativeX, relativeY)
		child.Draw(c.offscreenImage)

		// 座標を元に戻す
		child.SetPosition(originalX, originalY)
	}

	// 4. 完成したオフスクリーン画像を、スクリーン上の正しい位置に描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(containerX), float64(containerY))
	screen.DrawImage(c.offscreenImage, opts)
}

// HitTest は、指定された座標がコンテナまたはその子のいずれかにヒットするかをテストします。
// component.LayoutableWidgetのHitTestをオーバーライドします。
func (c *Container) HitTest(x, y int) component.Widget {
	// コンテナが非表示の場合はヒットしません。
	if !c.IsVisible() {
		return nil
	}
	// 描画順とは逆に、最前面の子（スライスの末尾）からヒットテストを行います。
	for i := len(c.children) - 1; i >= 0; i-- {
		child := c.children[i]
		// 非表示の子はヒットテストの対象外です。
		if !child.IsVisible() {
			continue
		}
		if target := child.HitTest(x, y); target != nil {
			return target // ヒットした子を返します。
		}
	}
	// どの子にもヒットしなかった場合、コンテナ自身がヒットするかテストします。
	// これにより、子の間の隙間やパディング部分でコンテナがイベントを受け取ることができます。
	if c.LayoutableWidget.HitTest(x, y) != nil {
		return c // コンテナ自身を返します。
	}
	return nil
}

// Cleanup は、コンテナとすべての子ウィジェットのリソースを解放します。
// UIツリーからコンテナが削除されるときや、アプリケーション終了時に呼び出されるべきです。
// [ロジック更新] オフスクリーン画像のリソースも解放します。
func (c *Container) Cleanup() {
	// まず、すべての子ウィジェットのクリーンアップを再帰的に呼び出します。
	for _, child := range c.children {
		child.Cleanup()
	}
	// 子のリストをクリアします。
	c.children = nil

	// [新規追加] オフスクリーン画像が存在すれば、そのリソースを解放します。
	if c.offscreenImage != nil {
		c.offscreenImage.Deallocate()
		c.offscreenImage = nil
	}

	// 最後に、埋め込まれたLayoutableWidget自身のクリーンアップ処理（イベントハンドラのクリアなど）を呼び出します。
	c.LayoutableWidget.Cleanup()
}

// --- Container Methods (from component.Container interface) ---

// detachChildは、親子関係のみを解消する内部ヘルパーです。
// Cleanupを呼び出さないため、ウィジェットをツリー間で安全に移動させるために使用します。
func (c *Container) detachChild(child component.Widget) bool {
	if child == nil {
		return false
	}
	for i, currentChild := range c.children {
		if currentChild == child {
			c.children = append(c.children[:i], c.children[i+1:]...)
			child.SetParent(nil)
			return true
		}
	}
	return false
}

// AddChild はコンテナに子ウィジェットを追加します。
func (c *Container) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	// 既に親が存在する場合は、その親から子をデタッチ（親子関係の解消のみ）します。
	// これにより、ウィジェットをUIツリー間で安全に「移動」させることができます。
	if oldParent := child.GetParent(); oldParent != nil {
		// Container型にキャストしてdetachChildを呼ぶことで、Cleanupが呼ばれるのを防ぐ
		if container, ok := oldParent.(*Container); ok {
			container.detachChild(child)
		}
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.MarkDirty(true)
}

// RemoveChild はコンテナから子ウィジェットを削除し、リソースを解放します。
func (c *Container) RemoveChild(child component.Widget) {
	// detachChildで親子関係を解消し、成功した場合のみCleanupを呼び出します。
	if c.detachChild(child) {
		// 子のリソースを解放します。
		child.Cleanup()
		// コンテナの再レイアウトを要求します。
		c.MarkDirty(true)
	}
}

// GetChildren はコンテナが保持するすべての子ウィジェットのスライスを返します。
func (c *Container) GetChildren() []component.Widget {
	return c.children
}

// --- Layout Container Methods ---

// GetPadding はレイアウト計算のためにパディング情報を返します。
// layout.Containerインターフェースを実装します。
func (c *Container) GetPadding() layout.Insets {
	style := c.GetStyle()
	if style.Padding != nil {
		return layout.Insets{
			// [修正] style.Insetsの各フィールドはintなので、キャストは不要です。
			Top:    style.Padding.Top,
			Right:  style.Padding.Right,
			Bottom: style.Padding.Bottom,
			Left:   style.Padding.Left,
		}
	}
	// パディングが設定されていない場合は、ゼロ値のInsetsを返します。
	return layout.Insets{}
}