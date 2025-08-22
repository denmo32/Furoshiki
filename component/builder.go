package component

import (
	"errors"
	"fmt"
	"furoshiki/style"
	"image/color"
	"reflect"
)

// positionSetter is an unexported interface used to check if a widget
// can have its requested position set. This is implemented by LayoutableWidget.
type positionSetter interface {
	SetRequestedPosition(x, y int)
}

// Builder is a generic base for all widget builders.
// It provides common methods for setting size, style, and layout properties.
// T is the concrete builder type (e.g., *LabelBuilder).
// W is the widget type being built (e.g., *Label).
type Builder[T any, W Widget] struct {
	Widget W
	errors []error
	Self   T
}

// Init initializes the base builder. It must be called from the concrete builder's constructor.
func (b *Builder[T, W]) Init(self T, widget W) {
	b.Self = self
	b.Widget = widget
}

// AddError adds a build error to the builder.
func (b *Builder[T, W]) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// AbsolutePosition sets the widget's desired relative position within its parent.
// This is only effective if the parent container uses an AbsoluteLayout (e.g., ZStack).
func (b *Builder[T, W]) AbsolutePosition(x, y int) T {
	if p, ok := any(b.Widget).(positionSetter); ok {
		p.SetRequestedPosition(x, y)
	} else {
		b.AddError(fmt.Errorf("%T does not support AbsolutePosition", b.Widget))
	}
	return b.Self
}

// Size sets the widget's size.
func (b *Builder[T, W]) Size(width, height int) T {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("size must be non-negative, got %dx%d", width, height))
	} else {
		b.Widget.SetSize(width, height)
	}
	return b.Self
}

// MinSize sets the widget's minimum size.
func (b *Builder[T, W]) MinSize(width, height int) T {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("min size must be non-negative, got %dx%d", width, height))
	} else {
		b.Widget.SetMinSize(width, height)
	}
	return b.Self
}

// Style merges the given style with the widget's existing style.
func (b *Builder[T, W]) Style(s style.Style) T {
	existingStyle := b.Widget.GetStyle()
	b.Widget.SetStyle(style.Merge(existingStyle, s))
	return b.Self
}

// Flex sets how the widget should stretch or shrink in a FlexLayout.
func (b *Builder[T, W]) Flex(flex int) T {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("flex must be non-negative, got %d", flex))
	} else {
		b.Widget.SetFlex(flex)
	}
	return b.Self
}

// --- Style Helpers ---

// BackgroundColor sets the widget's background color.
func (b *Builder[T, W]) BackgroundColor(c color.Color) T {
	return b.Style(style.Style{Background: style.PColor(c)})
}

// Margin sets the same margin value for all four sides.
func (b *Builder[T, W]) Margin(m int) T {
	return b.Style(style.Style{Margin: style.PInsets(style.Insets{Top: m, Right: m, Bottom: m, Left: m})})
}

// MarginInsets sets individual margin values for each side.
func (b *Builder[T, W]) MarginInsets(i style.Insets) T {
	return b.Style(style.Style{Margin: style.PInsets(i)})
}

// Padding sets the same padding value for all four sides.
func (b *Builder[T, W]) Padding(p int) T {
	return b.Style(style.Style{Padding: style.PInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})})
}

// PaddingInsets sets individual padding values for each side.
func (b *Builder[T, W]) PaddingInsets(i style.Insets) T {
	return b.Style(style.Style{Padding: style.PInsets(i)})
}

// BorderRadius sets the radius of the widget's corners.
func (b *Builder[T, W]) BorderRadius(radius float32) T {
	return b.Style(style.Style{BorderRadius: style.PFloat32(radius)})
}

// Border sets the widget's border width and color.
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

// Build finalizes the widget construction.
func (b *Builder[T, W]) Build() (W, error) {
	if len(b.errors) > 0 {
		var zero W
		// reflectを使用してウィジェットの型名を取得
		typeName := reflect.TypeOf(b.Widget).Elem().Name()
		return zero, fmt.Errorf("%s build errors: %w", typeName, errors.Join(b.errors...))
	}
	b.Widget.MarkDirty(true)
	return b.Widget, nil
}
