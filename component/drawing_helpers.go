package component

import (
	"furoshiki/style"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	// whitePixelImg は、頂点カラーで図形を描画する際のソース画像として使用します。
	// ebiten.WhiteImage がv2.2で削除されたため、代替として1x1の白い画像を一度だけ生成して再利用します。
	whitePixelImg *ebiten.Image
	initOnce      sync.Once
)

// ensureWhitePixelImg は、描画に必要な白い1x1画像をスレッドセーフに初期化します。
func ensureWhitePixelImg() {
	initOnce.Do(func() {
		whitePixelImg = ebiten.NewImage(1, 1)
		whitePixelImg.Fill(color.White)
	})
}

// --- Drawing Helper ---

// applyOpacity は元の色に不透明度を適用した新しい色を返します。
// color.NRGBAモデルに変換してアルファ値を操作することで、安全かつ意図通りに動作させます。
func applyOpacity(c color.Color, opacity *float64) color.Color {
	if c == nil || opacity == nil {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	nrgba.A = uint8(float64(nrgba.A) * (*opacity))
	return nrgba
}

// colorToScale は color.Color を ebiten.Vertex で使用する float32 の RGBA スケール値([0, 1])に変換します。
func colorToScale(clr color.Color) (float32, float32, float32, float32) {
	if clr == nil {
		return 0, 0, 0, 0
	}
	// RGBA() は alpha-premultiplied な [0, 65535] の範囲の値を返します。
	// 頂点カラーは [0, 1] の範囲の float であるため、65535.0 で割ります。
	r, g, b, a := clr.RGBA()
	return float32(r) / 65535.0, float32(g) / 65535.0, float32(b) / 65535.0, float32(a) / 65535.0
}

// createRoundedRectPath は、指定された半径で角丸矩形のパスを生成します。
// 角丸の半径が矩形サイズの半分を超える場合の描画崩れを防ぐため、半径を調整します。
func createRoundedRectPath(x, y, width, height, radius float32) *vector.Path {
	path := &vector.Path{}

	maxRadius := float32(math.Min(float64(width/2), float64(height/2)))
	radius = float32(math.Min(float64(radius), float64(maxRadius)))

	if radius <= 0 {
		path.MoveTo(x, y)
		path.LineTo(x+width, y)
		path.LineTo(x+width, y+height)
		path.LineTo(x, y+height)
		path.Close()
		return path
	}

	path.MoveTo(x+radius, y)
	path.LineTo(x+width-radius, y)
	path.QuadTo(x+width, y, x+width, y+radius)
	path.LineTo(x+width, y+height-radius)
	path.QuadTo(x+width, y+height, x+width-radius, y+height)
	path.LineTo(x+radius, y+height)
	path.QuadTo(x, y+height, x, y+height-radius)
	path.LineTo(x, y+radius)
	path.QuadTo(x, y, x+radius, y)
	path.Close()

	return path
}

// drawVectorPath は、vector.Pathを描画するための共通ヘルパー関数です。
// strokeOptsがnilでない場合は線を描画し、nilの場合は図形を塗りつぶします。
// これにより、背景と境界線の描画ロジックにおけるコードの重複を削減します。
func drawVectorPath(dst *ebiten.Image, path *vector.Path, clr color.Color, triOpts *ebiten.DrawTrianglesOptions, strokeOpts *vector.StrokeOptions) {
	var vertices []ebiten.Vertex
	var indices []uint16

	if strokeOpts != nil {
		vertices, indices = path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOpts)
	} else {
		vertices, indices = path.AppendVerticesAndIndicesForFilling(nil, nil)
	}

	if len(vertices) == 0 {
		return
	}

	cr, cg, cb, ca := colorToScale(clr)
	for i := range vertices {
		vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = cr, cg, cb, ca
	}

	dst.DrawTriangles(vertices, indices, whitePixelImg, triOpts)
}

// DrawStyledBackground は、指定されたスタイルでウィジェットの背景と境界線を描画します。
// この関数は、描画ロジックを内部ヘルパー関数(drawBackground, drawBorder)に委譲することで、
// コードの関心事を分離し、可読性を高めています。
func DrawStyledBackground(dst *ebiten.Image, x, y, width, height int, s style.Style) {
	if width <= 0 || height <= 0 {
		return
	}
	ensureWhitePixelImg()

	fx, fy := float32(x), float32(y)
	fw, fh := float32(width), float32(height)

	drawTrianglesOptions := &ebiten.DrawTrianglesOptions{AntiAlias: true}

	drawBackground(dst, fx, fy, fw, fh, s, drawTrianglesOptions)
	drawBorder(dst, fx, fy, fw, fh, s, drawTrianglesOptions)
}

