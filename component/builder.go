package component

import (
	"errors"
	"fmt"
	"furoshiki/event"
	"furoshiki/style"
	"image/color"
	"reflect"
)

// パッケージ全体で利用できるよう、共通エラーをエクスポートします。
var (
	ErrNilChild             = errors.New("child cannot be nil")
	ErrWidgetNotInitialized = errors.New("widget not properly initialized")
	ErrInvalidSize          = errors.New("size must be non-negative")
	ErrInvalidFlex          = errors.New("flex must be non-negative")
	ErrInvalidBorderWidth   = errors.New("border width must be non-negative")
)

// Buildableは、汎用のcomponent.Builderが操作するウィジェットが満たすべき能力を定義するインターフェースです。
// これにより、ビルダーは型安全にウィジェットのプロパティを操作できます。
type Buildable interface {
	Widget
	SizeSetter
	MinSizeSetter
	StyleGetterSetter
	SetFlex(flex int)
	GetFlex() int
	SetLayoutBoundary(isBoundary bool)
	SetLayoutData(data any)
	GetLayoutData() any
	EventProcessor
	AbsolutePositioner
}

// Builder は、すべてのウィジェットビルダーの汎用基底クラスです。
// サイズ、スタイル、レイアウトプロパティを設定するための共通メソッドを提供します。
// T は具体的なビルダー型（例: *LabelBuilder）です。
// W は構築中のウィジェット型（例: *Label）で、Buildableインターフェースを満たす必要があります。
type Builder[T any, W Buildable] struct {
	Widget W
	errors []error
	Self   T
}

// Init は基底ビルダーを初期化します。具象ビルダーのコンストラクタから呼び出す必要があります。
func (b *Builder[T, W]) Init(self T, widget W) {
	b.Self = self
	b.Widget = widget
}

// AddError はビルドエラーをビルダーに追加します。
func (b *Builder[T, W]) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// AbsolutePosition は、親コンテナ内でのウィジェットの希望相対位置を設定します。
// これは、親コンテナがAbsoluteLayout（例: ZStack）を使用している場合にのみ有効です。
func (b *Builder[T, W]) AbsolutePosition(x, y int) T {
	b.Widget.SetRequestedPosition(x, y)
	return b.Self
}

// Size はウィジェットのサイズを設定します。
func (b *Builder[T, W]) Size(width, height int) T {
	if err := validateSize(width, height); err != nil {
		b.AddError(err)
	} else {
		b.Widget.SetSize(width, height)
	}
	return b.Self
}

// MinSize はウィジェットの最小サイズを設定します。
func (b *Builder[T, W]) MinSize(width, height int) T {
	if err := validateSize(width, height); err != nil {
		b.AddError(err)
	} else {
		b.Widget.SetMinSize(width, height)
	}
	return b.Self
}

// validateSize はサイズが有効かどうかを検証します
func validateSize(width, height int) error {
	if width < 0 || height < 0 {
		return fmt.Errorf("%w, got %dx%d", ErrInvalidSize, width, height)
	}
	return nil
}

// Style は指定されたスタイルをウィジェットの既存のスタイルとマージします。
func (b *Builder[T, W]) Style(s style.Style) T {
	existingStyle := b.Widget.GetStyle()
	b.Widget.SetStyle(style.Merge(existingStyle, s))
	return b.Self
}

// Flex はFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (b *Builder[T, W]) Flex(flex int) T {
	if flex < 0 {
		b.AddError(fmt.Errorf("%w, got %d", ErrInvalidFlex, flex))
	} else {
		b.Widget.SetFlex(flex)
	}
	return b.Self
}

// --- スタイルヘルパーメソッド ---

// applyStyleProperty はスタイルプロパティを設定するための共通ヘルパー関数です
func (b *Builder[T, W]) applyStyleProperty(setter func(style.Style) style.Style) T {
	newStyle := setter(b.Widget.GetStyle())
	b.Widget.SetStyle(newStyle)
	return b.Self
}

// ApplyStyles はスタイルオプション関数を用いてウィジェットのスタイルを柔軟に設定します。
func (b *Builder[T, W]) ApplyStyles(opts ...style.StyleOption) T {
	newStyle := b.Widget.GetStyle()
	for _, opt := range opts {
		opt(&newStyle)
	}
	b.Widget.SetStyle(newStyle)
	return b.Self
}

// BackgroundColor はウィジェットの背景色を設定します。
func (b *Builder[T, W]) BackgroundColor(c color.Color) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.Background = style.PColor(c)
		return s
	})
}

