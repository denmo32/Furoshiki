package ui

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/widget"
)

// 【改善】このファイルは、レイアウトごとに特化した型付きビルダーを提供します。
// これにより、例えばGridのコンテキストでFlexbox用のメソッドを呼び出すといった
// 間違いをコンパイル時に防ぎ、APIの型安全性を向上させます。

// =================================================================
//
//	共通ヘルパー関数 (非公開)
//
// =================================================================

// builderConstraint は、uiパッケージ内のすべての型付きビルダーが満たすべき振る舞いを定義します。
type builderConstraint interface {
	component.ErrorAdder
	component.WidgetContainer
}

// buildContainerAndBuilder は、指定されたレイアウトでコンテナとビルダーを初期化する内部ヘルパーです。
// 【修正】型パラメータ`B`は構造体型、`PB`はそのポインタ型を想定しています。
func buildContainerAndBuilder[B any, PB interface {
	*B
	component.BuilderInitializer[*B, *container.Container]
}](l layout.Layout, buildFunc func(PB)) PB {
	c := &container.Container{}
	c.LayoutableWidget = component.NewLayoutableWidget()
	c.Init(c)
	c.SetLayout(l)

	var b PB = new(B)
	b.Init(b, c)

	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// addWidget は、ウィジェットをビルドして親ビルダーに追加する共通ロジックです。
func addWidget[B builderConstraint, W component.Widget, WB interface {
	component.BuilderFinalizer[W]
	component.ErrorAdder
}](parentBuilder B, widgetBuilder WB) {
	widget, err := widgetBuilder.Build()
	parentBuilder.AddError(err)
	parentBuilder.AddChild(widget)
}

// addNestedContainer は、ネストされたコンテナを追加する共通ロジックです。
// 【修正】型パラメータ`C`の制約を`component.Widget`に変更しました。
// `*container.Container`は`component.Widget`を満たすため、これは有効です。
func addNestedContainer[B builderConstraint, C component.Widget, CB interface {
	component.BuilderFinalizer[C]
	component.ErrorAdder
}](parentBuilder B, nestedBuilder CB) {
	containerWidget, err := nestedBuilder.Build()
	parentBuilder.AddError(err)
	parentBuilder.AddChild(containerWidget)
}

// =================================================================
//
//	FlexBuilder (VStack, HStack用)
//
// =================================================================

type FlexBuilder struct {
	component.Builder[*FlexBuilder, *container.Container]
}

func VStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	flexLayout := &layout.FlexLayout{Direction: layout.DirectionColumn}
	// 【修正】型パラメータには構造体型`FlexBuilder`を渡します。
	return buildContainerAndBuilder(flexLayout, buildFunc)
}

func HStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	flexLayout := &layout.FlexLayout{Direction: layout.DirectionRow}
	// 【修正】型パラメータには構造体型`FlexBuilder`を渡します。
	return buildContainerAndBuilder(flexLayout, buildFunc)
}

