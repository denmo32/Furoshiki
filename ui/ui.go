package ui

import (
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"
	"image/color"
)

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです
func VStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction: layout.DirectionColumn,
		})

	if buildFunc != nil {
		containerBuilder := &ContainerBuilder{ContainerBuilder: builder}
		buildFunc(containerBuilder)
	}

	return builder
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成するビルダーです
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

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成するビルダーです
func ZStack(buildFunc func(*ContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.AbsoluteLayout{})

	if buildFunc != nil {
		containerBuilder := &ContainerBuilder{ContainerBuilder: builder}
		buildFunc(containerBuilder)
	}

	return builder
}

// [追加] Grid は子要素を格子状に配置するGridLayoutコンテナを作成するビルダーです
func Grid(buildFunc func(*GridContainerBuilder)) *container.ContainerBuilder {
	builder := container.NewContainerBuilder().
		Layout(&layout.GridLayout{
			Columns: 1, // デフォルトは1列
		})

	if buildFunc != nil {
		// GridContainerBuilderにラップして、専用のメソッドを提供します
		gridBuilder := &GridContainerBuilder{ContainerBuilder: &ContainerBuilder{ContainerBuilder: builder}}
		buildFunc(gridBuilder)
	}

	return builder
}

// ContainerBuilder はコンテナに子要素を追加するためのヘルパーメソッドを提供します
type ContainerBuilder struct {
	*container.ContainerBuilder
}

// Label はコンテナにLabelを追加します
// [改善] 子ウィジェットのビルド時に発生したエラーを親のビルダーに伝播させます。
func (b *ContainerBuilder) Label(buildFunc func(*widget.LabelBuilder)) *ContainerBuilder {
	labelBuilder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(labelBuilder)
	}
	label, err := labelBuilder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(label)
	return b
}

