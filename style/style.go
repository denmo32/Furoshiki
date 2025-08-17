package style

import (
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
func Merge(base, overlay Style) Style {
	result := base
	if overlay.Background != nil {
		result.Background = overlay.Background
	}
	if overlay.BorderColor != nil {
		result.BorderColor = overlay.BorderColor
	}
	if overlay.BorderWidth != 0 {
		result.BorderWidth = overlay.BorderWidth
	}
	if overlay.Margin != (Insets{}) {
		result.Margin = overlay.Margin
	}
	if overlay.Padding != (Insets{}) {
		result.Padding = overlay.Padding
	}
	if overlay.Font != nil {
		result.Font = overlay.Font
	}
	if overlay.TextColor != nil {
		result.TextColor = overlay.TextColor
	}
	if overlay.BorderRadius != 0 {
		result.BorderRadius = overlay.BorderRadius
	}
	if overlay.Opacity != 0 {
		result.Opacity = overlay.Opacity
	}
	return result
}