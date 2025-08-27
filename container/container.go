package container

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/layout"
	"log"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
)

// Scroller は、Containerがクリッピング描画時にスクロールオフセットを
// 取得するために使用するインターフェースです。
// ScrollableContainerのような、スクロール機能を持つウィジェットがこのインターフェースを実装します。
type Scroller interface {
	GetScrollOffset() (x, y int)
}

// Containerは子Widgetを保持し、レイアウトを管理するコンポーネントです。
type Container struct {
	*component.LayoutableWidget
	children []component.Widget
	layout   layout.Layout
	warned   bool // サイズ警告を一度だけ出すためのフラグ

	clipsChildren  bool          // 子要素をクリッピングするかどうか
	offscreenImage *ebiten.Image // クリッピング描画用のオフスクリーンバッファ
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)

// NewContainer は、ビルダーを使わずに新しいContainerインスタンスを生成します。
// NOTE: 内部のInit呼び出しが失敗する可能性があるため、コンストラクタはerrorを返すように変更されました。
// NOTE: アプリケーション開発では、型安全で流暢なUI構築のために、
//       このコンストラクタを直接呼び出すのではなく、`ui.VStack`や`ui.HStack`、
//       または`container.NewContainerBuilder()`の使用を強く推奨します。
func NewContainer() (*Container, error) {
	c := &Container{
		children: make([]component.Widget, 0),
	}
	c.LayoutableWidget = component.NewLayoutableWidget()
	// NOTE: Initがエラーを返すようになったため、コンストラクタもエラーを返すように変更。
	// これにより、初期化の失敗を呼び出し元に安全に伝えることができます。
	if err := c.Init(c); err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}
	c.layout = &layout.FlexLayout{} // デフォルトはFlexLayout
	return c, nil
}