// Button はコンテナにButtonを追加します
// [改善] 子ウィジェットのビルド時に発生したエラーを親のビルダーに伝播させます。
func (b *ContainerBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *ContainerBuilder {
	buttonBuilder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(buttonBuilder)
	}
	button, err := buttonBuilder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(button)
	return b
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを追加します
// [改善] 子コンテナのビルド時に発生したエラーを親のビルダーに伝播させます。
func (b *ContainerBuilder) HStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := HStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(container)
	return b
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを追加します
// [改善] 子コンテナのビルド時に発生したエラーを親のビルダーに伝播させます。
func (b *ContainerBuilder) VStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := VStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(container)
	return b
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを追加します
// [改善] 子コンテナのビルド時に発生したエラーを親のビルダーに伝播させます。
func (b *ContainerBuilder) ZStack(buildFunc func(*ContainerBuilder)) *ContainerBuilder {
	builder := ZStack(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(container)
	return b
}

// [追加] Grid は子要素を格子状に配置するGridLayoutコンテナを追加します
func (b *ContainerBuilder) Grid(buildFunc func(*GridContainerBuilder)) *ContainerBuilder {
	builder := Grid(buildFunc)
	container, err := builder.Build()
	b.ContainerBuilder.AddError(err) // エラーを記録
	b.ContainerBuilder.AddChild(container)
	return b
}

// Padding はコンテナのパディングを設定します
// [修正] style.Insetsをポインタで設定するように変更
func (b *ContainerBuilder) Padding(padding int) *ContainerBuilder {
	b.ContainerBuilder.Style(style.Style{
		Padding: &style.Insets{
			Top:    padding,
			Right:  padding,
			Bottom: padding,
			Left:   padding,
		},
	})
	return b
}

// Gap はFlexLayoutの子要素間の間隔を設定します
// [改善] コンテナをビルドせずに、ビルダーから直接レイアウトオブジェクトを取得して変更します。
// これにより、不要な処理が削減され、効率が向上します。
func (b *ContainerBuilder) Gap(gap int) *ContainerBuilder {
	// 型アサーションでレイアウトがFlexLayoutか確認
	if flexLayout, ok := b.ContainerBuilder.GetLayout().(*layout.FlexLayout); ok {
		flexLayout.Gap = gap
		// レイアウトオブジェクトはポインタなので、直接変更すればコンテナに反映されます。
	}
	return b
}

// Justify はFlexLayoutの主軸方向の揃え位置を設定します
// [改善] コンテナをビルドせずに、ビルダーから直接レイアウトオブジェクトを取得して変更します。
func (b *ContainerBuilder) Justify(alignment layout.Alignment) *ContainerBuilder {
	if flexLayout, ok := b.ContainerBuilder.GetLayout().(*layout.FlexLayout); ok {
		flexLayout.Justify = alignment
	}
	return b
}

// AlignItems はFlexLayoutの交差軸方向の揃え位置を設定します
// [改善] コンテナをビルドせずに、ビルダーから直接レイアウトオブジェクトを取得して変更します。
func (b *ContainerBuilder) AlignItems(alignment layout.Alignment) *ContainerBuilder {
	if flexLayout, ok := b.ContainerBuilder.GetLayout().(*layout.FlexLayout); ok {
		flexLayout.AlignItems = alignment
	}
	return b
}

// Size はコンテナのサイズを設定します
func (b *ContainerBuilder) Size(width, height int) *ContainerBuilder {
	b.ContainerBuilder.Size(width, height)
	return b
}

// Style はコンテナのスタイルを設定します
func (b *ContainerBuilder) Style(s style.Style) *ContainerBuilder {
	b.ContainerBuilder.Style(s)
	return b
}

// Style は一般的に使用されるスタイルを簡単に作成するためのヘルパー関数です
// [修正] style.Styleのフィールドがポインタになったため、&でアドレスを渡すように変更
func Style(background color.Color, textColor color.Color) style.Style {
	return style.Style{
		Background: &background,
		TextColor:  &textColor,
		Padding: &style.Insets{
			Top: 5, Right: 10, Bottom: 5, Left: 10,
		},
	}
}

// --- GridContainerBuilder ---

// [追加] GridContainerBuilder はGridLayoutコンテナに特化した設定メソッドを提供します。
// ContainerBuilderをラップすることで、既存のメソッド（Size, Styleなど）も利用可能にします。
type GridContainerBuilder struct {
	*ContainerBuilder
}

// Columns はグリッドの列数を設定します。
func (b *GridContainerBuilder) Columns(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if count > 0 {
			gridLayout.Columns = count
		}
	}
	return b
}

// Rows はグリッドの行数を設定します。
// 0または負の値を設定すると、子の数と列数から自動的に計算されます。
func (b *GridContainerBuilder) Rows(count int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		gridLayout.Rows = count
	}
	return b
}

// HorizontalGap はセル間の水平方向の間隔を設定します。
func (b *GridContainerBuilder) HorizontalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		gridLayout.HorizontalGap = gap
	}
	return b
}

// VerticalGap はセル間の垂直方向の間隔を設定します。
func (b *GridContainerBuilder) VerticalGap(gap int) *GridContainerBuilder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		gridLayout.VerticalGap = gap
	}
	return b
}

// [修正] 以下のメソッドは、メソッドチェーンが正しく機能するように戻り値の型を *GridContainerBuilder にオーバーライドします。

// Size はコンテナのサイズを設定します。
func (b *GridContainerBuilder) Size(width, height int) *GridContainerBuilder {
	b.ContainerBuilder.Size(width, height)
	return b
}

// Style はコンテナのスタイルを設定します。
func (b *GridContainerBuilder) Style(s style.Style) *GridContainerBuilder {
	b.ContainerBuilder.Style(s)
	return b
}

// Padding はコンテナのパディングを設定します。
func (b *GridContainerBuilder) Padding(padding int) *GridContainerBuilder {
	b.ContainerBuilder.Padding(padding)
	return b
}

// Button はコンテナにButtonを追加します。
func (b *GridContainerBuilder) Button(buildFunc func(*widget.ButtonBuilder)) *GridContainerBuilder {
	b.ContainerBuilder.Button(buildFunc)
	return b
}

// Label はコンテナにLabelを追加します。
func (b *GridContainerBuilder) Label(buildFunc func(*widget.LabelBuilder)) *GridContainerBuilder {
	b.ContainerBuilder.Label(buildFunc)
	return b
}