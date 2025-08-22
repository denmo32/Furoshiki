package component

import (
    "errors"
    "fmt"
    "furoshiki/style"
    "image/color"
    "reflect"
)

// positionSetter は、ウィジェットが要求位置を設定できるかどうかをチェックする非公開インターフェースです
type positionSetter interface {
    SetRequestedPosition(x, y int)
}

// Builder は、すべてのウィジェットビルダーの汎用基底クラスです。
// サイズ、スタイル、レイアウトプロパティを設定するための共通メソッドを提供します。
// T は具体的なビルダー型（例: *LabelBuilder）です。
// W は構築中のウィジェット型（例: *Label）です。
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

// AbsolutePosition は、親コンテナ内でのウィジェットの希望相対位置を設定します。
// これは、親コンテナがAbsoluteLayout（例: ZStack）を使用している場合にのみ有効です。
func (b *Builder[T, W]) AbsolutePosition(x, y int) T {
    if p, ok := any(b.Widget).(positionSetter); ok {
        p.SetRequestedPosition(x, y)
    } else {
        b.AddError(fmt.Errorf("%T does not support AbsolutePosition", b.Widget))
    }
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
        return fmt.Errorf("size must be non-negative, got %dx%d", width, height)
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
        b.AddError(fmt.Errorf("flex must be non-negative, got %d", flex))
    } else {
        b.Widget.SetFlex(flex)
    }
    return b.Self
}

// --- スタイルヘルパーメソッド ---

// BackgroundColor はウィジェットの背景色を設定します。
func (b *Builder[T, W]) BackgroundColor(c color.Color) T {
    return b.Style(style.Style{Background: style.PColor(c)})
}

// Margin はすべての辺に同じマージン値を設定します。
func (b *Builder[T, W]) Margin(m int) T {
    return b.Style(style.Style{Margin: style.PInsets(style.Insets{Top: m, Right: m, Bottom: m, Left: m})})
}

// MarginInsets は各辺に個別のマージン値を設定します。
func (b *Builder[T, W]) MarginInsets(i style.Insets) T {
    return b.Style(style.Style{Margin: style.PInsets(i)})
}

// Padding はすべての辺に同じパディング値を設定します。
func (b *Builder[T, W]) Padding(p int) T {
    return b.Style(style.Style{Padding: style.PInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})})
}

// PaddingInsets は各辺に個別のパディング値を設定します。
func (b *Builder[T, W]) PaddingInsets(i style.Insets) T {
    return b.Style(style.Style{Padding: style.PInsets(i)})
}

// BorderRadius はウィジェットの角の半径を設定します。
func (b *Builder[T, W]) BorderRadius(radius float32) T {
    return b.Style(style.Style{BorderRadius: style.PFloat32(radius)})
}

// Border はウィジェットの境界線の幅と色を設定します。
func (b *Builder[T, W]) Border(width float32, c color.Color) T {
    if width < 0 {
        b.AddError(fmt.Errorf("border width must be non-negative, got %f", width))
        return b.Self
    }
    return b.Style(style.Style{
        BorderWidth: style.PFloat32(width),
        BorderColor: style.PColor(c),
    })
}

// Build はウィジェットの構築を完了します。
func (b *Builder[T, W]) Build() (W, error) {
    if len(b.errors) > 0 {
        var zero W
        typeName := reflect.TypeOf(b.Widget).Elem().Name()
        return zero, fmt.Errorf("%s build errors: %w", typeName, errors.Join(b.errors...))
    }
    b.Widget.MarkDirty(true)
    return b.Widget, nil
}