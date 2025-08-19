package component

import (
	"furoshiki/style"
	"image"
	"image/color"
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

// ensureWhitePixelImg は、描画に必要な白い1x1画像を初期化します。
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
	// NRGBAに変換してアルファ値を操作するのが最も安全
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	nrgba.A = uint8(float64(nrgba.A) * (*opacity))
	return nrgba
}

// colorToScale は color.Color を ebiten.Vertex で使用する float32 の RGBA スケール値に変換します。
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
// vectorパッケージのパス機能を使って図形を定義します。
func createRoundedRectPath(x, y, width, height, radius float32) *vector.Path {
	path := &vector.Path{}

	// 半径が矩形の幅や高さの半分を超える場合は調整
	maxRadius := width / 2
	if height/2 < maxRadius {
		maxRadius = height / 2
	}
	if radius > maxRadius {
		radius = maxRadius
	}

	if radius <= 0 {
		// 角丸が不要な場合は、単純な四角形のパスを生成
		path.MoveTo(x, y)
		path.LineTo(x+width, y)
		path.LineTo(x+width, y+height)
		path.LineTo(x, y+height)
		path.Close()
		return path
	}

	// 4つの角を円弧で結んで角丸矩形のパスを生成
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

// DrawStyledBackground は、指定されたスタイルでウィジェットの背景と境界線を描画します。
// [改善] Opacity(不透明度)とBorderRadius(角丸)に対応しました。
// [修正] Ebitengineのバージョンに依存しにくい、基本的な描画APIのみを使用するように修正しました。
func DrawStyledBackground(dst *ebiten.Image, x, y, width, height int, s style.Style) {
	if width <= 0 || height <= 0 {
		return
	}
	// 描画に必要な1x1の白ピクセル画像を準備
	ensureWhitePixelImg()

	// スタイルから値を取得 (nilの場合はゼロ値を使用)
	bgColorPtr := s.Background
	borderColorPtr := s.BorderColor
	borderWidth := float32(0)
	if s.BorderWidth != nil {
		borderWidth = *s.BorderWidth
	}
	radius := float32(0)
	if s.BorderRadius != nil {
		radius = *s.BorderRadius
	}

	// DrawTrianglesのオプションは共通
	drawTrianglesOptions := &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	}

	// 背景色の描画
	if bgColorPtr != nil && *bgColorPtr != color.Transparent {
		bgColor := *bgColorPtr
		if s.Opacity != nil {
			bgColor = applyOpacity(bgColor, s.Opacity)
		}

		if radius > 0 {
			// 角丸矩形のパスを生成して塗りつぶし
			path := createRoundedRectPath(float32(x), float32(y), float32(width), float32(height), radius)
			// 塗りつぶし用の頂点とインデックスを生成
			vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

			// 全ての頂点に背景色を設定
			cr, cg, cb, ca := colorToScale(bgColor)
			for i := range vertices {
				vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = cr, cg, cb, ca
			}

			// 三角形を描画
			dst.DrawTriangles(vertices, indices, whitePixelImg, drawTrianglesOptions)

		} else {
			// 通常の矩形を描画
			vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), bgColor, false)
		}
	}

	// 境界線の描画
	if borderColorPtr != nil && *borderColorPtr != color.Transparent && borderWidth > 0 {
		borderColor := *borderColorPtr
		if s.Opacity != nil {
			borderColor = applyOpacity(borderColor, s.Opacity)
		}

		if radius > 0 {
			// 境界線のパスは、図形の中心に描画されるため、幅の半分だけ内側にオフセットさせる
			halfBw := borderWidth / 2
			fx, fy := float32(x)+halfBw, float32(y)+halfBw
			fw, fh := float32(width)-borderWidth, float32(height)-borderWidth
			r := radius - halfBw
			if r < 0 {
				r = 0
			}
			insetPath := createRoundedRectPath(fx, fy, fw, fh, r)

			// 線描画用の頂点とインデックスを生成
			strokeOpts := &vector.StrokeOptions{Width: borderWidth, MiterLimit: 10}
			vertices, indices := insetPath.AppendVerticesAndIndicesForStroke(nil, nil, strokeOpts)

			// 全ての頂点に境界線の色を設定
			cr, cg, cb, ca := colorToScale(borderColor)
			for i := range vertices {
				vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = cr, cg, cb, ca
			}

			// 三角形を描画して線を描画
			dst.DrawTriangles(vertices, indices, whitePixelImg, drawTrianglesOptions)

		} else {
			// 通常の矩形の境界線を描画
			vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), borderWidth, borderColor, false)
		}
	}
}

// DrawAlignedText は、指定された矩形領域内にテキストを揃えて描画します。
// [修正] スタイルプロパティがポインタになったため、nilチェックを追加して安全に値を取得します。
func DrawAlignedText(screen *ebiten.Image, textContent string, area image.Rectangle, s style.Style) {
	// [改善] s.Fontがポインタになったため、nilチェックを追加。
	if textContent == "" || s.Font == nil || *s.Font == nil {
		return
	}

	// パディングの値を取得（nilの場合は0として扱う）
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
	// テキストの描画範囲を計算
	bounds := text.BoundString(*s.Font, textContent)
	// 水平方向の中央揃え
	textX := contentRect.Min.X + (contentRect.Dx()-bounds.Dx())/2

	// [改善] 垂直方向の中央揃えをより正確に計算
	// font.Metrics を使用して、アセント（ベースラインより上の高さ）とディセント（ベースラインより下の高さ）を取得します。
	metrics := (*s.Font).Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	// テキストの描画基準点（ベースライン）のY座標を計算します。
	// contentRectの中心にテキストの中心が来るように調整し、アセント分を足すことで正しいベースライン位置を求めます。
	textY := contentRect.Min.Y + (contentRect.Dy()-textHeight)/2 + metrics.Ascent.Ceil()

	// [改善] s.TextColorがポインタになったため、nilチェックとデリファレンスを追加。
	textColor := color.Color(color.Black) // デフォルトは黒
	if s.TextColor != nil {
		textColor = *s.TextColor
	}
	// [追加] テキストにもOpacityを適用
	if s.Opacity != nil {
		textColor = applyOpacity(textColor, s.Opacity)
	}
	text.Draw(screen, textContent, *s.Font, textX, textY, textColor)
}