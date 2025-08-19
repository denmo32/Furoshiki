package container

import (
	"errors"
	"fmt"

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
// [改善] デフォルトのレイアウトを、より汎用性が高く一般的に使用される FlexLayout に変更します。
// [修正] Containerの初期化を、self参照を渡す新しい方式に変更します。
func NewContainerBuilder() *ContainerBuilder {
	// まずコンテナのインスタンスを生成
	c := &Container{
		children: make([]component.Widget, 0),
	}
	// 次に、コンテナ自身をselfとして渡してLayoutableWidgetを初期化
	c.LayoutableWidget = component.NewLayoutableWidget(c)
	// デフォルトのレイアウトを設定
	c.layout = &layout.FlexLayout{}

	return &ContainerBuilder{
		container: c,
	}
}

// [追加] GetLayout は、ビルド中のコンテナが現在使用しているレイアウトを返します。
// これにより、uiパッケージのヘルパーなどが、コンテナをビルドせずにレイアウトプロパティを変更できます。
func (b *ContainerBuilder) GetLayout() layout.Layout {
	return b.container.GetLayout()
}

// [追加] AddError は、ビルドプロセス中に発生したエラーをビルダーに記録します。
// uiパッケージのヘルパー関数が、子のビルドエラーを親のビルダーに伝播させるために使用します。
func (b *ContainerBuilder) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// RelayoutBoundary は、このコンテナをレイアウトの境界として設定します。
// これにより、このコンテナ内部の変更が親コンテナのレイアウトに影響を与えなくなります。
// 動的なコンテンツを持つコンテナ（例：スクロールリスト）のパフォーマンスを向上させるのに役立ちます。
func (b *ContainerBuilder) RelayoutBoundary(isBoundary bool) *ContainerBuilder {
	b.container.SetRelayoutBoundary(isBoundary)
	return b
}

// Position はコンテナのスクリーン上の絶対位置を設定します。
// 注: このコンテナがFlexLayoutを持つ親の中にある場合、この位置はレイアウト計算によって上書きされます。
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
		b.errors = append(b.errors, errors.New("layout cannot be nil"))
		return b
	}
	b.container.SetLayout(layout)
	return b
}

// Style はコンテナのスタイルを設定します。
func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
	existingStyle := b.container.GetStyle()
	b.container.SetStyle(style.Merge(existingStyle, s))
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
		b.errors = append(b.errors, errors.New("child cannot be nil"))
		return b
	}
	b.container.AddChild(child)
	return b
}

// AddChildren はコンテナに複数の子ウィジェットを一度に追加します。
func (b *ContainerBuilder) AddChildren(children ...component.Widget) *ContainerBuilder {
	for _, child := range children {
		if child == nil {
			b.errors = append(b.errors, errors.New("child cannot be nil"))
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

	// [削除] ビルド時の警告チェックを削除 - Updateメソッド内でチェックするように変更済

	b.container.MarkDirty(true)
	return b.container, nil
}