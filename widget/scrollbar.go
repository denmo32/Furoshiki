package widget

import (
	"furoshiki/component"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

// NewScrollBarは、スクロールバーウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: 内部のInit呼び出しが失敗する可能性があるため、コンストラクタはerrorを返すように変更されました。
func NewScrollBar() (*ScrollBar, error) {
	s := &ScrollBar{
		trackColor: color.RGBA{220, 220, 220, 255},
		thumbColor: color.RGBA{180, 180, 180, 255},
	}
	s.LayoutableWidget = component.NewLayoutableWidget()
	// NOTE: Initがエラーを返すようになったため、エラーをチェックし、呼び出し元に伝播させます。
	if err := s.Init(s); err != nil {
		return nil, err
	}
	s.SetSize(10, 100)
	return s, nil
}

// Draw はScrollBarを描画します。
func (s *ScrollBar) Draw(screen *ebiten.Image) {
	if !s.IsVisible() || !s.HasBeenLaidOut() {
		return
	}
	x, y := s.GetPosition()
	width, height := s.GetSize()

	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), s.trackColor, false)

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
	thumbY := float32(y) + thumbYRange*float32(s.scrollRatio)

	vector.DrawFilledRect(screen, float32(x), thumbY, float32(width), thumbHeight, s.thumbColor, false)
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
	s, err := NewScrollBar()
	b := &ScrollBarBuilder{}
	b.Init(b, s)
	// NOTE: コンストラクタで発生した初期化エラーをビルダーに追加します。
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