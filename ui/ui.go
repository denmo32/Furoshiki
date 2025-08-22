package ui

import (
	"fmt"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"
	"image/color"
)

// newContainer は、指定されたレイアウトでコンテナビルダーを生成し、
// ユーザー提供のビルド関数を実行する内部ヘルパーです。
// これにより、VStack, HStack, ZStack といった類似のコンテナ生成関数のコード重複を排除します。
func newContainer(layout layout.Layout, buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	// containerパッケージの基本的なビルダーを作成し、指定されたレイアウトを設定します。
	builder := container.NewContainerBuilder().Layout(layout)

	if buildFunc != nil {
		// uiパッケージの高レベルなラッパービルダーを生成し、ユーザー提供の関数を実行します。
		// これにより、Label()やButton()などの便利なヘルパーメソッドが利用可能になります。
		uiBuilder := newContainerBuilder(builder)
		buildFunc(uiBuilder)
	}

	return builder
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func VStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	// 共通ヘルパー関数を、垂直方向のFlexLayoutで呼び出します。
	return newContainer(
		&layout.FlexLayout{
			Direction: layout.DirectionColumn,
		},
		buildFunc,
	)
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func HStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	// 共通ヘルパー関数を、水平方向のFlexLayoutで呼び出します。
	return newContainer(
		&layout.FlexLayout{
			Direction: layout.DirectionRow,
		},
		buildFunc,
	)
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成するビルダーです。
func ZStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	// 共通ヘルパー関数を、AbsoluteLayoutで呼び出します。
	return newContainer(&layout.AbsoluteLayout{}, buildFunc)
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを作成するビルダーです。
// この関数は buildFunc の引数の型が異なるため、newContainer ヘルパーは使用しません。
func Grid(buildFunc func(*GridContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.GridLayout{
			Columns: 1, // デフォルトは1列
		})

	if buildFunc != nil {
		// GridContainerBuilderにラップして、グリッド専用のメソッドを提供します。
		gridBuilder := newGridContainerBuilder(builder)
		buildFunc(gridBuilder)
	}

	return builder
}

// --- Generic Base Builder ---

// BaseBuilder は、uiパッケージ内のコンテナビルダー（ContainerBuilder, GridContainerBuilder）の
// 共通機能をまとめたジェネリックな基底構造体です。
// ジェネリクスの型パラメータ T は、この構造体を埋め込む具象ビルダー自身の型（例: *ContainerBuilder）を指定します。
// これにより、メソッドチェーンを維持したまま具象ビルダーの型を返すことが可能になり、コードの冗長性を排除します。
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
// デフォルトでFlex値が1に設定され、利用可能なスペースをすべて埋めようとします。
func (b *BaseBuilder[T]) Spacer() T {
	spacer, _ := widget.NewSpacerBuilder().Flex(1).Build()
	b.ContainerBuilder.AddChild(spacer)
	return b.self
}

// --- Nested Container Adders ---

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *BaseBuilder[T]) HStack(buildFunc func(*ContainerBuilder)) T {
	builder := HStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b.self
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *BaseBuilder[T]) VStack(buildFunc func(*ContainerBuilder)) T {
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

// --- Layout Property Helpers ---

// Gap はFlexLayoutの子要素間の間隔を設定します。
func (b *BaseBuilder[T]) Gap(gap int) T {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Gap != gap {
			flexLayout.Gap = gap
			b.Widget.MarkDirty(true)
		}
	}
	return b.self
}

// Justify はFlexLayoutの主軸方向の揃え位置を設定します。
func (b *BaseBuilder[T]) Justify(alignment layout.Alignment) T {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Justify != alignment {
			flexLayout.Justify = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b.self
}

// AlignItems はFlexLayoutの交差軸方向の揃え位置を設定します。
func (b *BaseBuilder[T]) AlignItems(alignment layout.Alignment) T {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.AlignItems != alignment {
			flexLayout.AlignItems = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b.self
}

// --- Common Property Helpers ---

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

// Position はコンテナの絶対位置を設定します。FlexLayout内では上書きされる可能性があります。
// component.BuilderのAbsolutePositionを呼び出します。
func (b *BaseBuilder[T]) Position(x, y int) T {
	b.ContainerBuilder.AbsolutePosition(x, y)
	return b.self
}

// --- Style Helpers ---

// BackgroundColor はコンテナの背景色を設定します。
func (b *BaseBuilder[T]) BackgroundColor(c color.Color) T {
	return b.Style(style.Style{Background: style.PColor(c)})
}

// Margin はコンテナのマージンを四方すべてに同じ値で設定します。
func (b *BaseBuilder[T]) Margin(margin int) T {
	return b.Style(style.Style{Margin: style.PInsets(style.Insets{
		Top: margin, Right: margin, Bottom: margin, Left: margin,
	})})
}

// MarginInsets はコンテナのマージンを各辺個別に設定します。
func (b *BaseBuilder[T]) MarginInsets(insets style.Insets) T {
	return b.Style(style.Style{Margin: style.PInsets(insets)})
}

// Padding はコンテナのパディングを四方すべてに同じ値で設定します。
func (b *BaseBuilder[T]) Padding(padding int) T {
	return b.Style(style.Style{Padding: style.PInsets(style.Insets{
		Top: padding, Right: padding, Bottom: padding, Left: padding,
	})})
}

// PaddingInsets はコンテナのパディングを各辺個別に設定します。
func (b *BaseBuilder[T]) PaddingInsets(insets style.Insets) T {
	return b.Style(style.Style{Padding: style.PInsets(insets)})
}

// BorderRadius はコンテナの角丸の半径を設定します。
func (b *BaseBuilder[T]) BorderRadius(radius float32) T {
	return b.Style(style.Style{BorderRadius: style.PFloat32(radius)})
}

// Border はコンテナの境界線を設定します。
func (b *BaseBuilder[T]) Border(width float32, c color.Color) T {
	if width < 0 {
		b.AddError(fmt.Errorf("border width must be non-negative, got %f", width))
		return b.self
	}
	return b.Style(style.Style{
		BorderWidth: style.PFloat32(width),
		BorderColor: style.PColor(c),
	})
}

// --- Concrete Builders ---

// ContainerBuilder はコンテナに子要素を追加するためのヘルパーメソッドを提供します。
// BaseBuilderを埋め込むことで、共通のビルダーメソッドを利用します。
type ContainerBuilder struct {
	BaseBuilder[*ContainerBuilder]
}

// newContainerBuilder は、内部で使用する新しいContainerBuilderを生成します。
func newContainerBuilder(cb *container.ContainerBuilder) *ContainerBuilder {
	b := &ContainerBuilder{}
	b.Init(b, cb)
	return b
}

// GridContainerBuilder はGridLayoutコンテナに特化した設定メソッドを提供します。
// BaseBuilderを埋め込むことで、共通のビルダーメソッドを利用します。
type GridContainerBuilder struct {
	BaseBuilder[*GridContainerBuilder]
}

// newGridContainerBuilder は、内部で使用する新しいGridContainerBuilderを生成します。
func newGridContainerBuilder(cb *container.ContainerBuilder) *GridContainerBuilder {
	b := &GridContainerBuilder{}
	b.Init(b, cb)
	return b
}

// --- Grid-specific Methods ---

// Columns はグリッドの列数を設定します。
func (b *GridContainerBuilder) Columns(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if count > 0 && gridLayout.Columns != count {
			gridLayout.Columns = count
			b.Widget.MarkDirty(true)
		}
	} else {
		// レイアウトがGridLayoutでない場合にエラーを追加する方が親切です。
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