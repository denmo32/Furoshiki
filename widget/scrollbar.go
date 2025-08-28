package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ScrollBar is a widget that indicates the state of a scrollable area.
type ScrollBar struct {
	// UPDATE: 複数の共通コンポーネントをWidgetCoreに集約
	*component.WidgetCore
	*component.Appearance
	*component.Interaction

	// ScrollBar specific fields
	contentRatio float64
	scrollRatio  float64
}

// --- Interface implementation verification ---
// UPDATE: StandardWidgetインターフェースの実装を検証
var _ component.StandardWidget = (*ScrollBar)(nil)
var _ component.ScrollBarWidget = (*ScrollBar)(nil) // Implements the more specific scrollbar interface
var _ event.EventTarget = (*ScrollBar)(nil)

// newScrollBar creates a new component-based ScrollBar.
func newScrollBar() (*ScrollBar, error) {
	s := &ScrollBar{}
	// UPDATE: WidgetCoreと、ScrollBarに特有のコンポーネントを初期化
	s.WidgetCore = component.NewWidgetCore(s)
	s.Appearance = component.NewAppearance(s)
	s.Interaction = component.NewInteraction(s)

	// Default styles
	s.SetStyle(style.Style{
		Background:  style.PColor(color.RGBA{220, 220, 220, 255}), // Track color
		BorderColor: style.PColor(color.RGBA{180, 180, 180, 255}), // Thumb color
	})

	s.SetSize(10, 100)
	return s, nil
}

// --- Interface implementations ---

// GetNode, GetLayoutProperties are now inherited from WidgetCore.
func (s *ScrollBar) Update()  {}
func (s *ScrollBar) Cleanup() { s.SetParent(nil) }

func (s *ScrollBar) Draw(info component.DrawInfo) {
	if !s.IsVisible() || !s.HasBeenLaidOut() {
		return
	}
	x, y := s.GetPosition()
	width, height := s.GetSize()

	finalX := float32(x + info.OffsetX)
	finalY := float32(y + info.OffsetY)

	st := s.ReadOnlyStyle()
	trackColor := color.RGBA{220, 220, 220, 255} // Default color
	if st.Background != nil {
		r, g, b, a := (*st.Background).RGBA()
		trackColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	thumbColor := color.RGBA{180, 180, 180, 255} // Default color
	if st.BorderColor != nil {
		r, g, b, a := (*st.BorderColor).RGBA()
		thumbColor = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}

	// Draw track
	vector.DrawFilledRect(info.Screen, finalX, finalY, float32(width), float32(height), trackColor, false)

	if s.contentRatio >= 1.0 {
		return
	}

	// Draw thumb
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

	vector.DrawFilledRect(info.Screen, finalX, thumbY, float32(width), thumbHeight, thumbColor, false)
}

// UPDATE: MarkDirty, SetPosition, SetSize, SetMinSize, GetMinSize, SetRequestedPosition,
// GetRequestedPositionはWidgetCoreに実装されているため削除。
// ScrollBarはコンテンツを持たないため、GetMinSizeはWidgetCoreのデフォルト実装で十分です。

func (s *ScrollBar) HitTest(x, y int) component.Widget {
	// TODO: Implement hit testing for thumb dragging
	return nil
}

func (s *ScrollBar) HandleEvent(e *event.Event) {
	s.Interaction.TriggerHandlers(e)
	// TODO: Implement event handling for thumb dragging
}

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
	sb, err := newScrollBar()
	b := &ScrollBarBuilder{}
	b.Init(b, sb)
	b.AddError(err)
	return b
}

func (b *ScrollBarBuilder) Build() (*ScrollBar, error) {
	return b.Builder.Build()
}

// TrackColor はスクロールバーのトラック（背景）の色を設定します。
func (b *ScrollBarBuilder) TrackColor(c color.Color) *ScrollBarBuilder {
	return b.BackgroundColor(c)
}

// ThumbColor はスクロールバーのつまみ（前景）の色を設定します。
func (b *ScrollBarBuilder) ThumbColor(c color.Color) *ScrollBarBuilder {
	// Thumbの色としてBorderColorを利用する
	st := b.Widget.GetStyle()
	st.BorderColor = style.PColor(c)
	b.Widget.SetStyle(st)
	return b
}