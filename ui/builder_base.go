package ui

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/widget"
)

// BaseContainerBuilder は、コンテナ系のビルダー（FlexBuilder, GridBuilderなど）の
// 共通機能を提供するジェネリックな基底ビルダーです。
// これにより、ウィジェット追加メソッド（Label, Buttonなど）のコード重複をなくし、
// メンテナンス性を向上させます。
// Tは、このビルダーを埋め込む具象ビルダーの型（例: *FlexBuilder）です。
type BaseContainerBuilder[T any] struct {
	component.Builder[T, *container.Container]
}

// Init は基底コンテナビルダーを初期化します。
func (b *BaseContainerBuilder[T]) Init(self T, c *container.Container) {
	b.Builder.Init(self, c)
}

// AddChild は子ウィジェットを追加します。
// このメソッドは component.WidgetContainer インターフェースを満たすために必要です。
func (b *BaseContainerBuilder[T]) AddChild(child component.Widget) {
	if child != nil {
		b.Widget.AddChild(child)
	} else {
		b.AddError(component.ErrNilChild)
	}
}

// --- 共通ウィジェット追加メソッド ---

// Label は、コンテナにLabelウィジェットを追加します。
func (b *BaseContainerBuilder[T]) Label(buildFunc func(*widget.LabelBuilder)) T {
	builder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b.Self
}

// Button は、コンテナにButtonウィジェットを追加します。
func (b *BaseContainerBuilder[T]) Button(buildFunc func(*widget.ButtonBuilder)) T {
	builder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b.Self
}

// Spacer は、コンテナにSpacerウィジェットを追加します。
// 主にFlexLayout内で使用され、利用可能なスペースを埋めるために伸縮します。
func (b *BaseContainerBuilder[T]) Spacer() T {
	// Spacerは通常Flex(1)で使われることが多いため、デフォルトで設定します。
	addWidget(b, widget.NewSpacerBuilder().Flex(1))
	return b.Self
}

// ScrollView は、コンテナにScrollViewウィジェットを追加します。
func (b *BaseContainerBuilder[T]) ScrollView(buildFunc func(*widget.ScrollViewBuilder)) T {
	builder := widget.NewScrollViewBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b.Self
}

// --- ネストされたコンテナ追加メソッド ---

// HStack は、コンテナに水平方向のFlexコンテナをネストして追加します。
func (b *BaseContainerBuilder[T]) HStack(buildFunc func(*FlexBuilder)) T {
	addNestedContainer(b, HStack(buildFunc))
	return b.Self
}

// VStack は、コンテナに垂直方向のFlexコンテナをネストして追加します。
func (b *BaseContainerBuilder[T]) VStack(buildFunc func(*FlexBuilder)) T {
	addNestedContainer(b, VStack(buildFunc))
	return b.Self
}

// ZStack は、コンテナにZ軸方向（重ね合わせ）のコンテナをネストして追加します。
func (b *BaseContainerBuilder[T]) ZStack(buildFunc func(*ZStackBuilder)) T {
	addNestedContainer(b, ZStack(buildFunc))
	return b.Self
}

// Grid は、コンテナにグリッドレイアウトコンテナをネストして追加します。
func (b *BaseContainerBuilder[T]) Grid(buildFunc func(*GridBuilder)) T {
	addNestedContainer(b, Grid(buildFunc))
	return b.Self
}

// AdvancedGrid は、コンテナに高度なグリッドレイアウトコンテナをネストして追加します。
func (b *BaseContainerBuilder[T]) AdvancedGrid(buildFunc func(*AdvancedGridBuilder)) T {
	addNestedContainer(b, AdvancedGrid(buildFunc))
	return b.Self
}

// --- 共通コンテナ設定メソッド ---

// SetLayoutBoundary はコンテナをレイアウト境界として設定します。
// このコンテナ内部でのレイアウト変更が、親コンテナの再レイアウトを引き起こさなくなります。
func (b *BaseContainerBuilder[T]) SetLayoutBoundary(isBoundary bool) T {
	b.Widget.SetLayoutBoundary(isBoundary)
	return b.Self
}

// ClipChildren は、コンテナの境界外に子要素がはみ出して描画されるのを防ぎます（クリッピング）。
func (b *BaseContainerBuilder[T]) ClipChildren(clips bool) T {
	b.Widget.SetClipsChildren(clips)
	return b.Self
}

// --- ビルドヘルパー (非公開) ---

// builderConstraint は、BaseContainerBuilderが内部で使用する制約です。
type builderConstraint interface {
	component.ErrorAdder
	component.WidgetContainer
}

// addWidget は、ウィジェットビルダーからウィジェットをビルドし、親コンテナビルダーに追加します。
func addWidget[B builderConstraint, W component.Widget, WB interface {
	component.BuilderFinalizer[W]
	component.ErrorAdder
}](parentBuilder B, widgetBuilder WB) {
	widget, err := widgetBuilder.Build()
	if err != nil {
		parentBuilder.AddError(err)
	}
	// エラーがあってもウィジェット自体は追加を試みる
	parentBuilder.AddChild(widget)
}

// addNestedContainer は、ネストされたコンテナビルダーをビルドし、親コンテナビルダーに追加します。
func addNestedContainer[B builderConstraint, C component.Widget, CB interface {
	component.BuilderFinalizer[C]
	component.ErrorAdder
}](parentBuilder B, nestedBuilder CB) {
	containerWidget, err := nestedBuilder.Build()
	if err != nil {
		parentBuilder.AddError(err)
	}
	parentBuilder.AddChild(containerWidget)
}