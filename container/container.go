package container

import (
	"fmt"

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
}

// コンパイル時にインターフェースの実装を検証
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)

// Updateはコンテナと子要素の状態を更新します。
// component.LayoutableWidgetのUpdateをオーバーライドして、レイアウト計算と子の更新を追加します。
// このメソッドはUIツリーのルートから毎フレーム再帰的に呼び出されます。
func (c *Container) Update() {
	if !c.IsVisible() {
		return
	}

	// このコンテナ、またはその子孫のいずれかで再レイアウトが必要な場合 (relayoutDirty=true)、
	// IsDirty() は true を返します。その場合、レイアウトを再計算します。
	if c.IsDirty() {
		if c.layout != nil {
			// レイアウト計算中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
			defer func() {
				if r := recover(); r != nil {
					// 将来的には、より高度なロギングやエラー報告メカニズムに置き換えることができます。
					fmt.Printf("Recovered from panic during layout calculation: %v\n", r)
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

// Drawはコンテナ自身と、そのすべての子を描画します。
// component.LayoutableWidgetのDrawをオーバーライドします。
func (c *Container) Draw(screen *ebiten.Image) {
	// コンテナが非表示の場合、自身も子も描画しない。
	if !c.IsVisible() {
		return
	}
	// まずコンテナ自身の背景などを描画
	c.LayoutableWidget.Draw(screen)
	// 次にすべての子ウィジェットを描画
	for _, child := range c.children {
		child.Draw(screen)
	}
}

// HitTest は、指定された座標がコンテナまたはその子のいずれかにヒットするかをテストします。
// component.LayoutableWidgetのHitTestをオーバーライドします。
func (c *Container) HitTest(x, y int) component.Widget {
	// コンテナが非表示の場合はヒットしない
	if !c.IsVisible() {
		return nil
	}
	// 描画順とは逆に、最前面の子からヒットテストを行う
	for i := len(c.children) - 1; i >= 0; i-- {
		child := c.children[i]
		// 非表示の子はヒットテストの対象外
		if !child.IsVisible() {
			continue
		}
		if target := child.HitTest(x, y); target != nil {
			return target // ヒットした子を返す
		}
	}
	// どの子にもヒットしなかった場合、コンテナ自身がヒットするかチェック
	if target := c.LayoutableWidget.HitTest(x, y); target != nil {
		return c // コンテナ自身を返す
	}
	return nil
}

// [追加] Cleanup は、コンテナとすべての子ウィジェットのリソースを解放します。
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

func (c *Container) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	// 既に親が存在する場合は、その親から子を削除
	if oldParent := child.GetParent(); oldParent != nil {
		oldParent.RemoveChild(child)
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.MarkDirty(true)
}

func (c *Container) RemoveChild(child component.Widget) {
	if child == nil {
		return
	}
	for i, currentChild := range c.children {
		if currentChild == child {
			// スライスから子を削除
			c.children = append(c.children[:i], c.children[i+1:]...)
			// 親への参照をクリア
			child.SetParent(nil)
			// 子のクリーンアップ処理を呼び出します。
			// component.WidgetインターフェースはCleanup()メソッドを保証しているため、型アサーションは不要です。
			child.Cleanup()
			// コンテナの再レイアウトを要求
			c.MarkDirty(true)
			return
		}
	}
}

func (c *Container) GetChildren() []component.Widget {
	return c.children
}

// --- Layout Container Methods ---

// GetPadding はレイアウト計算のためにパディング情報を返します。
// layout.Containerインターフェースを実装します。
func (c *Container) GetPadding() layout.Insets {
	// [改善] GetStyle()が値型を返すようになったため、ポインタアクセス(*)が不要になります。
	style := c.GetStyle()
	return layout.Insets{
		Top:    style.Padding.Top,
		Right:  style.Padding.Right,
		Bottom: style.Padding.Bottom,
		Left:   style.Padding.Left,
	}
}