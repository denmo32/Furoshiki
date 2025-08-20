package container

import (
	"fmt"
	"log"           // ログ出力のために追加
	"runtime/debug" // スタックトレース取得のために追加

	"furoshiki/component"
	"furoshiki/layout"

	"github.com/hajimehoshi/ebiten/v2"
)

// Containerは子Widgetを保持し、レイアウトを管理するコンポーネントです。
// component.Containerインターフェースを実装します。
type Container struct {
	*component.LayoutableWidget
	children []component.Widget
	layout   layout.Layout
	warned   bool // サイズ警告を一度だけ出すためのフラグ
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)

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

	// このコンテナ、またはその子孫のいずれかで再レイアウトが必要な場合、IsDirty()はtrueを返します。
	// その場合、レイアウトを再計算します。
	if c.IsDirty() {
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
		// レイアウトが完了したので、ダーティフラグをクリアします。
		// これにより、次のフレームで不要な再計算が行われるのを防ぎます。
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
func (c *Container) Draw(screen *ebiten.Image) {
	// コンテナが非表示の場合、自身も子も描画しません。
	if !c.IsVisible() {
		return
	}
	// まずコンテナ自身の背景などを描画ヘルパーで描画します。
	x, y := c.GetPosition()
	width, height := c.GetSize()
	component.DrawStyledBackground(screen, x, y, width, height, c.GetStyle())

	// 次にすべての子ウィジェットを描画します。
	for _, child := range c.children {
		child.Draw(screen)
	}
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
	// どの子にもヒットしなかった場合、コンテナ自身の領域（子の領域を除く部分）がヒットしたかテストします。
	if c.LayoutableWidget.HitTest(x, y) != nil {
		return c // コンテナ自身を返します。
	}
	return nil
}

// Cleanup は、コンテナとすべての子ウィジェットのリソースを解放します。
// UIツリーからコンテナが削除されるときや、アプリケーション終了時に呼び出されるべきです。
func (c *Container) Cleanup() {
	// まず、すべての子ウィジェットのクリーンアップを再帰的に呼び出します。
	for _, child := range c.children {
		child.Cleanup()
	}
	// 子のリストをクリアします。
	c.children = nil

	// 最後に、埋め込まれたLayoutableWidget自身のクリーンアップ処理（イベントハンドラのクリアなど）を呼び出します。
	c.LayoutableWidget.Cleanup()
}

// --- Container Methods (from component.Container interface) ---

// AddChild はコンテナに子ウィジェットを追加します。
func (c *Container) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	// 既に親が存在する場合は、その親から子を削除し、新しい親子関係を構築します。
	if oldParent := child.GetParent(); oldParent != nil {
		oldParent.RemoveChild(child)
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.MarkDirty(true)
}

// RemoveChild はコンテナから子ウィジェットを削除します。
func (c *Container) RemoveChild(child component.Widget) {
	if child == nil {
		return
	}
	for i, currentChild := range c.children {
		if currentChild == child {
			// スライスから子を削除します。
			c.children = append(c.children[:i], c.children[i+1:]...)
			// 親への参照をクリアします。
			child.SetParent(nil)
			// 子のリソースを解放します。
			child.Cleanup()
			// コンテナの再レイアウトを要求します。
			c.MarkDirty(true)
			return
		}
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
			Top:    style.Padding.Top,
			Right:  style.Padding.Right,
			Bottom: style.Padding.Bottom,
			Left:   style.Padding.Left,
		}
	}
	// パディングが設定されていない場合は、ゼロ値のInsetsを返します。
	return layout.Insets{}
}