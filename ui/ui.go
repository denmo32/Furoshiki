package ui

import (
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"
	"image/color"
)

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func VStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction: layout.DirectionColumn,
		})

	if buildFunc != nil {
		// ContainerBuilderでラップして、便利なヘルパーメソッドを提供します。
		containerBuilder := &ContainerBuilder{ContainerBuilder: builder}
		buildFunc(containerBuilder)
	}

	return builder
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです。
func HStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction: layout.DirectionRow,
		})

	if buildFunc != nil {
		containerBuilder := &ContainerBuilder{ContainerBuilder: builder}
		buildFunc(containerBuilder)
	}

	return builder
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成するビルダーです。
func ZStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.AbsoluteLayout{})

	if buildFunc != nil {
		containerBuilder := &ContainerBuilder{ContainerBuilder: builder}
		buildFunc(containerBuilder)
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
		// GridContainerBuilderにラップして、グリッド専用のメソッドを提供します。
		gridBuilder := &GridContainerBuilder{ContainerBuilder: &ContainerBuilder{ContainerBuilder: builder}}
		buildFunc(gridBuilder)
	}

	return builder
}

// ContainerBuilder はコンテナに子要素を追加するためのヘルパーメソッドを提供します。
type ContainerBuilder struct {
	*container.ContainerBuilder
}

// Label はコンテナにLabelを追加します。
func (b *ContainerBuilder) Label(buildFunc func(*widget.LabelBuilder)) *ContainerBuilder {
	labelBuilder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(labelBuilder)
	}
	label, err := labelBuilder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(label)
	return b
}

// Button はコンテナにButtonを追加します。
func (b *ContainerBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *ContainerBuilder {
	buttonBuilder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(buttonBuilder)
	}
	button, err := buttonBuilder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(button)
	return b
}

// [新機能] Spacer はコンテナに伸縮可能な空白を追加します。FlexLayout内でのみ効果があります。
// デフォルトでFlex値が1に設定され、利用可能なスペースをすべて埋めようとします。
func (b *ContainerBuilder) Spacer() *ContainerBuilder {
	spacer, _ := widget.NewSpacerBuilder().Flex(1).Build()
	b.ContainerBuilder.AddChild(spacer)
	return b
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *ContainerBuilder) HStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := HStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *ContainerBuilder) VStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := VStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを追加します。
func (b *ContainerBuilder) ZStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := ZStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを追加します。
func (b *ContainerBuilder) Grid(buildFunc func(*GridContainerBuilder)) *ContainerBuilder {
	builder := Grid(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err)
	b.ContainerBuilder.AddChild(container)
	return b
}

// --- Layout Property Helpers ---

