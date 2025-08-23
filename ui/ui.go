package ui

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/layout"
	"furoshiki/widget"
)

// Builder は、UI構造を宣言的に構築するための統一されたインターフェースです。
// component.Builderを直接埋め込むことで、Size(), Padding(), Flex()などの共通メソッドを
// 再定義することなく利用でき、冗長なラッパーコードを排除しています。
// ジェネリクスの型引数により、継承されたメソッドは正しく*Builderを返すようになります。
type Builder struct {
	component.Builder[*Builder, *container.Container]
}

// buildContainerAndBuilder は、指定されたレイアウトでコンテナとビルダーを初期化する内部ヘルパーです。
func buildContainerAndBuilder(l layout.Layout, buildFunc func(*Builder)) *Builder {
	// 1. UIの基礎となるContainerウィジェットを作成します。
	c := &container.Container{}
	// Container自身をselfとして渡して、基本的なウィジェット機能を初期化します。
	c.LayoutableWidget = component.NewLayoutableWidget(c)
	// 指定されたレイアウトを設定します。
	c.SetLayout(l)

	// 2. 新しいui.Builderインスタンスを作成します。
	b := &Builder{}
	// component.BuilderのInitメソッドを呼び出し、
	// 自分自身(*Builder)と構築対象のウィジェット(*container.Container)を関連付けます。
	b.Init(b, c)

	// 3. ユーザー提供の構築関数を実行し、コンテナに子要素などを追加します。
	if buildFunc != nil {
		buildFunc(b)
	}
	return b
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを作成します。
func VStack(buildFunc func(*Builder)) *Builder {
	flexLayout := &layout.FlexLayout{
		Direction: layout.DirectionColumn,
	}
	return buildContainerAndBuilder(flexLayout, buildFunc)
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを作成します。
func HStack(buildFunc func(*Builder)) *Builder {
	flexLayout := &layout.FlexLayout{
		Direction: layout.DirectionRow,
	}
	return buildContainerAndBuilder(flexLayout, buildFunc)
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを作成します。
func ZStack(buildFunc func(*Builder)) *Builder {
	return buildContainerAndBuilder(&layout.AbsoluteLayout{}, buildFunc)
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを作成します。
func Grid(buildFunc func(*Builder)) *Builder {
	// デフォルトで1列のグリッドレイアウトを設定します。
	return buildContainerAndBuilder(&layout.GridLayout{Columns: 1}, buildFunc)
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
// このメソッドは、ui.Builderがコンテナを操作する能力を持つために必要です。
func (b *Builder) AddChild(child component.Widget) *Builder {
	if child != nil {
		b.Widget.AddChild(child)
	} else {
		b.AddError(fmt.Errorf("cannot add a nil child widget"))
	}
	return b
}

// AddChildren はコンテナに複数の子ウィジェットを追加します。
func (b *Builder) AddChildren(children ...component.Widget) *Builder {
	for _, child := range children {
		b.AddChild(child) // AddChildを経由することでnilチェックを一元化
	}
	return b
}

// --- Nested Container Adders ---

// 【改善】ネストされたコンテナを追加するロジックを共通ヘルパーメソッド `addNestedContainer` に集約し、
// コードの重複を排除しました。これにより、各メソッドの意図がより明確になり、
// メンテナンス性が向上します。

// addNestedContainer は、ネストされたコンテナビルダーを受け取り、
// それをビルドして現在構築中のコンテナの子として追加する共通ロジックです。
func (b *Builder) addNestedContainer(nestedBuilder *Builder) *Builder {
	// ネストされたビルダーを使ってコンテナウィジェットを構築します。
	container, err := nestedBuilder.Build()
	// エラーがあれば現在のビルダーに記録します。
	b.AddError(err)
	// 構築したコンテナを、現在構築中のコンテナの子として追加します。
	b.AddChild(container)
	return b
}

// HStack は水平方向に子要素を配置するFlexLayoutコンテナを「子として」追加します。
func (b *Builder) HStack(buildFunc func(*Builder)) *Builder {
	// トップレベルのHStackビルダー関数を呼び出し、その結果を共通ヘルパーに渡して
	// 現在のコンテナに子として追加します。
	return b.addNestedContainer(HStack(buildFunc))
}

// VStack は垂直方向に子要素を配置するFlexLayoutコンテナを「子として」追加します。
func (b *Builder) VStack(buildFunc func(*Builder)) *Builder {
	// トップレベルのVStackビルダー関数を呼び出し、その結果を共通ヘルパーに渡します。
	return b.addNestedContainer(VStack(buildFunc))
}

// ZStack は子要素を重ねて配置するAbsoluteLayoutコンテナを「子として」追加します。
func (b *Builder) ZStack(buildFunc func(*Builder)) *Builder {
	// トップレベルのZStackビルダー関数を呼び出し、その結果を共通ヘルパーに渡します。
	return b.addNestedContainer(ZStack(buildFunc))
}

// Grid は子要素を格子状に配置するGridLayoutコンテナを「子として」追加します。
func (b *Builder) Grid(buildFunc func(*Builder)) *Builder {
	// トップレベルのGridビルダー関数を呼び出し、その結果を共通ヘルパーに渡します。
	return b.addNestedContainer(Grid(buildFunc))
}

// --- Container-specific Property Wrappers ---

// RelayoutBoundary は、このコンテナをレイアウト計算の境界として設定します。
// component.Builderには存在しないコンテナ固有の機能であるため、ui.Builderで実装します。
func (b *Builder) RelayoutBoundary(isBoundary bool) *Builder {
	b.Widget.SetRelayoutBoundary(isBoundary)
	return b
}

// [新規追加]
// ClipChildren は、コンテナがその境界外に子要素を描画しないように設定します（クリッピング）。
// component.Builderには存在しないコンテナ固有の機能であるため、ui.Builderで実装します。
func (b *Builder) ClipChildren(clips bool) *Builder {
	// b.Widget は *container.Container 型なので、SetClipsChildren を呼び出せます。
	// container.Container側には対応するメソッドが実装済みです。
	b.Widget.SetClipsChildren(clips)
	return b
}


// --- FlexLayout Specific Methods ---
// これらはui.Builderがラップするコンテナのレイアウトプロパティを操作する独自の機能です。

// Gap はFlexLayoutの子要素間の間隔を設定します。
// VStackまたはHStack内でのみ有効です。
func (b *Builder) Gap(gap int) *Builder {
	// b.Widgetは*container.Container型なので、GetLayout()を呼び出せます。
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
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
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
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
	if flexLayout, ok := b.Widget.GetLayout().(*layout.FlexLayout); ok {
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
// これらもui.Builderが提供する独自の機能です。

// Columns はグリッドの列数を設定します。
// Grid内でのみ有効です。
func (b *Builder) Columns(count int) *Builder {
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
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
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
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
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
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
	if gridLayout, ok := b.Widget.GetLayout().(*layout.GridLayout); ok {
		if gridLayout.VerticalGap != gap {
			gridLayout.VerticalGap = gap
			b.Widget.MarkDirty(true)
		}
	} else {
		b.AddError(fmt.Errorf("VerticalGap() is only applicable to GridLayout"))
	}
	return b
}

// Build はコンテナの構築を完了します。
// component.BuilderにBuildメソッドがあるため、このメソッドは厳密には不要ですが、
// 返り値の型を*container.Containerに明示するために定義しています。
func (b *Builder) Build() (*container.Container, error) {
	return b.Builder.Build()
}