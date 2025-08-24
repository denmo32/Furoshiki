package container

import (
    "furoshiki/component"
)

// ScrollableContainerBuilder はScrollableContainerを構築するためのビルダーです
type ScrollableContainerBuilder struct {
    component.Builder[*ScrollableContainerBuilder, *ScrollableContainer]
}

// NewScrollableContainerBuilder は新しいScrollableContainerBuilderを作成します
func NewScrollableContainerBuilder() *ScrollableContainerBuilder {
    sc := NewScrollableContainer()
    
    b := &ScrollableContainerBuilder{}
    b.Init(b, sc)
    return b
}

// AddChild は子ウィジェットを追加します
func (b *ScrollableContainerBuilder) AddChild(child component.Widget) *ScrollableContainerBuilder {
    if child == nil {
        b.AddError(component.ErrNilChild)
        return b
    }
    b.Widget.AddChild(child)
    return b
}

// Build はScrollableContainerの構築を完了します
func (b *ScrollableContainerBuilder) Build() (*ScrollableContainer, error) {
    container, err := b.Builder.Build()
    if err != nil {
        return nil, err
    }
    
    // ビルド時にコンテンツサイズを計算
    container.updateContentHeight()
    
    return container, nil
}