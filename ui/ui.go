package ui

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"
	"image/color"
)

// Builder は、UI構造を宣言的に構築するための統一されたインターフェースです。
// メソッドチェーンをサポートし、VStack, HStack, Gridなどのさまざまなコンテナタイプで
// 一貫した操作を提供します。
type Builder struct {
	*container.ContainerBuilder
}

// newBuilder は内部の container.ContainerBuilder をラップして新しい ui.Builder を作成します。
func newBuilder(cb *container.ContainerBuilder) *Builder {
	return &Builder{ContainerBuilder: cb}
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成します。
func VStack(buildFunc func(*Builder)) *Builder {
	cb := container.NewContainerBuilder().Layout(&layout.FlexLayout{
		Direction: layout.DirectionColumn,
	})
	builder := newBuilder(cb)
	if buildFunc != nil {
		buildFunc(builder)
	}
	return builder
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成します。
func HStack(buildFunc func(*Builder)) *Builder {
	cb := container.NewContainerBuilder().Layout(&layout.FlexLayout{
		Direction: layout.DirectionRow,
	})
	builder := newBuilder(cb)
	if buildFunc != nil {
		buildFunc(builder)
	}
	return builder
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成します。
func ZStack(buildFunc func(*Builder)) *Builder {
	cb := container.NewContainerBuilder().Layout(&layout.AbsoluteLayout{})
	builder := newBuilder(cb)
	if buildFunc != nil {
		buildFunc(builder)
	}
	return builder
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを作成します。
func Grid(buildFunc func(*Builder)) *Builder {
	cb := container.NewContainerBuilder().Layout(&layout.GridLayout{Columns: 1})
	builder := newBuilder(cb)
	if buildFunc != nil {
		buildFunc(builder)
	}
	return builder
}

// --- Child Widget Adders ---

// Label はコンテナにLabelを追加します。
func (b *Builder) Label(buildFunc func(*widget.LabelBuilder)) *Builder {
	labelBuilder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(labelBuilder)
	}
	label, err := labelBuilder.Build()
	b.AddError(err)
	b.AddChild(label)
	return b
}

// Button はコンテナにButtonを追加します。
func (b *Builder) Button(buildFunc func(*widget.ButtonBuilder)) *Builder {
	buttonBuilder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(buttonBuilder)
	}
	button, err := buttonBuilder.Build()
	b.AddError(err)
	b.AddChild(button)
	return b
}

// Spacer はコンテナに伸縮可能な空白を追加します。FlexLayout内でのみ効果があります。
func (b *Builder) Spacer() *Builder {
	spacer, _ := widget.NewSpacerBuilder().Flex(1).Build()
	b.AddChild(spacer)
	return b
}

// AddChild はコンテナに子ウィジェットを追加します。
func (b *Builder) AddChild(child component.Widget) *Builder {
	b.ContainerBuilder.AddChild(child)
	return b
}

// AddChildren はコンテナに複数の子ウィジェットを追加します。
func (b *Builder) AddChildren(children ...component.Widget) *Builder {
	b.ContainerBuilder.AddChildren(children...)
	return b
}

// --- Nested Container Adders ---

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *Builder) HStack(buildFunc func(*Builder)) *Builder {
	builder := HStack(buildFunc)
	container, err := builder.Build()
	b.AddError(err)
	b.AddChild(container)
	return b
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを追加します。
func (b *Builder) VStack(buildFunc func(*Builder)) *Builder {
	builder := VStack(buildFunc)
	container, err := builder.Build()
	b.AddError(err)
	b.AddChild(container)
	return b
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを追加します。
func (b *Builder) ZStack(buildFunc func(*Builder)) *Builder {
	builder := ZStack(buildFunc)
	container, err := builder.Build()
	b.AddError(err)
	b.AddChild(container)
	return b
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを追加します。
func (b *Builder) Grid(buildFunc func(*Builder)) *Builder {
	builder := Grid(buildFunc)
	container, err := builder.Build()
	b.AddError(err)
	b.AddChild(container)
	return b
}

// --- Common Property Wrappers ---

// Size はコンテナのサイズを設定します。
func (b *Builder) Size(width, height int) *Builder {
	b.ContainerBuilder.Size(width, height)
	return b
}

// MinSize はコンテナの最小サイズを設定します。
func (b *Builder) MinSize(width, height int) *Builder {
	b.ContainerBuilder.MinSize(width, height)
	return b
}

// Style はコンテナのスタイルを設定します。既存のスタイルとマージされます。
func (b *Builder) Style(s style.Style) *Builder {
	b.ContainerBuilder.Style(s)
	return b
}

// Flex はコンテナのFlexLayoutにおける伸縮係数を設定します。
func (b *Builder) Flex(flex int) *Builder {
	b.ContainerBuilder.Flex(flex)
	return b
}

// AbsolutePosition はコンテナの絶対位置を設定します。ZStack内でのみ有効です。
func (b *Builder) AbsolutePosition(x, y int) *Builder {
	b.ContainerBuilder.AbsolutePosition(x, y)
	return b
}

// RelayoutBoundary はコンテナをレイアウト境界として設定します。
func (b *Builder) RelayoutBoundary(isBoundary bool) *Builder {
	b.ContainerBuilder.RelayoutBoundary(isBoundary)
	return b
}

// [新規追加]
// AssignTo は、ビルド中のコンテナインスタンスへのポインタを変数に代入します。
// UIの宣言的な構築フローを維持したまま、特定のコンテナへの参照を取得するために使用します。
// 例: .VStack(func(b *ui.Builder){ b.AssignTo(&myContainer) }) (myContainerは *container.Container 型)
func (b *Builder) AssignTo(target any) *Builder {
	// 埋め込まれたContainerBuilderが持つAssignToメソッドを呼び出します。
	// このメソッドはcomponent.Builderに実装されており、継承を通じて利用可能です。
	b.ContainerBuilder.AssignTo(target)
	// メソッドチェーンを継続するために、*Builder型であるb自身を返します。
	return b
}

// --- Style Helper Wrappers ---

// BackgroundColor はコンテナの背景色を設定します。
func (b *Builder) BackgroundColor(c color.Color) *Builder {
	b.ContainerBuilder.BackgroundColor(c)
	return b
}

// Margin はコンテナのマージンを四方すべてに同じ値で設定します。
func (b *Builder) Margin(margin int) *Builder {
	b.ContainerBuilder.Margin(margin)
	return b
}

// MarginInsets はコンテナのマージンを各辺個別に設定します。
func (b *Builder) MarginInsets(insets style.Insets) *Builder {
	b.ContainerBuilder.MarginInsets(insets)
	return b
}

// Padding はコンテナのパディングを四方すべてに同じ値で設定します。
func (b *Builder) Padding(padding int) *Builder {
	b.ContainerBuilder.Padding(padding)
	return b
}

// PaddingInsets はコンテナのパディングを各辺個別に設定します。
func (b *Builder) PaddingInsets(insets style.Insets) *Builder {
	b.ContainerBuilder.PaddingInsets(insets)
	return b
}

// BorderRadius はコンテナの角丸の半径を設定します。
func (b *Builder) BorderRadius(radius float32) *Builder {
	b.ContainerBuilder.BorderRadius(radius)
	return b
}

// Border はコンテナの境界線を設定します。
func (b *Builder) Border(width float32, c color.Color) *Builder {
	b.ContainerBuilder.Border(width, c)
	return b
}

// --- FlexLayout Specific Methods ---

// Gap はFlexLayoutの子要素間の間隔を設定します。
// VStackまたはHStack内でのみ有効です。
func (b *Builder) Gap(gap int) *Builder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Gap != gap {
			flexLayout.Gap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Gap() is only applicable to FlexLayout (VStack, HStack)"))
	}
	return b
}

// Justify はFlexLayoutの主軸方向の揃え位置を設定します。
// VStackまたはHStack内でのみ有効です。
func (b *Builder) Justify(alignment layout.Alignment) *Builder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Justify != alignment {
			flexLayout.Justify = alignment
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Justify() is only applicable to FlexLayout (VStack, HStack)"))
	}
	return b
}

// AlignItems はFlexLayoutの交差軸方向の揃え位置を設定します。
// VStackまたはHStack内でのみ有効です。
func (b *Builder) AlignItems(alignment layout.Alignment) *Builder {
	if flexLayout, ok := b.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.AlignItems != alignment {
			flexLayout.AlignItems = alignment
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("AlignItems() is only applicable to FlexLayout (VStack, HStack)"))
	}
	return b
}

// --- GridLayout Specific Methods ---

// Columns はグリッドの列数を設定します。
// Grid内でのみ有効です。
func (b *Builder) Columns(count int) *Builder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if count > 0 && gridLayout.Columns != count {
			gridLayout.Columns = count
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Columns() is only applicable to GridLayout"))
	}
	return b
}

// Rows はグリッドの行数を設定します。0以下で自動計算されます。
// Grid内でのみ有効です。
func (b *Builder) Rows(count int) *Builder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.Rows != count {
			gridLayout.Rows = count
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("Rows() is only applicable to GridLayout"))
	}
	return b
}

// HorizontalGap はセル間の水平方向の間隔を設定します。
// Grid内でのみ有効です。
func (b *Builder) HorizontalGap(gap int) *Builder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.HorizontalGap != gap {
			gridLayout.HorizontalGap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("HorizontalGap() is only applicable to GridLayout"))
	}
	return b
}

// VerticalGap はセル間の垂直方向の間隔を設定します。
// Grid内でのみ有効です。
func (b *Builder) VerticalGap(gap int) *Builder {
	if gridLayout, ok := b.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.VerticalGap != gap {
			gridLayout.VerticalGap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("VerticalGap() is only applicable to GridLayout"))
	}
	return b
}