// TextColor はウィジェットのテキスト色を設定します。
func (b *Builder[T, W]) TextColor(c color.Color) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.TextColor = style.PColor(c)
		return s
	})
}

// Margin はすべての辺に同じマージン値を設定します。
func (b *Builder[T, W]) Margin(m int) T {
	return b.MarginInsets(style.Insets{Top: m, Right: m, Bottom: m, Left: m})
}

// MarginInsets は各辺に個別のマージン値を設定します。
func (b *Builder[T, W]) MarginInsets(i style.Insets) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.Margin = style.PInsets(i)
		return s
	})
}

// Padding はすべての辺に同じパディング値を設定します。
func (b *Builder[T, W]) Padding(p int) T {
	return b.PaddingInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})
}

// PaddingInsets は各辺に個別のパディング値を設定します。
func (b *Builder[T, W]) PaddingInsets(i style.Insets) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.Padding = style.PInsets(i)
		return s
	})
}

// BorderRadius はウィジェットの角の半径を設定します。
func (b *Builder[T, W]) BorderRadius(radius float32) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.BorderRadius = style.PFloat32(radius)
		return s
	})
}

// Border はウィジェットの境界線の幅と色を設定します。
func (b *Builder[T, W]) Border(width float32, c color.Color) T {
	if width < 0 {
		b.AddError(fmt.Errorf("%w, got %f", ErrInvalidBorderWidth, width))
		return b.Self
	}
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.BorderWidth = style.PFloat32(width)
		s.BorderColor = style.PColor(c)
		return s
	})
}

// TextAlign はテキストの水平方向の揃え位置を設定します。
func (b *Builder[T, W]) TextAlign(align style.TextAlignType) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.TextAlign = style.PTextAlignType(align)
		return s
	})
}

// VerticalAlign はテキストの垂直方向の揃え位置を設定します。
func (b *Builder[T, W]) VerticalAlign(align style.VerticalAlignType) T {
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.VerticalAlign = style.PVerticalAlignType(align)
		return s
	})
}

// --- 汎用イベントハンドラ設定メソッド ---
// AddOnClick は、ウィジェットがクリックされたときに実行されるイベントハンドラを追加します。
func (b *Builder[T, W]) AddOnClick(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.EventClick, handler)
	return b.Self
}

// AddOnMouseEnter は、マウスカーソルがウィジェット上に入ったときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseEnter(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseEnter, handler)
	return b.Self
}

// AddOnMouseLeave は、マウスカーソルがウィジェットから離れたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseLeave(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseLeave, handler)
	return b.Self
}

// AddOnMouseMove は、マウスカーソルがウィジェット上で移動したときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseMove(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseMove, handler)
	return b.Self
}

// AddOnMouseDown は、マウスボタンがウィジェット上で押されたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseDown(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseDown, handler)
	return b.Self
}

// AddOnMouseUp は、マウスボタンがウィジェット上で解放されたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseUp(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseUp, handler)
	return b.Self
}

// AddOnMouseScroll は、マウスホイールがウィジェット上でスクロールされたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseScroll(handler event.EventHandler) T {
	b.Widget.AddEventHandler(event.MouseScroll, handler)
	return b.Self
}

// AssignTo は、ビルド中のウィジェットインスタンスへのポインタを変数に代入します。
func (b *Builder[T, W]) AssignTo(target any) T {
	if target == nil {
		b.AddError(errors.New("AssignTo target cannot be nil"))
		return b.Self
	}

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		b.AddError(fmt.Errorf("AssignTo target must be a non-nil pointer, got %T", target))
		return b.Self
	}

	targetElem := targetVal.Elem()
	widgetVal := reflect.ValueOf(b.Widget)
	if !widgetVal.Type().AssignableTo(targetElem.Type()) {
		b.AddError(fmt.Errorf("cannot assign widget of type %s to target of type %s", widgetVal.Type(), targetElem.Type()))
		return b.Self
	}

	if !targetElem.CanSet() {
		b.AddError(fmt.Errorf("AssignTo target cannot be set"))
		return b.Self
	}
	targetElem.Set(widgetVal)

	return b.Self
}

// Build はウィジェットの構築を完了します。
func (b *Builder[T, W]) Build() (W, error) {
	if len(b.errors) > 0 {
		var zero W
		// エラー時にどのウィジェットで問題が起きたか分かりやすいように型名を含める
		typeName := reflect.TypeOf(b.Widget)
		if typeName.Kind() == reflect.Ptr {
			typeName = typeName.Elem()
		}
		return zero, fmt.Errorf("%s build errors: %w", typeName.Name(), errors.Join(b.errors...))
	}
	b.Widget.MarkDirty(true)
	return b.Widget, nil
}