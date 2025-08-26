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
	c := NewContainer() // コンストラクタを呼び出す
	b := &ContainerBuilder{}
	b.Init(b, c)
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

// SetLayoutBoundary はコンテナをレイアウト境界として設定します。
func (b *ContainerBuilder) SetLayoutBoundary(isBoundary bool) *ContainerBuilder {
	b.Widget.SetLayoutBoundary(isBoundary)
	return b
}

// SetClipsChildren はコンテナのクリッピング動作を設定します。
func (b *ContainerBuilder) SetClipsChildren(clips bool) *ContainerBuilder {
	b.Widget.SetClipsChildren(clips)
	return b
}

// Build はコンテナの構築を完了します。
func (b *ContainerBuilder) Build() (*Container, error) {
	return b.Builder.Build()
}
