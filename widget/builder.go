package widget

import (
	"errors"
	"fmt"
	"furoshiki/component"
	"furoshiki/style"
	"image/color"
)

// textWidget は、component.TextWidget を埋め込むウィジェットが満たすインターフェースです。
// これにより、ジェネリックビルダーがテキスト関連ウィジェットの共通メソッドを呼び出せるようになります。
type textWidget interface {
	component.Widget
	SetText(string)
	SetRequestedPosition(x, y int)
}

// Builder は、component.TextWidget をベースにしたウィジェットビルダーのための汎用的なベースです。
// ジェネリクスを使用することで、コードの重複を避けつつ、型安全なメソッドチェーンを実現します。
// T は具象ビルダーの型 (例: *LabelBuilder)
// W はビルドされるウィジェットの型 (例: *Label)
type Builder[T any, W textWidget] struct {
	Widget W // 具象ビルダーからアクセスできるように公開します
	errors []error
	self   T // メソッドチェーンを可能にするために、具象ビルダー自身への参照を保持します
}

// Init はベースビルダーを初期化します。具象ビルダーのコンストラクタから呼び出す必要があります。
func (b *Builder[T, W]) Init(self T, widget W) {
	b.self = self
	b.Widget = widget
}

// Text はウィジェットのテキストを設定します。
func (b *Builder[T, W]) Text(text string) T {
	b.Widget.SetText(text)
	return b.self
}

// Positionは、ウィジェットの希望する相対位置を設定します。
//
// 重要: このメソッドは、親コンテナが `AbsoluteLayout` (主に `ui.ZStack` で作成) を
// 使用している場合にのみ有効です。`FlexLayout` (`VStack`, `HStack`) や `GridLayout` の
// 中にあるウィジェットに対してこのメソッドを使用しても、設定はレイアウトシステムによって
// 無視されるため効果はありません。
func (b *Builder[T, W]) Position(x, y int) T {
	b.Widget.SetRequestedPosition(x, y)
	return b.self
}

// Size はウィジェットのサイズを設定します。
func (b *Builder[T, W]) Size(width, height int) T {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("size must be non-negative, got %dx%d", width, height))
	} else {
		b.Widget.SetSize(width, height)
	}
	return b.self
}

// MinSize はウィジェットの最小サイズを設定します。
func (b *Builder[T, W]) MinSize(width, height int) T {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("min size must be non-negative, got %dx%d", width, height))
	} else {
		b.Widget.SetMinSize(width, height)
	}
	return b.self
}

// Style はウィジェットの基本スタイルを設定します。既存のスタイルとマージされます。
func (b *Builder[T, W]) Style(s style.Style) T {
	existingStyle := b.Widget.GetStyle()
	b.Widget.SetStyle(style.Merge(existingStyle, s))
	return b.self
}

// Flex は、FlexLayout 内でウィジェットがどのように伸縮するかを設定します。
func (b *Builder[T, W]) Flex(flex int) T {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("flex must be non-negative, got %d", flex))
	} else {
		b.Widget.SetFlex(flex)
	}
	return b.self
}

// --- Style Helpers ---

// BackgroundColor はウィジェットの背景色を設定します。
func (b *Builder[T, W]) BackgroundColor(c color.Color) T {
	return b.Style(style.Style{Background: style.PColor(c)})
}

// TextColor はウィジェットのテキスト色を設定します。
func (b *Builder[T, W]) TextColor(c color.Color) T {
	return b.Style(style.Style{TextColor: style.PColor(c)})
}

// Margin はウィジェットのマージンを四方同じ値で設定します。
func (b *Builder[T, W]) Margin(m int) T {
	return b.Style(style.Style{Margin: style.PInsets(style.Insets{Top: m, Right: m, Bottom: m, Left: m})})
}

// MarginInsets はウィジェットのマージンを各辺個別に設定します。
func (b *Builder[T, W]) MarginInsets(i style.Insets) T {
	return b.Style(style.Style{Margin: style.PInsets(i)})
}

// Padding はウィジェットのパディングを四方同じ値で設定します。
func (b *Builder[T, W]) Padding(p int) T {
	return b.Style(style.Style{Padding: style.PInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})})
}

// PaddingInsets はウィジェットのパディングを各辺個別に設定します。
func (b *Builder[T, W]) PaddingInsets(i style.Insets) T {
	return b.Style(style.Style{Padding: style.PInsets(i)})
}

// BorderRadius はウィジェットの角丸の半径を設定します。
func (b *Builder[T, W]) BorderRadius(radius float32) T {
	return b.Style(style.Style{BorderRadius: style.PFloat32(radius)})
}

// Border はウィジェットの境界線を設定します。
func (b *Builder[T, W]) Border(width float32, c color.Color) T {
	if width < 0 {
		b.errors = append(b.errors, fmt.Errorf("border width must be non-negative, got %f", width))
		return b.self
	}
	return b.Style(style.Style{
		BorderWidth: style.PFloat32(width),
		BorderColor: style.PColor(c),
	})
}

// TextAlign はテキストの水平方向の揃えを設定します。
func (b *Builder[T, W]) TextAlign(align style.TextAlignType) T {
	return b.Style(style.Style{TextAlign: style.PTextAlignType(align)})
}

// VerticalAlign はテキストの垂直方向の揃えを設定します。
func (b *Builder[T, W]) VerticalAlign(align style.VerticalAlignType) T {
	return b.Style(style.Style{VerticalAlign: style.PVerticalAlignType(align)})
}

// Build はウィジェットの構築を完了します。
func (b *Builder[T, W]) Build(typeName string) (W, error) {
	if len(b.errors) > 0 {
		var zero W
		return zero, fmt.Errorf("%s build errors: %w", typeName, errors.Join(b.errors...))
	}
	// ウィジェットがダーティマークされていることを保証し、最初のフレームでレイアウトが実行されるようにします。
	b.Widget.MarkDirty(true)
	return b.Widget, nil
}