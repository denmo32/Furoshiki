package container

import (
	"fmt"
	"log"           // [追加] ログ出力のために追加
	"runtime/debug" // [追加] スタックトレース取得のために追加

	"furoshiki/component"
	"furoshiki/layout"

	"github.com/hajimehoshi/ebiten/v2"
)

// Containerは子Widgetを保持し、レイアウトを管理するコンポーネントです。
// component.Containerインターフェースを実装します。
type Container struct {
	*component.LayoutableWidget
	children []component.Widget
	layout   layout.Layout // 非公開フィールド
	warned   bool          // サイズ警告を出したかどうかのフラグ
}

// コンパイル時にインターフェースの実装を検証
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)

// GetLayout はコンテナが使用しているレイアウトを返します
func (c *Container) GetLayout() layout.Layout {
	return c.layout
}

// SetLayout はコンテナが使用するレイアウトを設定します
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

	// サイズ警告のチェック - Updateメソッド内でのみチェックするように変更
	c.checkSizeWarning()

	// このコンテナ、またはその子孫のいずれかで再レイアウトが必要な場合 (relayoutDirty=true)、
	// IsDirty() は true を返します。その場合、レイアウトを再計算します。
	if c.IsDirty() {
		if c.layout != nil {
			// レイアウト計算中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
			defer func() {
				if r := recover(); r != nil {
					// [改善] パニック発生時に、より詳細なデバッグ情報（スタックトレース）をログに出力します。
					// 将来的には、より高度なロギングやエラー報告メカニズムに置き換えることができます。
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

// checkSizeWarning はコンテナのサイズに関する警告を出力します
func (c *Container) checkSizeWarning() {
	// すでに警告を出している場合は何もしない
	if c.warned {
		return
	}

	// Flexが設定されておらず、かつサイズが両方0のコンテナは描画されない可能性が高いため警告する。
	// 片方の軸のサイズが0なのは、親のFlexLayout(AlignStretch)に依存する一般的なパターンのため、警告対象外とする。
	if c.GetFlex() == 0 {
		width, height := c.GetSize()
		if width == 0 && height == 0 {
			parent := c.GetParent()
			// 親がいない（ルートコンテナ）の場合のみ警告
			if parent == nil {
				fmt.Printf("Warning: Root container created with no flex and zero size (width=0, height=0). It may not be visible.\n")
				c.warned = true // 警告を出したことを記録
			}
		}
	}
}

// Drawはコンテナ自身と、そのすべての子を描画します。
// component.LayoutableWidgetのDrawをオーバーライドします。
// [修正] 埋め込み先のDrawメソッド呼び出しをやめ、コンテナ自身の背景描画と子の描画呼び出しを明確に分離します。
func (c *Container) Draw(screen *ebiten.Image) {
	// コンテナが非表示の場合、自身も子も描画しない。
	if !c.IsVisible() {
		return
	}
	// まずコンテナ自身の背景などを描画ヘルパーで描画
	x, y := c.GetPosition()
	width, height := c.GetSize()
	component.DrawStyledBackground(screen, x, y, width, height, c.GetStyle())

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
	// [修正] base_widget.HitTestがselfを返すようになったため、戻り値がnilでなければコンテナ自身(c)がヒットしたことになる。
	if c.LayoutableWidget.HitTest(x, y) != nil {
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
// [修正] スタイルのPaddingがポインタ型になったため、nilの場合はゼロ値のInsetsを返します。
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
	// パディングが設定されていない場合は、ゼロ値のInsetsを返す
	return layout.Insets{}
}