// --- FlexBuilder Methods ---
func (b *FlexBuilder) AddChild(child component.Widget) {
	if child != nil {
		b.Widget.AddChild(child)
	} else {
		b.AddError(component.ErrNilChild)
	}
}
func (b *FlexBuilder) Label(buildFunc func(*widget.LabelBuilder)) *FlexBuilder {
	builder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *FlexBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *FlexBuilder {
	builder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *FlexBuilder) Spacer() *FlexBuilder {
	addWidget(b, widget.NewSpacerBuilder().Flex(1))
	return b
}
func (b *FlexBuilder) ScrollView(buildFunc func(*widget.ScrollViewBuilder)) *FlexBuilder {
	builder := widget.NewScrollViewBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *FlexBuilder) HStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	addNestedContainer(b, HStack(buildFunc))
	return b
}
func (b *FlexBuilder) VStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	addNestedContainer(b, VStack(buildFunc))
	return b
}
func (b *FlexBuilder) ZStack(buildFunc func(*ZStackBuilder)) *FlexBuilder {
	addNestedContainer(b, ZStack(buildFunc))
	return b
}
func (b *FlexBuilder) Grid(buildFunc func(*GridBuilder)) *FlexBuilder {
	addNestedContainer(b, Grid(buildFunc))
	return b
}
func (b *FlexBuilder) RelayoutBoundary(isBoundary bool) *FlexBuilder {
	b.Widget.SetRelayoutBoundary(isBoundary)
	return b
}
func (b *FlexBuilder) ClipChildren(clips bool) *FlexBuilder {
	b.Widget.SetClipsChildren(clips)
	return b
}
func (b *FlexBuilder) Gap(gap int) *FlexBuilder {
	if flexLayout := b.Widget.GetLayout().(*layout.FlexLayout); flexLayout.Gap != gap {
		flexLayout.Gap = gap
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *FlexBuilder) Justify(alignment layout.Alignment) *FlexBuilder {
	if flexLayout := b.Widget.GetLayout().(*layout.FlexLayout); flexLayout.Justify != alignment {
		flexLayout.Justify = alignment
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *FlexBuilder) AlignItems(alignment layout.Alignment) *FlexBuilder {
	if flexLayout := b.Widget.GetLayout().(*layout.FlexLayout); flexLayout.AlignItems != alignment {
		flexLayout.AlignItems = alignment
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *FlexBuilder) Build() (*container.Container, error) {
	return b.Builder.Build()
}

// =================================================================
//
//	GridBuilder (Grid用)
//
// =================================================================

type GridBuilder struct {
	component.Builder[*GridBuilder, *container.Container]
}

func Grid(buildFunc func(*GridBuilder)) *GridBuilder {
	gridLayout := &layout.GridLayout{Columns: 1}
	// 【修正】型パラメータには構造体型`GridBuilder`を渡します。
	return buildContainerAndBuilder(gridLayout, buildFunc)
}

// --- GridBuilder Methods ---
func (b *GridBuilder) AddChild(child component.Widget) {
	if child != nil {
		b.Widget.AddChild(child)
	} else {
		b.AddError(component.ErrNilChild)
	}
}
func (b *GridBuilder) Label(buildFunc func(*widget.LabelBuilder)) *GridBuilder {
	builder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *GridBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *GridBuilder {
	builder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}

// Note: Spacer() is omitted as it has no effect in GridLayout.
func (b *GridBuilder) ScrollView(buildFunc func(*widget.ScrollViewBuilder)) *GridBuilder {
	builder := widget.NewScrollViewBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *GridBuilder) HStack(buildFunc func(*FlexBuilder)) *GridBuilder {
	addNestedContainer(b, HStack(buildFunc))
	return b
}
func (b *GridBuilder) VStack(buildFunc func(*FlexBuilder)) *GridBuilder {
	addNestedContainer(b, VStack(buildFunc))
	return b
}
func (b *GridBuilder) ZStack(buildFunc func(*ZStackBuilder)) *GridBuilder {
	addNestedContainer(b, ZStack(buildFunc))
	return b
}
func (b *GridBuilder) Grid(buildFunc func(*GridBuilder)) *GridBuilder {
	addNestedContainer(b, Grid(buildFunc))
	return b
}
func (b *GridBuilder) RelayoutBoundary(isBoundary bool) *GridBuilder {
	b.Widget.SetRelayoutBoundary(isBoundary)
	return b
}
func (b *GridBuilder) ClipChildren(clips bool) *GridBuilder {
	b.Widget.SetClipsChildren(clips)
	return b
}
func (b *GridBuilder) Columns(count int) *GridBuilder {
	if gridLayout := b.Widget.GetLayout().(*layout.GridLayout); count > 0 && gridLayout.Columns != count {
		gridLayout.Columns = count
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *GridBuilder) Rows(count int) *GridBuilder {
	if gridLayout := b.Widget.GetLayout().(*layout.GridLayout); gridLayout.Rows != count {
		gridLayout.Rows = count
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *GridBuilder) HorizontalGap(gap int) *GridBuilder {
	if gridLayout := b.Widget.GetLayout().(*layout.GridLayout); gridLayout.HorizontalGap != gap {
		gridLayout.HorizontalGap = gap
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *GridBuilder) VerticalGap(gap int) *GridBuilder {
	if gridLayout := b.Widget.GetLayout().(*layout.GridLayout); gridLayout.VerticalGap != gap {
		gridLayout.VerticalGap = gap
		b.Widget.MarkDirty(true)
	}
	return b
}
func (b *GridBuilder) Build() (*container.Container, error) {
	return b.Builder.Build()
}

// =================================================================
//
//	ZStackBuilder (ZStack用)
//
// =================================================================

type ZStackBuilder struct {
	component.Builder[*ZStackBuilder, *container.Container]
}

func ZStack(buildFunc func(*ZStackBuilder)) *ZStackBuilder {
	// 【修正】型パラメータには構造体型`ZStackBuilder`を渡します。
	return buildContainerAndBuilder(&layout.AbsoluteLayout{}, buildFunc)
}

// --- ZStackBuilder Methods ---
func (b *ZStackBuilder) AddChild(child component.Widget) {
	if child != nil {
		b.Widget.AddChild(child)
	} else {
		b.AddError(component.ErrNilChild)
	}
}
func (b *ZStackBuilder) Label(buildFunc func(*widget.LabelBuilder)) *ZStackBuilder {
	builder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *ZStackBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *ZStackBuilder {
	builder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}

// 【新規追加】ZStackBuilderに他のコンテナやウィジェットを追加するメソッドを拡張しました。
// これにより、FlexBuilderなど他のビルダーとのAPIの一貫性が向上し、
// ZStack内でより複雑なレイアウトを宣言的に構築できるようになります。
func (b *ZStackBuilder) ScrollView(buildFunc func(*widget.ScrollViewBuilder)) *ZStackBuilder {
	builder := widget.NewScrollViewBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	addWidget(b, builder)
	return b
}
func (b *ZStackBuilder) HStack(buildFunc func(*FlexBuilder)) *ZStackBuilder {
	addNestedContainer(b, HStack(buildFunc))
	return b
}
func (b *ZStackBuilder) VStack(buildFunc func(*FlexBuilder)) *ZStackBuilder {
	addNestedContainer(b, VStack(buildFunc))
	return b
}
func (b *ZStackBuilder) ZStack(buildFunc func(*ZStackBuilder)) *ZStackBuilder {
	addNestedContainer(b, ZStack(buildFunc))
	return b
}
func (b *ZStackBuilder) Grid(buildFunc func(*GridBuilder)) *ZStackBuilder {
	addNestedContainer(b, Grid(buildFunc))
	return b
}

func (b *ZStackBuilder) Build() (*container.Container, error) {
	return b.Builder.Build()
}
