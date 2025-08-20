package layout

import (
	"furoshiki/component"
)

// Layout は、コンテナ内の子要素をどのように配置するかを決定するロジックのインターフェースです。
type Layout interface {
	// Layout は、指定されたコンテナの子要素のサイズと位置を計算し、設定します。
	Layout(container Container)
}

// Insets はパディングやマージンの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Container は、レイアウトが必要とするコンテナの振る舞いを定義するインターフェースです。
// これにより、layoutパッケージは具体的なコンテナの実装から独立しています。
type Container interface {
	GetSize() (width, height int)
	GetPosition() (x, y int)
	GetChildren() []component.Widget
	GetPadding() Insets
}

// LayoutType はレイアウトの種類を定義します。現在は直接使用されていません。
type LayoutType int

const (
	LayoutAbsolute LayoutType = iota
	LayoutFlex
)

// Alignment は要素の揃え位置を定義します。
type Alignment int

const (
	AlignStart   Alignment = iota // 要素を開始位置に揃えます (左揃え or 上揃え)
	AlignCenter                 // 要素を中央に揃えます
	AlignEnd                    // 要素を終了位置に揃えます (右揃え or 下揃え)
	AlignStretch                // 要素をコンテナのサイズいっぱいに引き伸ばします (交差軸でのみ有効)
)

// Direction は要素を並べる方向を定義します。
type Direction int

const (
	DirectionRow    Direction = iota // 水平方向 (左から右へ)
	DirectionColumn                  // 垂直方向 (上から下へ)
)