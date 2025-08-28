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
	// UPDATE: 型アサーション失敗時に使用するエラーを追加
	ErrUnsupportedOperation = errors.New("unsupported operation for this widget type")
)

// Builder は、すべてのウィジェットビルダーの汎用基底クラスです。
// サイズ、スタイル、レイアウトプロパティを設定するための共通メソッドを提供します。
// T は具体的なビルダー型（例: *LabelBuilder）です。
// W は構築中のウィジェット型（例: *Label）です。
// UPDATE: Buildableインターフェース制約を削除
type Builder[T any, W Widget] struct {
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

// UPDATE: 型アサーションに失敗した場合のエラーを生成するヘルパー関数を追加
func (b *Builder[T, W]) newUnsupportedError(operation string) error {
	// any(b.Widget) を使うことで、インターフェース型のnilポインタでも安全に型名を取得
	typeName := "Unknown"
	if t := reflect.TypeOf(any(b.Widget)); t != nil {
		if t.Kind() == reflect.Ptr {
			typeName = t.Elem().Name()
		} else {
			typeName = t.Name()
		}
	}
	return fmt.Errorf("%w: %s on %s", ErrUnsupportedOperation, operation, typeName)
}

// AbsolutePosition は、親コンテナ内でのウィジェットの希望相対位置を設定します。
// これは、親コンテナがAbsoluteLayout（例: ZStack）を使用している場合にのみ有効です。
func (b *Builder[T, W]) AbsolutePosition(x, y int) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if ap, ok := any(b.Widget).(AbsolutePositioner); ok {
		ap.SetRequestedPosition(x, y)
	} else {
		b.AddError(b.newUnsupportedError("AbsolutePosition"))
	}
	return b.Self
}

// Size はウィジェットのサイズを設定します。
func (b *Builder[T, W]) Size(width, height int) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if ss, ok := any(b.Widget).(SizeSetter); ok {
		if err := validateSize(width, height); err != nil {
			b.AddError(err)
		} else {
			ss.SetSize(width, height)
		}
	} else {
		b.AddError(b.newUnsupportedError("Size"))
	}
	return b.Self
}

// MinSize はウィジェットの最小サイズを設定します。
func (b *Builder[T, W]) MinSize(width, height int) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if mss, ok := any(b.Widget).(MinSizeSetter); ok {
		if err := validateSize(width, height); err != nil {
			b.AddError(err)
		} else {
			mss.SetMinSize(width, height)
		}
	} else {
		b.AddError(b.newUnsupportedError("MinSize"))
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
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if sgs, ok := any(b.Widget).(StyleGetterSetter); ok {
		existingStyle := sgs.GetStyle()
		sgs.SetStyle(style.Merge(existingStyle, s))
	} else {
		b.AddError(b.newUnsupportedError("Style"))
	}
	return b.Self
}

// Flex はFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (b *Builder[T, W]) Flex(flex int) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	// LayoutPropertiesOwnerは全てのウィジェットが実装すべき基本機能の一つとみなし、
	// LayoutPropertiesコンポーネント経由でアクセスします。
	if lpo, ok := any(b.Widget).(LayoutPropertiesOwner); ok {
		if flex < 0 {
			b.AddError(fmt.Errorf("%w, got %d", ErrInvalidFlex, flex))
		} else {
			lpo.GetLayoutProperties().SetFlex(flex)
		}
	} else {
		b.AddError(b.newUnsupportedError("Flex"))
	}
	return b.Self
}

// --- スタイルヘルパーメソッド ---

// applyStyleProperty はスタイルプロパティを設定するための共通ヘルパー関数です
func (b *Builder[T, W]) applyStyleProperty(setter func(style.Style) style.Style) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if sgs, ok := any(b.Widget).(StyleGetterSetter); ok {
		newStyle := setter(sgs.GetStyle())
		sgs.SetStyle(newStyle)
	} else {
		// エラーは各呼び出し元で追加されるため、ここでは何もしません。
	}
	return b.Self
}

// ApplyStyles はスタイルオプション関数を用いてウィジェットのスタイルを柔軟に設定します。
func (b *Builder[T, W]) ApplyStyles(opts ...style.StyleOption) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if sgs, ok := any(b.Widget).(StyleGetterSetter); ok {
		newStyle := sgs.GetStyle()
		for _, opt := range opts {
			opt(&newStyle)
		}
		sgs.SetStyle(newStyle)
	} else {
		b.AddError(b.newUnsupportedError("ApplyStyles"))
	}
	return b.Self
}

