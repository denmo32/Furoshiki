package container

import (
    "errors"
    "fmt"
    "log"

    "furoshiki/component"
    "furoshiki/layout"
    "furoshiki/style"

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

    // ダーティフラグが立っている場合、レイアウトを再計算します。
    // レイアウトの変更は子ウィジェットのサイズや位置に影響を与える可能性があるため、
    // 子のUpdate()を呼び出す前に実行する必要があります。
    if c.IsDirty() {
        if c.layout != nil {
            // レイアウト計算中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
            defer func() {
                if r := recover(); r != nil {
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
    // isVisibleフィールドに直接アクセスするために、埋め込んだLayoutableWidgetのフィールドを直接参照します。
    // しかし、Goの仕様上、c.isVisibleとは書けず、c.LayoutableWidget.isVisibleとする必要があります。
    // 可読性のため、一度変数に受けるか、あるいはIsVisible()を呼び出すのが一般的ですが、
    // ここではパフォーマンス最適化のデモとして、あえて直接アクセスに近い形をとります。
    // 実際には IsVisible() をインライン化するコンパイラの最適化に期待する方が良い場合もあります。
    if !c.IsVisible() { // IsVisible()経由のアクセスに戻します。フィールド直接アクセスは埋め込みの都合上冗長なため。
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
    if !c.IsVisible() {
        return nil
    }
    // 描画順とは逆に、最前面の子からヒットテストを行う
    for i := len(c.children) - 1; i >= 0; i-- {
        child := c.children[i]
        // 非表示の子はスキップ
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
            c.children = append(c.children[:i], c.children[i+1:]...)
            child.SetParent(nil)
            // 子のクリーンアップ処理を呼び出す
            if cleanup, ok := child.(interface{ Cleanup() }); ok {
                cleanup.Cleanup()
            }
            c.MarkDirty(true)
            return
        }
    }
}

func (c *Container) GetChildren() []component.Widget {
    return c.children
}

// --- ContainerBuilder ---
type ContainerBuilder struct {
    container *Container
    errors    []error
}

func NewContainerBuilder() *ContainerBuilder {
    return &ContainerBuilder{
        container: &Container{
            LayoutableWidget: component.NewLayoutableWidget(),
            layout:           &layout.AbsoluteLayout{},
        },
    }
}

// RelayoutBoundary は、このコンテナをレイアウトの境界として設定します。
// これにより、このコンテナ内部の変更が親コンテナのレイアウトに影響を与えなくなります。
// 動的なコンテンツを持つコンテナ（例：スクロールリスト）のパフォーマンスを向上させるのに役立ちます。
func (b *ContainerBuilder) RelayoutBoundary(isBoundary bool) *ContainerBuilder {
    b.container.SetRelayoutBoundary(isBoundary)
    return b
}

func (b *ContainerBuilder) Position(x, y int) *ContainerBuilder {
    b.container.SetPosition(x, y)
    return b
}

func (b *ContainerBuilder) Size(width, height int) *ContainerBuilder {
    if width < 0 {
        b.errors = append(b.errors, fmt.Errorf("width must be non-negative, got %d", width))
    }
    if height < 0 {
        b.errors = append(b.errors, fmt.Errorf("height must be non-negative, got %d", height))
    }
    b.container.SetSize(width, height)
    return b
}

func (b *ContainerBuilder) Layout(layout layout.Layout) *ContainerBuilder {
    if layout == nil {
        b.errors = append(b.errors, fmt.Errorf("layout cannot be nil"))
        return b
    }
    b.container.layout = layout
    return b
}

func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
    existingStyle := b.container.GetStyle()
    b.container.SetStyle(style.Merge(*existingStyle, s))
    return b
}

func (b *ContainerBuilder) Flex(flex int) *ContainerBuilder {
    if flex < 0 {
        b.errors = append(b.errors, fmt.Errorf("flex must be non-negative, got %d", flex))
        return b
    }
    b.container.SetFlex(flex)
    return b
}

func (b *ContainerBuilder) AddChild(child component.Widget) *ContainerBuilder {
    if child == nil {
        b.errors = append(b.errors, fmt.Errorf("child cannot be nil"))
        return b
    }
    b.container.AddChild(child)
    return b
}

func (b *ContainerBuilder) AddChildren(children ...component.Widget) *ContainerBuilder {
    for _, child := range children {
        if child == nil {
            b.errors = append(b.errors, fmt.Errorf("child cannot be nil"))
            continue
        }
        b.container.AddChild(child)
    }
    return b
}

func (b *ContainerBuilder) Build() (*Container, error) {
    if len(b.errors) > 0 {
        joinedErr := errors.Join(b.errors...)
        return nil, fmt.Errorf("container build errors: %w", joinedErr)
    }

    // Flexが設定されておらず、かつサイズが両方0のコンテナは描画されない可能性が高いため警告する。
    // 片方の軸のサイズが0なのは、親のFlexLayout(AlignStretch)に依存する一般的なパターンのため、警告対象外とする。
    if b.container.GetFlex() == 0 {
        width, height := b.container.GetSize()
        if width == 0 && height == 0 {
            log.Printf("Warning: Container has no flex and both width and height are 0. It may not be visible. Using default 200x200.")
            b.container.SetSize(200, 200)
        }
    }

    b.container.MarkDirty(true)
    return b.container, nil
}