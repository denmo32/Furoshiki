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
func NewContainer() *Container {
	c := &Container{
		children: make([]component.Widget, 0),
	}
	c.LayoutableWidget = component.NewLayoutableWidget()
	c.Init(c)
	c.layout = &layout.FlexLayout{} // デフォルトはFlexLayout
	return c
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
				// レイアウト計算は複雑なため、パニックから保護します。
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic during layout calculation: %v\n%s", r, debug.Stack())
					}
				}()
				c.layout.Layout(c)
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

// Drawはコンテナ自身と、そのすべての子を描画します。
func (c *Container) Draw(screen *ebiten.Image) {
	if !c.IsVisible() {
		return
	}

	if c.clipsChildren {
		c.drawWithClipping(screen)
	} else {
		c.drawWithoutClipping(screen)
	}
}

// drawWithoutClipping は、クリッピングを行わずにコンテナと子を描画します。
func (c *Container) drawWithoutClipping(screen *ebiten.Image) {
	x, y := c.GetPosition()
	width, height := c.GetSize()
	component.DrawStyledBackground(screen, x, y, width, height, c.GetStyle())

	for _, child := range c.children {
		child.Draw(screen)
	}
}

// drawWithClipping は、オフスクリーンバッファを利用してクリッピングを行いながら描画します。
func (c *Container) drawWithClipping(screen *ebiten.Image) {
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

	// コンテナ自身の背景をオフスクリーン画像に描画
	component.DrawStyledBackground(c.offscreenImage, 0, 0, containerWidth, containerHeight, c.GetStyle())

	// コンテナ自身がScrollerインターフェースを実装しているかチェック
	var scrollOffsetX, scrollOffsetY int
	if scroller, ok := any(c).(Scroller); ok {
		scrollOffsetX, scrollOffsetY = scroller.GetScrollOffset()
	}

	// 子ウィジェットとその子孫の座標を再帰的に変更するためのヘルパー関数
	var offsetWidgets func(w component.Widget, dx, dy int)
	offsetWidgets = func(w component.Widget, dx, dy int) {
		x, y := w.GetPosition()
		// NOTE: Drawループ中にSetPositionを呼ぶとダーティフラグが立ち、意図しない再描画や再計算が起きる可能性がある。
		// この問題を根本的に解決するには、ライブラリの描画アーキテクチャの変更が必要になる。
		w.SetPosition(x+dx, y+dy)
		if cont, ok := w.(component.Container); ok {
			for _, grandChild := range cont.GetChildren() {
				offsetWidgets(grandChild, dx, dy)
			}
		}
	}

	// 子要素をオフスクリーン画像に描画
	for _, child := range c.children {
		// オフセットを計算
		// 子ウィジェットを、コンテナのローカル座標系(0,0)を基準とした位置に描画するためのオフセット
		offsetX := -(containerX - scrollOffsetX)
		offsetY := -(containerY - scrollOffsetY)

		// 座標を一時的に変更
		offsetWidgets(child, offsetX, offsetY)

		// オフセットされた座標でオフスクリーン画像に描画
		child.Draw(c.offscreenImage)

		// 座標を元に戻す
		offsetWidgets(child, -offsetX, -offsetY)
	}

	// 完成したオフスクリーン画像をスクリーンに描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(containerX), float64(containerY))
	screen.DrawImage(c.offscreenImage, opts)
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
	style := c.GetStyle()
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
