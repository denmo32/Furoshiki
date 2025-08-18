package layout

import (
	"furoshiki/core"
	"furoshiki/style"
)

// Layout は、コンテナ内の子要素をどのように配置するかを決定するロジックのインターフェースです。
type Layout interface {
	Layout(container Container)
}

// Container は、レイアウトが必要とするコンテナの振る舞いを定義するインターフェースです。
// これにより、layoutパッケージは具体的なコンテナの実装から独立しています。
type Container interface {
	GetSize() (width, height int)
	GetPosition() (x, y int)
	GetChildren() []core.Widget
	GetStyle() *style.Style
}

// LayoutType はレイアウトの種類を定義します。
type LayoutType int

const (
	LayoutAbsolute LayoutType = iota
	LayoutFlex
)

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