package style

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"
)

// Styleはコンポーネントの視覚的プロパティを定義します。
type Style struct {
	Background   color.Color
	BorderColor  color.Color
	BorderWidth  float32
	Margin       Insets
	Padding      Insets
	Font         font.Face
	TextColor    color.Color
	BorderRadius float32
	Opacity      float64
}

// Insetsはマージンやパディングの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Mergeは2つのスタイルをマージします。
// baseスタイルをベースに、overlayスタイルの値を上書きします。
// overlayスタイルでゼロ値（nil, 0, 空の構造体など）が指定されているフィールドは、
// 「設定されていない」と見なされ、baseスタイルの値が保持されます。
func Merge(base, overlay Style) Style {
	result := base

	// Background
	if overlay.Background != nil {
		result.Background = overlay.Background
	}

	// BorderColor
	if overlay.BorderColor != nil {
		result.BorderColor = overlay.BorderColor
	}

	// BorderWidth
	// [注意] ゼロ値チェックのため、overlay.BorderWidthを意図的に0に設定して
	// baseの値を上書きすることはできません。
	if overlay.BorderWidth != 0 {
		result.BorderWidth = overlay.BorderWidth
	}

	// Margin
	// [注意] overlay.Marginがゼロ値(Insets{})の場合、baseの値が維持されます。
	// 特定の辺のマージンのみを0にしたい場合は、overlay側で他の辺の値を明示的に設定する必要があります。
	if overlay.Margin != (Insets{}) {
		result.Margin = overlay.Margin
	}

	// Padding
	if overlay.Padding != (Insets{}) {
		result.Padding = overlay.Padding
	}

	// Font
	if overlay.Font != nil {
		result.Font = overlay.Font
	}

	// TextColor
	if overlay.TextColor != nil {
		result.TextColor = overlay.TextColor
	}

	// BorderRadius
	// [注意] ゼロ値チェックのため、overlay.BorderRadiusを意図的に0に設定して
	// baseの値を上書きすることはできません。
	if overlay.BorderRadius != 0 {
		result.BorderRadius = overlay.BorderRadius
	}

	// Opacity
	// [注意] ゼロ値チェックのため、overlay.Opacityを意図的に0.0に設定して
	// baseの値を上書きすることはできません。
	if overlay.Opacity != 0 {
		result.Opacity = overlay.Opacity
	}

	return result
}

// Validate はスタイルの値が有効かどうかを検証します。
// 無効な値が見つかった場合はエラーを返します。
func (s *Style) Validate() error {
	if s.BorderWidth < 0 {
		return fmt.Errorf("border width must be non-negative, got %f", s.BorderWidth)
	}

	if s.BorderRadius < 0 {
		return fmt.Errorf("border radius must be non-negative, got %f", s.BorderRadius)
	}

	if s.Opacity < 0 || s.Opacity > 1.0 {
		return fmt.Errorf("opacity must be between 0.0 and 1.0, got %f", s.Opacity)
	}

	if s.Margin.Top < 0 || s.Margin.Right < 0 || s.Margin.Bottom < 0 || s.Margin.Left < 0 {
		return fmt.Errorf("margin values must be non-negative, got %v", s.Margin)
	}

	if s.Padding.Top < 0 || s.Padding.Right < 0 || s.Padding.Bottom < 0 || s.Padding.Left < 0 {
		return fmt.Errorf("padding values must be non-negative, got %v", s.Padding)
	}

	return nil
}