package widget

import (
	"errors"
	"fmt"
	"furoshiki/component"
	"furoshiki/style"
)

// textWidget は、component.TextWidget を埋め込むウィジェットが満たすインターフェースです。
// これにより、ジェネリックビルダーがテキスト関連ウィジェットの共通メソッドを呼び出せるようになります。
type textWidget interface {
	component.Widget
	SetText(string)
	CalculateMinSize() (int, int)
	// SetRequestedPositionはcomponent.LayoutableWidgetに実装されているため、
	// それを埋め込むことでこのインターフェースを満たします。
	SetRequestedPosition(x, y int)
}

// Builder は、component.TextWidget をベースにしたウィジェットビルダーのための汎用的なベースです。
// T は具象ビルダーの型 (例: *LabelBuilder)
// W はビルドされるウィジェットの型 (例: *Label)
type Builder[T any, W textWidget] struct {
	Widget W // 具象ビルダーからアクセスできるように公開します
	errors []error
	self   T
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

// [追加] Positionは、ウィジェットの希望する相対位置を設定します。
// この値は、親コンテナがAbsoluteLayoutを使用している場合に、子の配置位置として利用されます。
// FlexLayoutなど他のレイアウトでは無視されることがあります。
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

// Style はウィジェットの基本スタイルを設定します。
func (b *Builder[T, W]) Style(s style.Style) T {
	// [改善] GetStyle()が値型を返すようになったため、ポインタアクセス(*)が不要になります。
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

// Build はウィジェットの構築を完了します。
// エラーをチェックし、最小サイズを計算して設定します。
func (b *Builder[T, W]) Build(typeName string) (W, error) {
	if len(b.errors) > 0 {
		var zero W
		return zero, fmt.Errorf("%s build errors: %w", typeName, errors.Join(b.errors...))
	}
	minWidth, minHeight := b.Widget.CalculateMinSize()
	b.Widget.SetMinSize(minWidth, minHeight)
	return b.Widget, nil
}