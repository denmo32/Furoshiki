package component

import (
	"furoshiki/style"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Drawing Helper ---

func DrawStyledBackground(dst *ebiten.Image, x, y, width, height int, s style.Style) {
	if width <= 0 || height <= 0 {
		return
	}
	if s.Background != nil && s.Background != color.Transparent {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), s.Background, false)
	}
	if s.BorderColor != nil && s.BorderWidth > 0 {
		vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), s.BorderWidth, s.BorderColor, false)
	}
}

func DrawAlignedText(screen *ebiten.Image, textContent string, area image.Rectangle, s style.Style) {
	if textContent == "" || s.Font == nil {
		return
	}
	contentRect := image.Rect(
		area.Min.X+s.Padding.Left,
		area.Min.Y+s.Padding.Top,
		area.Max.X-s.Padding.Right,
		area.Max.Y-s.Padding.Bottom,
	)
	if contentRect.Dx() <= 0 || contentRect.Dy() <= 0 {
		return
	}
	bounds := text.BoundString(s.Font, textContent)
	textX := contentRect.Min.X + (contentRect.Dx()-bounds.Dx())/2
	textY := contentRect.Min.Y + (contentRect.Dy()+bounds.Dy())/2
	var textColor color.Color = color.Black
	if s.TextColor != nil {
		textColor = s.TextColor
	}
	text.Draw(screen, textContent, s.Font, textX, textY, textColor)
}