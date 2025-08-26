package ui

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/widget"
	"reflect" // reflectパッケージをインポート
)

// このファイルは、レイアウトごとに特化した型付きビルダーを提供します。
// これらのビルダーは、新しく導入された `BaseContainerBuilder` を埋め込むことで、
// ウィジェット追加などの共通ロジックを再利用し、コードの冗長性を排除しています。

// --- FlexBuilder (VStack, HStack用) ---

// FlexBuilder は、FlexLayoutを持つコンテナを構築するためのビルダーです。
type FlexBuilder struct {
	*BaseContainerBuilder[*FlexBuilder]
}

// VStack は子要素を垂直方向に配置するコンテナを構築します。
func VStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	flexLayout := &layout.FlexLayout{Direction: layout.DirectionColumn}
	return buildFlexContainer(flexLayout, buildFunc)
}

// HStack は子要素を水平方向に配置するコンテナを構築します。
func HStack(buildFunc func(*FlexBuilder)) *FlexBuilder {
	flexLayout := &layout.FlexLayout{Direction: layout.DirectionRow}
	return buildFlexContainer(flexLayout, buildFunc)
}

// buildFlexContainer はFlexBuilderのインスタンスを生成する内部ヘルパーです。
func buildFlexContainer(l *layout.FlexLayout, buildFunc func(*FlexBuilder)) *FlexBuilder {
	c := container.NewContainer()
	c.SetLayout(l)

	b := &FlexBuilder{
		BaseContainerBuilder: &BaseContainerBuilder[*FlexBuilder]{},
	}
	b.Init(b, c) // BaseContainerBuilderのInitを呼び出す

	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// Gap は、FlexLayout内の子要素間の間隔を設定します。
func (b *FlexBuilder) Gap(gap int) *FlexBuilder {
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Gap != gap {
			flexLayout.Gap = gap
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// Justify は、FlexLayoutの主軸方向の揃え位置を設定します。
func (b *FlexBuilder) Justify(alignment layout.Alignment) *FlexBuilder {
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.Justify != alignment {
			flexLayout.Justify = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// AlignItems は、FlexLayoutの交差軸方向の揃え位置を設定します。
func (b *FlexBuilder) AlignItems(alignment layout.Alignment) *FlexBuilder {
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
		if flexLayout.AlignItems != alignment {
			flexLayout.AlignItems = alignment
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// Build はコンテナの構築を完了します。
func (b *FlexBuilder) Build() (*container.Container, error) { return b.Builder.Build() }

// --- GridBuilder (Grid用) ---

// GridBuilder は、GridLayoutを持つコンテナを構築するためのビルダーです。
type GridBuilder struct {
	*BaseContainerBuilder[*GridBuilder]
}

// Grid は子要素を格子状に配置するコンテナを構築します。
func Grid(buildFunc func(*GridBuilder)) *GridBuilder {
	gridLayout := &layout.GridLayout{Columns: 1} // デフォルトは1列
	c := container.NewContainer()
	c.SetLayout(gridLayout)

	b := &GridBuilder{
		BaseContainerBuilder: &BaseContainerBuilder[*GridBuilder]{},
	}
	b.Init(b, c)

	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// Columns は、グリッドの列数を設定します。
func (b *GridBuilder) Columns(count int) *GridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
		if count > 0 && gridLayout.Columns != count {
			gridLayout.Columns = count
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// Rows は、グリッドの行数を設定します。
func (b *GridBuilder) Rows(count int) *GridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.Rows != count {
			gridLayout.Rows = count
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// HorizontalGap は、グリッドの水平方向の間隔を設定します。
func (b *GridBuilder) HorizontalGap(gap int) *GridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.HorizontalGap != gap {
			gridLayout.HorizontalGap = gap
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// VerticalGap は、グリッドの垂直方向の間隔を設定します。
func (b *GridBuilder) VerticalGap(gap int) *GridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.VerticalGap != gap {
			gridLayout.VerticalGap = gap
			b.Widget.MarkDirty(true)
		}
	}
	return b
}

// Build はコンテナの構築を完了します。
func (b *GridBuilder) Build() (*container.Container, error) { return b.Builder.Build() }

// --- ZStackBuilder (ZStack用) ---

// ZStackBuilder は、AbsoluteLayoutを持つコンテナを構築するためのビルダーです。
type ZStackBuilder struct {
	*BaseContainerBuilder[*ZStackBuilder]
}

// ZStack は子要素を重ねて配置するコンテナを構築します。
func ZStack(buildFunc func(*ZStackBuilder)) *ZStackBuilder {
	c := container.NewContainer()
	c.SetLayout(&layout.AbsoluteLayout{})

	b := &ZStackBuilder{
		BaseContainerBuilder: &BaseContainerBuilder[*ZStackBuilder]{},
	}
	b.Init(b, c)

	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// Build はコンテナの構築を完了します。
func (b *ZStackBuilder) Build() (*container.Container, error) { return b.Builder.Build() }

// --- AdvancedGridBuilder ---

// AdvancedGridBuilder は、AdvancedGridLayoutを持つコンテナを構築するためのビルダーです。
type AdvancedGridBuilder struct {
	*BaseContainerBuilder[*AdvancedGridBuilder]
}

// AdvancedGrid は、セル結合や可変サイズの列・行を持つ高度なグリッドコンテナを構築します。
func AdvancedGrid(buildFunc func(*AdvancedGridBuilder)) *AdvancedGridBuilder {
	gridLayout := &layout.AdvancedGridLayout{}
	c := container.NewContainer()
	c.SetLayout(gridLayout)

	b := &AdvancedGridBuilder{
		BaseContainerBuilder: &BaseContainerBuilder[*AdvancedGridBuilder]{},
	}
	b.Init(b, c)

	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// Fixed は、固定ピクセルサイズのトラック定義を返します。
func Fixed(pixels float64) layout.TrackDefinition {
	return layout.TrackDefinition{Sizing: layout.TrackSizingFixed, Value: pixels}
}

// Weight は、重み付けによる可変サイズのトラック定義を返します。
func Weight(weight float64) layout.TrackDefinition {
	return layout.TrackDefinition{Sizing: layout.TrackSizingWeighted, Value: weight}
}

// Columns は、グリッドの列定義を設定します。
func (b *AdvancedGridBuilder) Columns(defs ...layout.TrackDefinition) *AdvancedGridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.AdvancedGridLayout); ok {
		gridLayout.ColumnDefinitions = defs
		b.Widget.MarkDirty(true)
	}
	return b
}

// Rows は、グリッドの行定義を設定します。
func (b *AdvancedGridBuilder) Rows(defs ...layout.TrackDefinition) *AdvancedGridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.AdvancedGridLayout); ok {
		gridLayout.RowDefinitions = defs
		b.Widget.MarkDirty(true)
	}
	return b
}

// HorizontalGap は、グリッドの水平方向の間隔を設定します。
func (b *AdvancedGridBuilder) HorizontalGap(gap int) *AdvancedGridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.AdvancedGridLayout); ok {
		gridLayout.HorizontalGap = gap
		b.Widget.MarkDirty(true)
	}
	return b
}

// VerticalGap は、グリッドの垂直方向の間隔を設定します。
func (b *AdvancedGridBuilder) VerticalGap(gap int) *AdvancedGridBuilder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.AdvancedGridLayout); ok {
		gridLayout.VerticalGap = gap
		b.Widget.MarkDirty(true)
	}
	return b
}

// Gap は、水平・垂直両方の間隔を同じ値に設定します。
func (b *AdvancedGridBuilder) Gap(gap int) *AdvancedGridBuilder {
	b.HorizontalGap(gap)
	b.VerticalGap(gap)
	return b
}

// Add は、指定された行と列にウィジェットを追加します。
// このメソッドは、ウィジェットビルダーを直接受け取るジェネリックなヘルパーです。
func Add[W component.Widget, WB interface {
	component.BuilderFinalizer[W]
	component.ErrorAdder
}](b *AdvancedGridBuilder, row, col, rowSpan, colSpan int, builder WB) {
	widget, err := builder.Build()
	if err != nil {
		b.AddError(err)
		// エラーがあっても不完全なウィジェットを追加することで、レイアウトの崩れを確認しやすくします
	}

	// 【修正点】
	// ジェネリック型Wの変数は直接nilと比較できないため、reflectを使用してnilチェックを行います。
	// widgetはインターフェース(component.Widget)を満たすポインタ型(*widget.Buttonなど)であるため、
	// この方法で安全にチェックできます。
	v := reflect.ValueOf(widget)
	isNil := !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil())

	if !isNil {
		placement := layout.GridPlacementData{
			Row: row, Col: col, RowSpan: rowSpan, ColSpan: colSpan,
		}
		widget.SetLayoutData(placement)
		b.AddChild(widget)
	}
}

// ButtonAt は、指定された位置とスパンでButtonウィジェットをグリッドに追加します。
func (b *AdvancedGridBuilder) ButtonAt(row, col, rowSpan, colSpan int, buildFunc func(*widget.ButtonBuilder)) *AdvancedGridBuilder {
	builder := widget.NewButtonBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	Add(b, row, col, rowSpan, colSpan, builder)
	return b
}

// LabelAt は、指定された位置とスパンでLabelウィジェットをグリッドに追加します。
func (b *AdvancedGridBuilder) LabelAt(row, col, rowSpan, colSpan int, buildFunc func(*widget.LabelBuilder)) *AdvancedGridBuilder {
	builder := widget.NewLabelBuilder()
	if buildFunc != nil {
		buildFunc(builder)
	}
	Add(b, row, col, rowSpan, colSpan, builder)
	return b
}

// Build はコンテナの構築を完了します。
func (b *AdvancedGridBuilder) Build() (*container.Container, error) { return b.Builder.Build() }