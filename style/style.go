package style

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"
)

// [追加] テキストの水平方向の揃え位置を定義します。
type TextAlignType int

const (
	TextAlignLeft   TextAlignType = iota // 左揃え
	TextAlignCenter                      // 中央揃え
	TextAlignRight                       // 右揃え
)

// [追加] テキストの垂直方向の揃え位置を定義します。
type VerticalAlignType int

const (
	VerticalAlignTop    VerticalAlignType = iota // 上揃え
	VerticalAlignMiddle                          // 中央揃え (垂直)
	VerticalAlignBottom                          // 下揃え
)

// Styleはコンポーネントの視覚的プロパティを定義します。
// ゼロ値（例: 0, nil）と「未設定」を区別できるように、
// 多くのフィールドをポインタ型に変更しています。これにより、Merge関数がより柔軟になります。
type Style struct {
	Background   *color.Color
	BorderColor  *color.Color
	BorderWidth  *float32
	Margin       *Insets
	Padding      *Insets
	Font         *font.Face
	TextColor    *color.Color
	BorderRadius *float32
	Opacity      *float64 // 0.0 (完全に透明) から 1.0 (完全に不透明)
	// [追加] テキストの水平・垂直方向の揃え位置
	TextAlign    *TextAlignType
	VerticalAlign *VerticalAlignType
}

// Insetsはマージンやパディングの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Mergeは2つのスタイルをマージします。
// baseスタイルをベースに、overlayスタイルのプロパティがnilでない場合に値を上書きします。
// これにより、意図的にBorderWidthを0に設定するなどの操作が可能になります。
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
	// [追加] テキスト揃えプロパティのマージ
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

// --- Pointer Helpers ---
// 以下のヘルパー関数は、スタイル設定をより直感的かつシンプルにするために提供されます。
// これらを使用することで、一時変数を宣言することなく、直接スタイル構造体に値を設定できます。

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

// [追加] PTextAlignType は TextAlignType の値からそのポインタを生成して返します。
func PTextAlignType(t TextAlignType) *TextAlignType {
	return &t
}

// [追加] PVerticalAlignType は VerticalAlignType の値からそのポインタを生成して返します。
func PVerticalAlignType(v VerticalAlignType) *VerticalAlignType {
	return &v
}