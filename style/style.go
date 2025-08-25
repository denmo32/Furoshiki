package style

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"
)

// TextAlignType はテキストの水平方向の揃え位置を定義します。
type TextAlignType int

const (
	TextAlignLeft TextAlignType = iota
	TextAlignCenter
	TextAlignRight
)

// VerticalAlignType はテキストの垂直方向の揃え位置を定義します。
type VerticalAlignType int

const (
	VerticalAlignTop VerticalAlignType = iota
	VerticalAlignMiddle
	VerticalAlignBottom
)

// Styleはコンポーネントの視覚的プロパティを定義します。
// 多くのフィールドがポインタ型になっており、「未設定」の状態を区別できます。
type Style struct {
	Background    *color.Color
	BorderColor   *color.Color
	BorderWidth   *float32
	Margin        *Insets
	Padding       *Insets
	Font          *font.Face
	TextColor     *color.Color
	BorderRadius  *float32
	Opacity       *float64
	TextAlign     *TextAlignType
	VerticalAlign *VerticalAlignType
}

// Insetsはマージンやパディングの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Mergeは2つのスタイルをマージします。
// baseスタイルをベースに、overlayスタイルのプロパティがnilでない場合に値を上書きします。
func Merge(base, overlay Style) Style {
	result := base
	if overlay.Background != nil {
		result.Background = overlay.Background
	}
	if overlay.BorderColor != nil {
		result.BorderColor = overlay.BorderColor
	}
	if overlay.BorderWidth != nil {
		result.BorderWidth = overlay.BorderWidth
	}
	if overlay.Margin != nil {
		result.Margin = overlay.Margin
	}
	if overlay.Padding != nil {
		result.Padding = overlay.Padding
	}
	if overlay.Font != nil {
		result.Font = overlay.Font
	}
	if overlay.TextColor != nil {
		result.TextColor = overlay.TextColor
	}
	if overlay.BorderRadius != nil {
		result.BorderRadius = overlay.BorderRadius
	}
	if overlay.Opacity != nil {
		result.Opacity = overlay.Opacity
	}
	if overlay.TextAlign != nil {
		result.TextAlign = overlay.TextAlign
	}
	if overlay.VerticalAlign != nil {
		result.VerticalAlign = overlay.VerticalAlign
	}
	return result
}

// Validate はスタイルの値が有効かどうかを検証します。
func (s *Style) Validate() error {
	if s.BorderWidth != nil && *s.BorderWidth < 0 {
		return fmt.Errorf("border width must be non-negative, got %f", *s.BorderWidth)
	}
	if s.BorderRadius != nil && *s.BorderRadius < 0 {
		return fmt.Errorf("border radius must be non-negative, got %f", *s.BorderRadius)
	}
	if s.Opacity != nil && (*s.Opacity < 0 || *s.Opacity > 1.0) {
		return fmt.Errorf("opacity must be between 0.0 and 1.0, got %f", *s.Opacity)
	}
	return nil
}

// DeepCopy はスタイルのディープコピーを生成します。
// これにより、コピー元のスタイルに影響を与えることなく安全に変更できます。
func (s Style) DeepCopy() Style {
	newStyle := s
	if s.Background != nil {
		newStyle.Background = PColor(*s.Background)
	}
	if s.BorderColor != nil {
		newStyle.BorderColor = PColor(*s.BorderColor)
	}
	if s.BorderWidth != nil {
		newStyle.BorderWidth = PFloat32(*s.BorderWidth)
	}
	if s.Margin != nil {
		newStyle.Margin = PInsets(*s.Margin)
	}
	if s.Padding != nil {
		newStyle.Padding = PInsets(*s.Padding)
	}
	if s.TextColor != nil {
		newStyle.TextColor = PColor(*s.TextColor)
	}
	if s.BorderRadius != nil {
		newStyle.BorderRadius = PFloat32(*s.BorderRadius)
	}
	if s.Opacity != nil {
		newStyle.Opacity = PFloat64(*s.Opacity)
	}
	if s.TextAlign != nil {
		newStyle.TextAlign = PTextAlignType(*s.TextAlign)
	}
	if s.VerticalAlign != nil {
		newStyle.VerticalAlign = PVerticalAlignType(*s.VerticalAlign)
	}
	// s.Font (*font.Face) はインターフェースなのでディープコピーしない
	return newStyle
}

// --- Pointer Helpers ---
// これらを使用することで、一時変数を宣言することなく、直接スタイル構造体に値を設定できます。
// 例: style.Style{ Background: style.PColor(color.White) }

func PColor(c color.Color) *color.Color                         { return &c }
func PFloat32(f float32) *float32                               { return &f }
func PFloat64(f float64) *float64                               { return &f }
func PInsets(i Insets) *Insets                                  { return &i }
func PFont(f font.Face) *font.Face                              { return &f }
func PTextAlignType(t TextAlignType) *TextAlignType             { return &t }
func PVerticalAlignType(v VerticalAlignType) *VerticalAlignType { return &v }
