package ui

import (
	"furoshiki/container"
	"furoshiki/layout"
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
