package container

import (
	"errors"
	"fmt"
	"log"

	"furoshiki/component"
	"furoshiki/layout"
	"furoshiki/style"
)

// ContainerBuilder は、Containerを安全かつ流れるように構築するためのビルダーです。
type ContainerBuilder struct {
	container *Container
	errors    []error
}

// NewContainerBuilder は、デフォルト値で初期化されたContainerBuilderを返します。
// デフォルトのレイアウトは AbsoluteLayout です。
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

// Position はコンテナの初期位置を設定します。
func (b *ContainerBuilder) Position(x, y int) *ContainerBuilder {
	b.container.SetPosition(x, y)
	return b
}

// Size はコンテナのサイズを設定します。
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

// Layout はコンテナが使用するレイアウトマネージャーを設定します。
func (b *ContainerBuilder) Layout(layout layout.Layout) *ContainerBuilder {
	if layout == nil {
		b.errors = append(b.errors, fmt.Errorf("layout cannot be nil"))
		return b
	}
	b.container.layout = layout
	return b
}

// Style はコンテナのスタイルを設定します。
func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
	existingStyle := b.container.GetStyle()
	b.container.SetStyle(style.Merge(*existingStyle, s))
	return b
}

// Flex は、親がFlexLayoutの場合にコンテナがどのように伸縮するかを設定します。
func (b *ContainerBuilder) Flex(flex int) *ContainerBuilder {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("flex must be non-negative, got %d", flex))
		return b
	}
	b.container.SetFlex(flex)
	return b
}

// AddChild はコンテナに子ウィジェットを追加します。
func (b *ContainerBuilder) AddChild(child component.Widget) *ContainerBuilder {
	if child == nil {
		b.errors = append(b.errors, fmt.Errorf("child cannot be nil"))
		return b
	}
	b.container.AddChild(child)
	return b
}

// AddChildren はコンテナに複数の子ウィジェットを一度に追加します。
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

// Build は、設定に基づいて最終的なContainerを構築して返します。
// 構築中にエラーが発生した場合は、エラーを返します。
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
