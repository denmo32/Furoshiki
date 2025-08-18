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
	// テキストの描画範囲を計算
	bounds := text.BoundString(s.Font, textContent)
	// 水平方向の中央揃え
	textX := contentRect.Min.X + (contentRect.Dx()-bounds.Dx())/2

	// [改善] 垂直方向の中央揃えをより正確に計算
	// font.Metrics を使用して、アセント（ベースラインより上の高さ）とディセント（ベースラインより下の高さ）を取得します。
	metrics := s.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	// テキストの描画基準点（ベースライン）のY座標を計算します。
	// contentRectの中心にテキストの中心が来るように調整し、アセント分を足すことで正しいベースライン位置を求めます。
	textY := contentRect.Min.Y + (contentRect.Dy()-textHeight)/2 + metrics.Ascent.Ceil()

	var textColor color.Color = color.Black
	if s.TextColor != nil {
		textColor = s.TextColor
	}
	text.Draw(screen, textContent, s.Font, textX, textY, textColor)
}