// BackgroundColor はウィジェットの背景色を設定します。
func (b *Builder[T, W]) BackgroundColor(c color.Color) T {
	// UPDATE: 型アサーションのチェックはapplyStyleProperty内で行われるが、
	// 失敗時のエラーメッセージを明確にするためにここでもチェックします。
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("BackgroundColor"))
		return b.Self
	}
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.Background = style.PColor(c)
		return s
	})
}

// TextColor はウィジェットのテキスト色を設定します。
func (b *Builder[T, W]) TextColor(c color.Color) T {
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("TextColor"))
		return b.Self
	}
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
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("MarginInsets"))
		return b.Self
	}
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
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("PaddingInsets"))
		return b.Self
	}
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.Padding = style.PInsets(i)
		return s
	})
}

// BorderRadius はウィジェットの角の半径を設定します。
func (b *Builder[T, W]) BorderRadius(radius float32) T {
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("BorderRadius"))
		return b.Self
	}
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
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("Border"))
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
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("TextAlign"))
		return b.Self
	}
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.TextAlign = style.PTextAlignType(align)
		return s
	})
}

// VerticalAlign はテキストの垂直方向の揃え位置を設定します。
func (b *Builder[T, W]) VerticalAlign(align style.VerticalAlignType) T {
	if _, ok := any(b.Widget).(StyleGetterSetter); !ok {
		b.AddError(b.newUnsupportedError("VerticalAlign"))
		return b.Self
	}
	return b.applyStyleProperty(func(s style.Style) style.Style {
		s.VerticalAlign = style.PVerticalAlignType(align)
		return s
	})
}

// --- 汎用イベントハンドラ設定メソッド ---
// addEventHandler はイベントハンドラを追加するための共通ヘルパーです。
func (b *Builder[T, W]) addEventHandler(eventType event.EventType, handler event.EventHandler) T {
	// UPDATE: 型アサーションでインターフェースサポートをチェック
	if ep, ok := any(b.Widget).(EventProcessor); ok {
		ep.AddEventHandler(eventType, handler)
	} else {
		b.AddError(b.newUnsupportedError(fmt.Sprintf("AddOn%s", eventType)))
	}
	return b.Self
}

// AddOnClick は、ウィジェットがクリックされたときに実行されるイベントハンドラを追加します。
func (b *Builder[T, W]) AddOnClick(handler event.EventHandler) T {
	return b.addEventHandler(event.EventClick, handler)
}

// AddOnMouseEnter は、マウスカーソルがウィジェット上に入ったときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseEnter(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseEnter, handler)
}

// AddOnMouseLeave は、マウスカーソルがウィジェットから離れたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseLeave(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseLeave, handler)
}

// AddOnMouseMove は、マウスカーソルがウィジェット上で移動したときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseMove(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseMove, handler)
}

// AddOnMouseDown は、マウスボタンがウィジェット上で押されたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseDown(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseDown, handler)
}

// AddOnMouseUp は、マウスボタンがウィジェット上で解放されたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseUp(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseUp, handler)
}

// AddOnMouseScroll は、マウスホイールがウィジェット上でスクロールされたときに実行されるハンドラを追加します。
func (b *Builder[T, W]) AddOnMouseScroll(handler event.EventHandler) T {
	return b.addEventHandler(event.MouseScroll, handler)
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
	// any()でラップすることで、インターフェース型のnilポインタを正しく扱います
	widgetVal := reflect.ValueOf(any(b.Widget))
	if !widgetVal.IsValid() { // ウィジェットがnilの場合
		// ターゲットをnilに設定しようと試みる（ターゲットがnil代入可能なら）
		if targetElem.Type().Kind() == reflect.Interface || targetElem.Type().Kind() == reflect.Ptr {
			if !targetElem.IsNil() && targetElem.CanSet() {
				targetElem.Set(reflect.Zero(targetElem.Type()))
			}
		}
		return b.Self
	}

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
		if typeName != nil && typeName.Kind() == reflect.Ptr {
			typeName = typeName.Elem()
		}
		// NOTE: 複数のエラーを結合して返すように変更
		return zero, fmt.Errorf("%s build errors: %w", typeName.Name(), errors.Join(b.errors...))
	}
	b.Widget.MarkDirty(true)
	return b.Widget, nil
}