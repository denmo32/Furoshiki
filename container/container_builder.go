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
	c.Init(c)                       // ContainerはLayoutableWidgetを埋め込んでいるため、Initメソッドを直接呼び出せます。
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
		b.AddError(component.ErrNilChild)
		return b
	}
	b.Widget.AddChild(child)
	return b
}

// AddChildren adds multiple child widgets to the container.
func (b *ContainerBuilder) AddChildren(children ...component.Widget) *ContainerBuilder {
	for _, child := range children {
		if child == nil {
			b.AddError(component.ErrNilChild)
			continue
		}
		b.Widget.AddChild(child)
	}
	return b
}

// 【改善】SetLayoutBoundary はコンテナをレイアウト境界として設定します。
// メソッド名を SetRelayoutBoundary から SetLayoutBoundary に変更して直感性を向上させました。
func (b *ContainerBuilder) SetLayoutBoundary(isBoundary bool) *ContainerBuilder {
	b.Widget.SetLayoutBoundary(isBoundary)
	return b
}

// 【改善】SetRelayoutBoundary はコンテナをレイアウト境界として設定します。
// 後方互換性のために残します。内部的には SetLayoutBoundary を呼び出します。
func (b *ContainerBuilder) SetRelayoutBoundary(isBoundary bool) *ContainerBuilder {
	b.Widget.SetRelayoutBoundary(isBoundary)
	return b
}

// 【改善】SetClipsChildren はコンテナのクリッピング動作を設定します。
func (b *ContainerBuilder) SetClipsChildren(clips bool) *ContainerBuilder {
	b.Widget.SetClipsChildren(clips)
	return b
}

// Build はコンテナの構築を完了します。
func (b *ContainerBuilder) Build() (*Container, error) {
	widget, err := b.Builder.Build()
	if err != nil {
		return nil, err
	}
	return widget, nil
}
