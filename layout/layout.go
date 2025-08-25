package layout

import (
	"furoshiki/component"
)

// Layout は、コンテナ内の子要素をどのように配置するかを決定するロジックのインターフェースです。
type Layout interface {
	Layout(container Container)
}

// Insets はパディングやマージンの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Container は、レイアウトが必要とするコンテナの振る舞いを定義するインターフェースです。
type Container interface {
	GetSize() (width, height int)
	GetPosition() (x, y int)
	GetChildren() []component.Widget
	GetPadding() Insets
}

// ScrollViewer は、ScrollViewLayoutがScrollViewウィジェットを操作するために必要なメソッドを定義します。
type ScrollViewer interface {
	Container
	GetContentContainer() component.Widget
	GetVScrollBar() component.ScrollBarWidget
	GetScrollY() float64
	SetScrollY(y float64)
	SetContentHeight(h int)
}

// Alignment は要素の揃え位置を定義します。
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// Direction は要素を並べる方向を定義します。
type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)
