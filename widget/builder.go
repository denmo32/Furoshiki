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
	// [変更] CalculateMinSizeはGetMinSizeに統合されたため、ここでは不要になります。
	// GetMinSizeはcomponent.Widgetインターフェースに含まれています。
	// SetRequestedPositionはcomponent.LayoutableWidgetに実装されているため、
	// それを埋め込むことでこのインターフェースを満たします。
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
// [重要] この設定は、親コンテナが `AbsoluteLayout` (例: `ui.ZStack` ヘルパーで作成) を
// 使用している場合にのみ有効です。`FlexLayout` (VStack, HStack) など他のレイアウトでは無視されます。
func (b *Builder[T, W]) Position(x, y int) T {
	b.Widget.SetRequestedPosition(x, y)
	return b.self
}

// Size はウィジェットのサイズを設定します。
// FlexLayout内では、この値はflex値が設定されていない場合の「基本サイズ」として扱われ、
// レイアウト計算によって上書きされる可能性があります。
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

// [改良] Build はウィジェットの構築を完了します。
// 以前はここでサイズ計算を行っていましたが、レイアウトシステムとの責務を明確に分離するため、
// このメソッドはエラーチェックのみを行い、ウィジェットインスタンスを返すように変更されました。
// 最終的なサイズと位置は、UIツリーに追加された後、親コンテナのレイアウトシステムによって決定されます。
func (b *Builder[T, W]) Build(typeName string) (W, error) {
	if len(b.errors) > 0 {
		var zero W
		return zero, fmt.Errorf("%s build errors: %w", typeName, errors.Join(b.errors...))
	}

	// ビルダーの責務はウィジェットのプロパティを設定することまで。
	// 最終的なサイズ計算と適用は、親コンテナのレイアウトシステムの役割です。

	// ウィジェットがダーティマークされていることを保証し、最初のフレームでレイアウトが実行されるようにします。
	b.Widget.MarkDirty(true)

	return b.Widget, nil
}

// max は2つの整数のうち大きい方を返します。
// (この関数はcomponent.text_widgetでも使用するため、将来的に共通パッケージに移動することを検討しても良いでしょう)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}