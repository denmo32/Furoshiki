package widget

import (
	"furoshiki/component"
)

// ScrollViewBuilder は、ScrollViewを宣言的に構築するためのビルダーです。
type ScrollViewBuilder struct {
	component.Builder[*ScrollViewBuilder, *ScrollView]
}

// NewScrollViewBuilder は新しいScrollViewBuilderを生成します。
func NewScrollViewBuilder() *ScrollViewBuilder {
	sv := NewScrollView()
	b := &ScrollViewBuilder{}
	b.Init(b, sv)
	return b
}

// Content はScrollViewのスクロール可能なコンテンツとしてウィジェットを設定します。
func (b *ScrollViewBuilder) Content(content component.Widget) *ScrollViewBuilder {
	if content == nil {
		b.AddError(component.ErrNilChild)
		return b
	}
	b.Widget.SetContent(content)
	return b
}

// Build は、設定に基づいて最終的なScrollViewを構築して返します。
func (b *ScrollViewBuilder) Build() (*ScrollView, error) {
	return b.Builder.Build()
}