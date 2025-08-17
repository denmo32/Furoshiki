package container

import (
	"errors"
	"fmt"
	"log"

	"furoshiki/component"
	"furoshiki/layout"
	"furoshiki/style"
)

// Containerは子コンポーネントを保持し、レイアウトを管理します。
type Container struct {
	*component.LayoutableComponent
	layout layout.Layout
	// clipChildren bool // 子要素をクリッピングするかどうかのフラグ (将来的な拡張用)
}

// layout.Containerインターフェースを実装していることを確認します。
var _ layout.Container = (*Container)(nil)

// Updateはコンテナと子要素の状態を更新します。
func (c *Container) Update() {
	// isVisible チェックは埋め込み元の Update で行われる
	if c.IsDirty() && c.layout != nil {
		c.layout.Layout(c)
	}
	c.LayoutableComponent.Update()
}

// Container専用のDrawメソッドは定義しません。
// 埋め込まれた component.LayoutableComponent の Draw メソッドが使用されます。

// --- ContainerBuilder ---
type ContainerBuilder struct {
	container *Container
	errors    []error
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		container: &Container{
			// NewLayoutableComponent を使用して、isVisible などのデフォルト値を正しく設定
			LayoutableComponent: component.NewLayoutableComponent(),
			layout:              &layout.AbsoluteLayout{},
		},
	}
}

// RelayoutBoundary は、このコンテナをレイアウトの境界として設定します。
// これにより、このコンテナ内部の変更が親コンテナのレイアウトに影響を与えなくなります。
// 動的なコンテンツを持つコンテナ（例：スクロールリスト）のパフォーマンスを向上させるのに役立ちます。
func (b *ContainerBuilder) RelayoutBoundary(isBoundary bool) *ContainerBuilder {
	// LayoutableComponentの公開されていないメソッドを直接呼び出す代わりに、
	// コンテナが持つLayoutableComponentのポインタを介して設定する必要があります。
	// 現状のLayoutableComponentは値型で埋め込まれているため、component側にも修正が必要です。
	// ここでは、将来的なインターフェースを想定した仮実装とします。
	// → component.go側でSetRelayoutBoundaryを実装済みのため、直接呼び出し可能。
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
	b.container.layout = layout
	return b
}

func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
	existingStyle := b.container.GetStyle()
	b.container.SetStyle(style.Merge(*existingStyle, s))
	return b
}

func (b *ContainerBuilder) Flex(flex int) *ContainerBuilder {
	b.container.SetFlex(flex)
	return b
}

func (b *ContainerBuilder) AddChild(child component.Component) *ContainerBuilder {
	b.container.AddChild(child)
	return b
}

func (b *ContainerBuilder) AddChildren(children ...component.Component) *ContainerBuilder {
	for _, child := range children {
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

	b.container.MarkDirty()
	return b.container, nil
}
