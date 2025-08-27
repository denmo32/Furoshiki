package widget

import (
	"furoshiki/component"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ScrollBar は、スクロール可能な領域の状態を視覚的に示すウィジェットです。
type ScrollBar struct {
	*component.LayoutableWidget
	trackColor   color.Color
	thumbColor   color.Color
	contentRatio float64
	scrollRatio  float64
}

var _ component.ScrollBarWidget = (*ScrollBar)(nil)

// newScrollBarは、スクロールバーウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//
//	常にNewScrollBarBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newScrollBar() (*ScrollBar, error) {
	s := &ScrollBar{
		trackColor: color.RGBA{220, 220, 220, 255},
		thumbColor: color.RGBA{180, 180, 180, 255},
	}
	s.LayoutableWidget = component.NewLayoutableWidget()
	if err := s.Init(s); err != nil {
		return nil, err
	}
	s.SetSize(10, 100)
	return s, nil
}

// UPDATE: DrawメソッドのシグネチャをDrawInfoを受け取るように変更
// Draw はScrollBarを描画します。
func (s *ScrollBar) Draw(info component.DrawInfo) {
	if !s.IsVisible() || !s.HasBeenLaidOut() {
		return
	}
	x, y := s.GetPosition()
	width, height := s.GetSize()

	// UPDATE: 親から渡されたオフセットを描画座標に適用
	finalX := float32(x + info.OffsetX)
	finalY := float32(y + info.OffsetY)

	vector.DrawFilledRect(info.Screen, finalX, finalY, float32(width), float32(height), s.trackColor, false)

	if s.contentRatio >= 1.0 {
		return
	}
	thumbHeight := float32(float64(height) * s.contentRatio)
	minThumbHeight := float32(10)
	if thumbHeight < minThumbHeight {
		thumbHeight = minThumbHeight
	}
	if height < int(minThumbHeight) {
		return
	}

	thumbYRange := float32(height) - thumbHeight
	thumbY := finalY + thumbYRange*float32(s.scrollRatio)

	vector.DrawFilledRect(info.Screen, finalX, thumbY, float32(width), thumbHeight, s.thumbColor, false)
}

// SetRatios は、つまみのサイズと位置を計算するための比率を設定します。
func (s *ScrollBar) SetRatios(contentRatio, scrollRatio float64) {
	if s.contentRatio != contentRatio || s.scrollRatio != scrollRatio {
		s.contentRatio = contentRatio
		s.scrollRatio = scrollRatio
		s.MarkDirty(false)
	}
}

// --- ScrollBarBuilder ---
type ScrollBarBuilder struct {
	component.Builder[*ScrollBarBuilder, *ScrollBar]
}

func NewScrollBarBuilder() *ScrollBarBuilder {
	s, err := newScrollBar()
	b := &ScrollBarBuilder{}
	b.Init(b, s)
	b.AddError(err)
	return b
}

func (b *ScrollBarBuilder) Build() (*ScrollBar, error) {
	return b.Builder.Build()
}

func (b *ScrollBarBuilder) TrackColor(c color.Color) *ScrollBarBuilder {
	b.Widget.trackColor = c
	return b
}

func (b *ScrollBarBuilder) ThumbColor(c color.Color) *ScrollBarBuilder {
	b.Widget.thumbColor = c
	return b
}
