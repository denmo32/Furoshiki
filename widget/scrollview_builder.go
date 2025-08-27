package widget

import (
	"fmt"
	"furoshiki/component"
)

// ScrollViewBuilder は、ScrollViewを宣言的に構築するためのビルダーです。
type ScrollViewBuilder struct {
	component.Builder[*ScrollViewBuilder, *ScrollView]
}

// NewScrollViewBuilder は新しいScrollViewBuilderを生成します。
func NewScrollViewBuilder() *ScrollViewBuilder {
	sv, err := NewScrollView()
	b := &ScrollViewBuilder{}
	b.Init(b, sv)
	// NOTE: コンストラクタで発生した初期化エラーをビルダーに追加します。
	b.AddError(err)
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

// ScrollSensitivity は、マウスホイールのスクロール感度を設定します。
func (b *ScrollViewBuilder) ScrollSensitivity(sensitivity float64) *ScrollViewBuilder {
	if sensitivity < 0 {
		b.AddError(fmt.Errorf("scroll sensitivity cannot be negative, got %f", sensitivity))
	} else {
		b.Widget.ScrollSensitivity = sensitivity
	}
	return b
}

// Build は、最終的なScrollViewを構築して返します。
func (b *ScrollViewBuilder) Build() (*ScrollView, error) {
	return b.Builder.Build()
}
