package ui

import (
	"fmt"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"
	"image/color"
)

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func VStack(buildFunc func(*FlexContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().Layout(&layout.FlexLayout{
		Direction: layout.DirectionColumn,
	})
	if buildFunc != nil {
		uiBuilder := newFlexContainerBuilder(builder)
		buildFunc(uiBuilder)
	}
	return builder
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func HStack(buildFunc func(*FlexContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().Layout(&layout.FlexLayout{
		Direction: layout.DirectionRow,
	})
	if buildFunc != nil {
		uiBuilder := newFlexContainerBuilder(builder)
		buildFunc(uiBuilder)
	}
	return builder
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成するビルダーです。
func ZStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().Layout(&layout.AbsoluteLayout{})
	if buildFunc != nil {
		uiBuilder := newContainerBuilder(builder)
		buildFunc(uiBuilder)
	}
	return builder
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを作成するビルダーです。
func Grid(buildFunc func(*GridContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.GridLayout{
			Columns: 1, // デフォルトは1列
		})
	if buildFunc != nil {
		gridBuilder := newGridContainerBuilder(builder)
		buildFunc(gridBuilder)
	}
	return builder
}

// --- Generic Base Builder ---

// BaseBuilderは、uiパッケージ内のすべての具象コンテナビルダーの共通機能を提供します。
// これを埋め込むことで、具象ビルダーは子ウィジェット追加機能と、下位ビルダーのプロパティ設定機能の両方を、
// 一貫したメソッドチェーンで利用できます。
type BaseBuilder[T any] struct {
	*container.ContainerBuilder
	self T
}

// Init は、BaseBuilderを初期化します。具象ビルダーのコンストラクタから呼び出す必要があります。
func (b *BaseBuilder[T]) Init(self T, cb *container.ContainerBuilder) {
	b.self = self
	b.ContainerBuilder = cb
}

// --- Child Widget Adders ---

// Label はコンテナにLabelを追加します。
func (b *BaseBuilder[T]) Label(buildFunc func(*widget.LabelBuilder)) T {
	labelBuilder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(labelBuilder)
	}
	label, err := labelBuilder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(label)
	return b.self
}

// Button はコンテナにButtonを追加します。
func (b *BaseBuilder[T]) Button(buildFunc func(*widget.ButtonBuilder)) T {
	buttonBuilder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(buttonBuilder)
	}
	button, err := buttonBuilder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(button)
	return b.self
}

// Spacer はコンテナに伸縮可能な空白を追加します。FlexLayout内でのみ効果があります。
func (b *BaseBuilder[T]) Spacer() T {
	spacer, _ := widget.NewSpacerBuilder().Flex(1).Build()
	b.ContainerBuilder.AddChild(spacer)
	return b.self
}

// --- Nested Container Adders ---

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *BaseBuilder[T]) HStack(buildFunc func(*FlexContainerBuilder)) T {
	builder := HStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b.self
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *BaseBuilder[T]) VStack(buildFunc func(*FlexContainerBuilder)) T {
	builder := VStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b.self
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを追加します。
func (b *BaseBuilder[T]) ZStack(buildFunc func(*ContainerBuilder)) T {
	builder := ZStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b.self
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを追加します。
func (b *BaseBuilder[T]) Grid(buildFunc func(*GridContainerBuilder)) T {
	builder := Grid(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b.self
}

// --- Common Property Wrappers ---
// 以下のメソッドは container.ContainerBuilder のメソッドをラップし、
// メソッドチェーンが途切れないように具象ビルダー型(T)を返すためのものです。

// Size はコンテナのサイズを設定します。
func (b *BaseBuilder[T]) Size(width, height int) T {
	b.ContainerBuilder.Size(width, height)
	return b.self
}

// Style はコンテナのスタイルを設定します。既存のスタイルとマージされます。
func (b *BaseBuilder[T]) Style(s style.Style) T {
	b.ContainerBuilder.Style(s)
	return b.self
}

// Flex はコンテナのFlexLayoutにおける伸縮係数を設定します。
func (b *BaseBuilder[T]) Flex(flex int) T {
	b.ContainerBuilder.Flex(flex)
	return b.self
}

// Position はコンテナの絶対位置を設定します。ZStack内でのみ有効です。
func (b *BaseBuilder[T]) Position(x, y int) T {
	b.ContainerBuilder.AbsolutePosition(x, y)
	return b.self
}

// RelayoutBoundary はコンテナをレイアウト境界として設定します。
func (b *BaseBuilder[T]) RelayoutBoundary(isBoundary bool) T {
	b.ContainerBuilder.RelayoutBoundary(isBoundary)
	return b.self
}

// --- Style Helper Wrappers ---

// BackgroundColor はコンテナの背景色を設定します。
func (b *BaseBuilder[T]) BackgroundColor(c color.Color) T {
	b.ContainerBuilder.BackgroundColor(c)
	return b.self
}

// Margin はコンテナのマージンを四方すべてに同じ値で設定します。
func (b *BaseBuilder[T]) Margin(margin int) T {
	b.ContainerBuilder.Margin(margin)
	return b.self
}

// MarginInsets はコンテナのマージンを各辺個別に設定します。
func (b *BaseBuilder[T]) MarginInsets(insets style.Insets) T {
	b.ContainerBuilder.MarginInsets(insets)
	return b.self
}

// Padding はコンテナのパディングを四方すべてに同じ値で設定します。
func (b *BaseBuilder[T]) Padding(padding int) T {
	b.ContainerBuilder.Padding(padding)
	return b.self
}

// PaddingInsets はコンテナのパディングを各辺個別に設定します。
func (b *BaseBuilder[T]) PaddingInsets(insets style.Insets) T {
	b.ContainerBuilder.PaddingInsets(insets)
	return b.self
}

// BorderRadius はコンテナの角丸の半径を設定します。
func (b *BaseBuilder[T]) BorderRadius(radius float32) T {
	b.ContainerBuilder.BorderRadius(radius)
	return b.self
}

// Border はコンテナの境界線を設定します。
func (b *BaseBuilder[T]) Border(width float32, c color.Color) T {
	b.ContainerBuilder.Border(width, c)
	return b.self
}

// --- Concrete Builders ---

// ContainerBuilder は、特定のレイアウトに依存しない共通のコンテナビルダーです。
type ContainerBuilder struct {
	BaseBuilder[*ContainerBuilder]
}

func newContainerBuilder(cb *container.ContainerBuilder) *ContainerBuilder {
	b := &ContainerBuilder{}
	b.Init(b, cb)
	return b
}

// FlexContainerBuilder は、FlexLayoutを持つコンテナ（VStack, HStack）に特化した設定メソッドを提供します。
type FlexContainerBuilder struct {
	BaseBuilder[*FlexContainerBuilder]
}

func newFlexContainerBuilder(cb *container.ContainerBuilder) *FlexContainerBuilder {
	b := &FlexContainerBuilder{}
	b.Init(b, cb)
	return b
}

// --- FlexLayout-specific Methods ---

// Gap はFlexLayoutの子要素間の間隔を設定します。
func (b *FlexContainerBuilder) Gap(gap int) *FlexContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Gap != gap {
			flexLayout.Gap = gap
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// Justify はFlexLayoutの主軸方向の揃え位置を設定します。
func (b *FlexContainerBuilder) Justify(alignment layout.Alignment) *FlexContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Justify != alignment {
			flexLayout.Justify = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// AlignItems はFlexLayoutの交差軸方向の揃え位置を設定します。
func (b *FlexContainerBuilder) AlignItems(alignment layout.Alignment) *FlexContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.AlignItems != alignment {
			flexLayout.AlignItems = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// GridContainerBuilder はGridLayoutコンテナに特化した設定メソッドを提供します。
type GridContainerBuilder struct {
	BaseBuilder[*GridContainerBuilder]
}

func newGridContainerBuilder(cb *container.ContainerBuilder) *GridContainerBuilder {
	b := &GridContainerBuilder{}
	b.Init(b, cb)
	return b
}

// --- GridLayout-specific Methods ---

// Columns はグリッドの列数を設定します。
func (b *GridContainerBuilder) Columns(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if count > 0 && gridLayout.Columns != count {
			gridLayout.Columns = count
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Columns() can only be used on a Grid container"))
	}
	return b
}

// Rows はグリッドの行数を設定します。0以下で自動計算されます。
func (b *GridContainerBuilder) Rows(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.Rows != count {
			gridLayout.Rows = count
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Rows() can only be used on a Grid container"))
	}
	return b
}

// HorizontalGap はセル間の水平方向の間隔を設定します。
func (b *GridContainerBuilder) HorizontalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.HorizontalGap != gap {
			gridLayout.HorizontalGap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("HorizontalGap() can only be used on a Grid container"))
	}
	return b
}

// VerticalGap はセル間の垂直方向の間隔を設定します。
func (b *GridContainerBuilder) VerticalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.VerticalGap != gap {
			gridLayout.VerticalGap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("VerticalGap() can only be used on a Grid container"))
	}
	return b
}