// SetClipsChildren はコンテナのクリッピング動作を設定します。
func (c *Container) SetClipsChildren(clips bool) {
	if c.clipsChildren != clips {
		c.clipsChildren = clips
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
// このメソッドはUIツリーのルートから毎フレーム再帰的に呼び出されます。
func (c *Container) Update() {
	if !c.IsVisible() {
		return
	}

	c.checkSizeWarning()

	if c.IsDirty() {
		if c.NeedsRelayout() {
			if c.layout != nil {
				// NOTE: レイアウト計算がエラーを返すように変更されたため、ここでハンドリングします。
				//       以前のpanic/recoverモデルから移行し、より予測可能なエラー処理を実現します。
				if err := c.layout.Layout(c); err != nil {
					// レイアウト計算中にエラーが発生した場合、ログに出力します。
					// これにより、開発者はレイアウトに関する問題を早期に発見できます。
					log.Printf("Error during layout calculation: %v\n%s", err, debug.Stack())
				}
			}
		}
		c.ClearDirty()
	}

	for _, child := range c.children {
		child.Update()
	}
}

// checkSizeWarning はコンテナのサイズに関する警告を出力します。
func (c *Container) checkSizeWarning() {
	if c.warned {
		return
	}
	if c.GetFlex() == 0 {
		width, height := c.GetSize()
		if width == 0 && height == 0 && c.GetParent() == nil {
			fmt.Printf("Warning: Root container created with no flex and zero size. It may not be visible.\n")
			c.warned = true
		}
	}
}

// UPDATE: DrawメソッドのシグネチャをDrawInfoを受け取るように変更
// Drawはコンテナ自身と、そのすべての子を描画します。
func (c *Container) Draw(info component.DrawInfo) {
	if !c.IsVisible() {
		return
	}

	if c.clipsChildren {
		c.drawWithClipping(info)
	} else {
		c.drawWithoutClipping(info)
	}
}

// UPDATE: drawWithoutClippingのシグネチャをDrawInfoを受け取るように変更
// drawWithoutClipping は、クリッピングを行わずにコンテナと子を描画します。
func (c *Container) drawWithoutClipping(info component.DrawInfo) {
	x, y := c.GetPosition()
	width, height := c.GetSize()
	// NOTE: パフォーマンス向上のためReadOnlyStyle()を使用します。
	// UPDATE: 親から渡されたオフセットを描画座標に適用
	finalX := x + info.OffsetX
	finalY := y + info.OffsetY
	component.DrawStyledBackground(info.Screen, finalX, finalY, width, height, c.ReadOnlyStyle())

	for _, child := range c.children {
		// UPDATE: 子の描画にもオフセット情報を伝播
		child.Draw(info)
	}
}

// UPDATE: drawWithClippingのシグネチャをDrawInfoを受け取るように変更し、副作用を完全排除
// drawWithClipping は、オフスクリーンバッファを利用してクリッピングを行いながら描画します。
func (c *Container) drawWithClipping(info component.DrawInfo) {
	containerX, containerY := c.GetPosition()
	containerWidth, containerHeight := c.GetSize()

	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	// オフスクリーン画像の準備
	if c.offscreenImage == nil || c.offscreenImage.Bounds().Dx() != containerWidth || c.offscreenImage.Bounds().Dy() != containerHeight {
		if c.offscreenImage != nil {
			c.offscreenImage.Deallocate()
		}
		c.offscreenImage = ebiten.NewImage(containerWidth, containerHeight)
	}
	c.offscreenImage.Clear()

	// コンテナ自身の背景をオフスクリーン画像に描画(オフセットは(0,0))
	// NOTE: パフォーマンス向上のためReadOnlyStyle()を使用します。
	component.DrawStyledBackground(c.offscreenImage, 0, 0, containerWidth, containerHeight, c.ReadOnlyStyle())

	// コンテナ自身がScrollerインターフェースを実装しているかチェック
	var scrollOffsetX, scrollOffsetY int
	if scroller, ok := any(c).(Scroller); ok {
		scrollOffsetX, scrollOffsetY = scroller.GetScrollOffset()
	}

	// UPDATE: 子ウィジェットの座標を変更する代わりに、描画オフセットを計算して渡します。
	// これにより、描画処理中の状態変更（副作用）がなくなり、コードの堅牢性が向上します。
	// 子ウィジェットは、自身の絶対座標をこのオフセットに基づいてオフスクリーンバッファ上の
	// ローカル座標に変換して描画します。
	childDrawInfo := component.DrawInfo{
		Screen:  c.offscreenImage,
		OffsetX: -(containerX - scrollOffsetX),
		OffsetY: -(containerY - scrollOffsetY),
	}

	// 子要素をオフスクリーン画像に描画
	for _, child := range c.children {
		// オフセットされた座標でオフスクリーン画像に描画
		child.Draw(childDrawInfo)
	}

	// 完成したオフスクリーン画像をスクリーンに描画
	opts := &ebiten.DrawImageOptions{}
	// UPDATE: 親から渡されたオフセットを最終的な描画位置に適用
	finalX := float64(containerX + info.OffsetX)
	finalY := float64(containerY + info.OffsetY)
	opts.GeoM.Translate(finalX, finalY)
	info.Screen.DrawImage(c.offscreenImage, opts)
}

// HitTest は、指定された座標がコンテナまたはその子のいずれかにヒットするかをテストします。
func (c *Container) HitTest(x, y int) component.Widget {
	if !c.IsVisible() {
		return nil
	}
	// 描画順とは逆に、最前面の子からヒットテストします。
	for i := len(c.children) - 1; i >= 0; i-- {
		child := c.children[i]
		if !child.IsVisible() {
			continue
		}
		if target := child.HitTest(x, y); target != nil {
			return target
		}
	}
	// どの子にもヒットしなかった場合、コンテナ自身をテストします。
	// LayoutableWidget.HitTestは、ヒットした場合にコンテナ自身(c)を返します。
	return c.LayoutableWidget.HitTest(x, y)
}

// Cleanup は、コンテナとすべての子ウィジェットのリソースを解放します。
func (c *Container) Cleanup() {
	for _, child := range c.children {
		child.Cleanup()
	}
	c.children = nil

	if c.offscreenImage != nil {
		c.offscreenImage.Deallocate()
		c.offscreenImage = nil
	}

	c.LayoutableWidget.Cleanup()
}

// detachChildは、親子関係のみを解消する内部ヘルパーです。
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
	if oldParent := child.GetParent(); oldParent != nil {
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
	if c.detachChild(child) {
		child.Cleanup()
		c.MarkDirty(true)
	}
}

// GetChildren はコンテナが保持するすべての子ウィジェットのスライスを返します。
func (c *Container) GetChildren() []component.Widget {
	return c.children
}

// GetPadding はレイアウト計算のためにパディング情報を返します。
func (c *Container) GetPadding() layout.Insets {
	// NOTE: パフォーマンス向上のためReadOnlyStyle()を使用します。
	style := c.ReadOnlyStyle()
	if style.Padding != nil {
		return layout.Insets{
			Top:    style.Padding.Top,
			Right:  style.Padding.Right,
			Bottom: style.Padding.Bottom,
			Left:   style.Padding.Left,
		}
	}
	return layout.Insets{}
}