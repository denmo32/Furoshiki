package container

import (
	"errors"
	"furoshiki/component"
	"furoshiki/layout"
)

// ContainerBuilder is a fluent builder for Container widgets.
type ContainerBuilder struct {
	component.Builder[*ContainerBuilder, *Container]
}

// NewContainerBuilder creates a new builder for a container.
func NewContainerBuilder() *ContainerBuilder {
	c := &Container{
		children: make([]component.Widget, 0),
	}
	// 【改善】コンストラクタからself引数を削除し、Initメソッドでself参照を設定する方式に統一します。
	// これにより、コンパイルエラーが解消され、初期化ロジックが直感的になります。
	c.LayoutableWidget = component.NewLayoutableWidget()
	c.Init(c) // ContainerはLayoutableWidgetを埋め込んでいるため、Initメソッドを直接呼び出せます。
	c.layout = &layout.FlexLayout{} // Default layout

	b := &ContainerBuilder{}
	b.Init(b, c) // Init the embedded component.Builder
	return b
}

// GetLayout returns the container's layout manager.
func (b *ContainerBuilder) GetLayout() layout.Layout {
	return b.Widget.GetLayout()
}

// Layout sets the container's layout manager.
func (b *ContainerBuilder) Layout(layout layout.Layout) *ContainerBuilder {
	if layout == nil {
		b.AddError(errors.New("layout cannot be nil"))
		return b
	}
	b.Widget.SetLayout(layout)
	return b
}

// AddChild adds a child widget to the container.
func (b *ContainerBuilder) AddChild(child component.Widget) *ContainerBuilder {
	if child == nil {
		b.AddError(errors.New("child cannot be nil"))
		return b
	}
	b.Widget.AddChild(child)
	return b
}

// AddChildren adds multiple child widgets to the container.
func (b *ContainerBuilder) AddChildren(children ...component.Widget) *ContainerBuilder {
	for _, child := range children {
		if child == nil {
			b.AddError(errors.New("child cannot be nil"))
			continue
		}
		b.Widget.AddChild(child)
	}
	return b
}

// RelayoutBoundary sets the container as a relayout boundary.
func (b *ContainerBuilder) RelayoutBoundary(isBoundary bool) *ContainerBuilder {
	b.Widget.SetRelayoutBoundary(isBoundary)
	return b
}

// [新規追加]
// ClipChildren は、コンテナがその境界外に子要素を描画しないように設定します（クリッピング）。
// trueに設定すると、子要素はコンテナのパディング領域の内側にのみ描画されます。
// これは、スクロール可能な領域などを作成する際の基礎となります。
func (b *ContainerBuilder) ClipChildren(clips bool) *ContainerBuilder {
	b.Widget.SetClipsChildren(clips)
	return b
}

// Build finalizes the container construction.
func (b *ContainerBuilder) Build() (*Container, error) {
	// The embedded builder's Build method returns (W, error), which is (*Container, error)
	return b.Builder.Build()
}