// Gap はFlexLayoutの子要素間の間隔を設定します。
func (b *ContainerBuilder) Gap(gap int) *ContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Gap != gap {
			flexLayout.Gap = gap
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// Justify はFlexLayoutの主軸方向の揃え位置を設定します。
func (b *ContainerBuilder) Justify(alignment layout.Alignment) *ContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Justify != alignment {
			flexLayout.Justify = alignment
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// AlignItems はFlexLayoutの交差軸方向の揃え位置を設定します。
func (b *ContainerBuilder) AlignItems(alignment layout.Alignment) *ContainerBuilder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.AlignItems != alignment {
			flexLayout.AlignItems = alignment
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// --- Common Property Helpers ---
// [改良] メソッドチェーンが途切れないように、戻り値の型を *ContainerBuilder にします。

// Size はコンテナのサイズを設定します。
func (b *ContainerBuilder) Size(width, height int) *ContainerBuilder {
	b.ContainerBuilder.Size(width, height)
	return b
}

// Style はコンテナのスタイルを設定します。既存のスタイルとマージされます。
func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
	b.ContainerBuilder.Style(s)
	return b
}

// Flex はコンテナのFlexLayoutにおける伸縮係数を設定します。
func (b *ContainerBuilder) Flex(flex int) *ContainerBuilder {
	b.ContainerBuilder.Flex(flex)
	return b
}

// Position はコンテナの絶対位置を設定します。FlexLayout内では上書きされる可能性があります。
func (b *ContainerBuilder) Position(x, y int) *ContainerBuilder {
	b.ContainerBuilder.Position(x, y)
	return b
}

// --- Style Helpers ---
// [新機能] スタイルをより直感的に設定するためのヘルパーメソッド群です。

// BackgroundColor はコンテナの背景色を設定します。
func (b *ContainerBuilder) BackgroundColor(c color.Color) *ContainerBuilder {
	return b.Style(style.Style{Background: style.PColor(c)})
}

// Padding はコンテナのパディングを四方すべてに同じ値で設定します。
func (b *ContainerBuilder) Padding(padding int) *ContainerBuilder {
	return b.Style(style.Style{Padding: style.PInsets(style.Insets{
		Top: padding, Right: padding, Bottom: padding, Left: padding,
	})})
}

// PaddingInsets はコンテナのパディングを各辺個別に設定します。
func (b *ContainerBuilder) PaddingInsets(insets style.Insets) *ContainerBuilder {
	return b.Style(style.Style{Padding: style.PInsets(insets)})
}

// BorderRadius はコンテナの角丸の半径を設定します。
func (b *ContainerBuilder) BorderRadius(radius float32) *ContainerBuilder {
	return b.Style(style.Style{BorderRadius: style.PFloat32(radius)})
}

// StyleHelper は一般的に使用されるスタイルを簡単に作成するためのヘルパー関数です。
func StyleHelper(background color.Color, textColor color.Color) style.Style {
	return style.Style{
		Background: style.PColor(background),
		TextColor:  style.PColor(textColor),
		Padding: style.PInsets(style.Insets{
			Top: 5, Right: 10, Bottom: 5, Left: 10,
		}),
	}
}

// --- GridContainerBuilder ---

// GridContainerBuilder はGridLayoutコンテナに特化した設定メソッドを提供します。
type GridContainerBuilder struct {
	*ContainerBuilder
}

// Columns はグリッドの列数を設定します。
func (b *GridContainerBuilder) Columns(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if count > 0 && gridLayout.Columns != count {
			gridLayout.Columns = count
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// Rows はグリッドの行数を設定します。0以下で自動計算されます。
func (b *GridContainerBuilder) Rows(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.Rows != count {
			gridLayout.Rows = count
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// HorizontalGap はセル間の水平方向の間隔を設定します。
func (b *GridContainerBuilder) HorizontalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.HorizontalGap != gap {
			gridLayout.HorizontalGap = gap
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// VerticalGap はセル間の垂直方向の間隔を設定します。
func (b *GridContainerBuilder) VerticalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.VerticalGap != gap {
			gridLayout.VerticalGap = gap
			// [修正] レイアウトプロパティの変更は再レイアウトを必要とするため、ダーティフラグを立てます。
			b.MarkDirty(true)
		}
	}
	return b
}

// [改良] メソッドチェーンがGridContainerBuilderを返すようにオーバーライドします。
func (b *GridContainerBuilder) Size(width, height int) *GridContainerBuilder {
	b.ContainerBuilder.Size(width, height)
	return b
}
func (b *GridContainerBuilder) Style(s style.Style) *GridContainerBuilder {
	b.ContainerBuilder.Style(s)
	return b
}
func (b *GridContainerBuilder) Padding(padding int) *GridContainerBuilder {
	b.ContainerBuilder.Padding(padding)
	return b
}
func (b *GridContainerBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *GridContainerBuilder {
	b.ContainerBuilder.Button(buildFunc)
	return b
}
func (b *GridContainerBuilder) Label(buildFunc func(*widget.LabelBuilder)) *GridContainerBuilder {
	b.ContainerBuilder.Label(buildFunc)
	return b
}