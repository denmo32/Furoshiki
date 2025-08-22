package style

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"
)

// TextAlignType はテキストの水平方向の揃え位置を定義します。
type TextAlignType int

const (
	TextAlignLeft   TextAlignType = iota // 左揃え
	TextAlignCenter                      // 中央揃え
	TextAlignRight                       // 右揃え
)

// VerticalAlignType はテキストの垂直方向の揃え位置を定義します。
type VerticalAlignType int

const (
	VerticalAlignTop    VerticalAlignType = iota // 上揃え
	VerticalAlignMiddle                          // 中央揃え (垂直)
	VerticalAlignBottom                          // 下揃え
)

// Styleはコンポーネントの視覚的プロパティを定義します。
// 多くのフィールドがポインタ型になっています。これにより、ゼロ値（例: 0, nil）と
// 「未設定」の状態を区別でき、Merge関数がより柔軟に動作します。
type Style struct {
	Background    *color.Color
	BorderColor   *color.Color
	BorderWidth   *float32
	Margin        *Insets
	Padding       *Insets
	Font          *font.Face
	TextColor     *color.Color
	BorderRadius  *float32
	Opacity       *float64 // 0.0 (完全に透明) から 1.0 (完全に不透明)
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
// 無効な値が見つかった場合はエラーを返します。
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
	if s.Margin != nil && (s.Margin.Top < 0 || s.Margin.Right < 0 || s.Margin.Bottom < 0 || s.Margin.Left < 0) {
		return fmt.Errorf("margin values must be non-negative, got %v", *s.Margin)
	}
	if s.Padding != nil && (s.Padding.Top < 0 || s.Padding.Right < 0 || s.Padding.Bottom < 0 || s.Padding.Left < 0) {
		return fmt.Errorf("padding values must be non-negative, got %v", *s.Padding)
	}
	return nil
}

// DeepCopy はスタイルのディープコピーを生成します。
// ポインタフィールドが指す先の値もコピーされるため、
// コピー元のスタイルに影響を与えることなく安全に変更できます。
// font.Faceはインターフェースでありコピーできないため、ポインタは共有されます。
func (s Style) DeepCopy() Style {
	newStyle := s // シャローコピーで基本構造をコピー

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
	// s.Font (*font.Face) はインターフェースなのでディープコピーしない（できない）
	// ポインタのコピー（シャローコピー）のままとする
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

	return newStyle
}

// --- Pointer Helpers ---
// 以下のヘルパー関数は、スタイル設定をより直感的かつシンプルにするために提供されます。
// これらを使用することで、一時変数を宣言することなく、直接スタイル構造体に値を設定できます。
// 例: style.Style{ Background: style.PColor(color.White) }

// PColor は color.Color の値からそのポインタを生成して返します。
func PColor(c color.Color) *color.Color {
	return &c
}

// PFloat32 は float32 の値からそのポインタを生成して返します。
func PFloat32(f float32) *float32 {
	return &f
}

// PFloat64 は float64 の値からそのポインタを生成して返します。
func PFloat64(f float64) *float64 {
	return &f
}

// PInsets は Insets の値からそのポインタを生成して返します。
func PInsets(i Insets) *Insets {
	return &i
}

// PFont は font.Face の値からそのポインタを生成して返します。
func PFont(f font.Face) *font.Face {
	return &f
}

// PTextAlignType は TextAlignType の値からそのポインタを生成して返します。
func PTextAlignType(t TextAlignType) *TextAlignType {
	return &t
}

// PVerticalAlignType は VerticalAlignType の値からそのポインタを生成して返します。
func PVerticalAlignType(v VerticalAlignType) *VerticalAlignType {
	return &v
}