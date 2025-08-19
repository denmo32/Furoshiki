package style

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"
)

// Styleはコンポーネントの視覚的プロパティを定義します。
// [修正] ゼロ値（例: 0, nil）と「未設定」を区別できるように、
// 多くのフィールドをポインタ型に変更します。これにより、Merge関数がより柔軟になります。
// [改善] Background, BorderColor, Font, TextColor もポインタ型に変更し、一貫性を高めました。
type Style struct {
	Background   *color.Color
	BorderColor  *color.Color
	BorderWidth  *float32
	Margin       *Insets
	Padding      *Insets
	Font         *font.Face
	TextColor    *color.Color
	BorderRadius *float32
	Opacity      *float64
}

// Insetsはマージンやパディングの四方の値を表します。
type Insets struct {
	Top, Right, Bottom, Left int
}

// Mergeは2つのスタイルをマージします。
// baseスタイルをベースに、overlayスタイルのプロパティがnilでない場合に値を上書きします。
// [修正] ポインタ型のフィールドをチェックするようにロジックを更新しました。
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

	return result
}

// Validate はスタイルの値が有効かどうかを検証します。
// 無効な値が見つかった場合はエラーを返します。
// [修正] ポインタ型のフィールドに対応するため、nilチェックを追加します。
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