// drawBackground は、ウィジェットの背景を描画する内部ヘルパーです。
func drawBackground(dst *ebiten.Image, x, y, width, height float32, s style.Style, opts *ebiten.DrawTrianglesOptions) {
	bgColorPtr := s.Background
	if bgColorPtr == nil || *bgColorPtr == color.Transparent {
		return
	}

	bgColor := *bgColorPtr
	if s.Opacity != nil {
		bgColor = applyOpacity(bgColor, s.Opacity)
	}

	radius := float32(0)
	if s.BorderRadius != nil {
		radius = *s.BorderRadius
	}

	if radius > 0 {
		// パスを生成し、共通描画ヘルパーを呼び出します（塗りつぶしモード）。
		path := createRoundedRectPath(x, y, width, height, radius)
		drawVectorPath(dst, path, bgColor, opts, nil)
	} else {
		vector.DrawFilledRect(dst, x, y, width, height, bgColor, false)
	}
}

// drawBorder は、ウィジェットの境界線を描画する内部ヘルパーです。
func drawBorder(dst *ebiten.Image, x, y, width, height float32, s style.Style, opts *ebiten.DrawTrianglesOptions) {
	borderColorPtr := s.BorderColor
	borderWidth := float32(0)
	if s.BorderWidth != nil {
		borderWidth = *s.BorderWidth
	}
	if borderColorPtr == nil || *borderColorPtr == color.Transparent || borderWidth <= 0 {
		return
	}

	borderColor := *borderColorPtr
	if s.Opacity != nil {
		borderColor = applyOpacity(borderColor, s.Opacity)
	}

	radius := float32(0)
	if s.BorderRadius != nil {
		radius = *s.BorderRadius
	}

	if radius > 0 {
		// 境界線のパスは、図形の中心に描画されるため、幅の半分だけ内側にオフセットさせます。
		halfBw := borderWidth / 2
		insetPath := createRoundedRectPath(x+halfBw, y+halfBw, width-borderWidth, height-borderWidth, radius-halfBw)

		// 線描画用のオプションを作成し、共通描画ヘルパーを呼び出します。
		strokeOpts := &vector.StrokeOptions{Width: borderWidth, MiterLimit: 10}
		drawVectorPath(dst, insetPath, borderColor, opts, strokeOpts)
	} else {
		vector.StrokeRect(dst, x, y, width, height, borderWidth, borderColor, false)
	}
}

// DrawAlignedText は、指定された矩形領域内にテキストを揃えて描画します。
func DrawAlignedText(screen *ebiten.Image, textContent string, area image.Rectangle, s style.Style) {
	if textContent == "" || s.Font == nil || *s.Font == nil {
		return
	}

	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}

	contentRect := image.Rect(
		area.Min.X+padding.Left,
		area.Min.Y+padding.Top,
		area.Max.X-padding.Right,
		area.Max.Y-padding.Bottom,
	)
	if contentRect.Dx() <= 0 || contentRect.Dy() <= 0 {
		return
	}
	bounds := text.BoundString(*s.Font, textContent)

	var textX int
	textAlign := style.TextAlignLeft
	if s.TextAlign != nil {
		textAlign = *s.TextAlign
	}
	switch textAlign {
	case style.TextAlignCenter:
		textX = contentRect.Min.X + (contentRect.Dx()-bounds.Dx())/2
	case style.TextAlignRight:
		textX = contentRect.Max.X - bounds.Dx()
	default: // style.TextAlignLeft
		textX = contentRect.Min.X
	}

	var textY int
	metrics := (*s.Font).Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	verticalAlign := style.VerticalAlignMiddle
	if s.VerticalAlign != nil {
		verticalAlign = *s.VerticalAlign
	}
	switch verticalAlign {
	case style.VerticalAlignTop:
		textY = contentRect.Min.Y + metrics.Ascent.Ceil()
	case style.VerticalAlignBottom:
		textY = contentRect.Max.Y - metrics.Descent.Ceil()
	default: // style.VerticalAlignMiddle
		textY = contentRect.Min.Y + (contentRect.Dy()-textHeight)/2 + metrics.Ascent.Ceil()
	}

	textColor := color.Color(color.Black)
	if s.TextColor != nil {
		textColor = *s.TextColor
	}
	if s.Opacity != nil {
		textColor = applyOpacity(textColor, s.Opacity)
	}

	text.Draw(screen, textContent, *s.Font, textX, textY, textColor)
}
