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
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Interaction
	*component.Visibility
	*component.Dirty

	// UPDATE: hasBeenLaidOutフィールドはVisibilityコンポーネントに統合されたため削除されました。
	// hasBeenLaidOut bool
	contentRatio float64
	scrollRatio  float64
}

// --- Interface implementation verification ---
var _ component.Widget = (*ScrollBar)(nil)
var _ component.ScrollBarWidget = (*ScrollBar)(nil)
var _ component.NodeOwner = (*ScrollBar)(nil)
var _ component.LayoutPropertiesOwner = (*ScrollBar)(nil)
var _ component.VisibilityOwner = (*ScrollBar)(nil)
var _ component.DirtyManager = (*ScrollBar)(nil)
var _ component.AbsolutePositioner = (*ScrollBar)(nil)
var _ event.EventTarget = (*ScrollBar)(nil)
var _ component.EventProcessor = (*ScrollBar)(nil)
var _ component.StyleGetterSetter = (*ScrollBar)(nil) // For ScrollBarBuilder

// newScrollBar creates a new component-based ScrollBar.
func newScrollBar() (*ScrollBar, error) {
	s := &ScrollBar{}
	s.Node = component.NewNode(s)
	s.Transform = component.NewTransform()
	s.LayoutProperties = component.NewLayoutProperties()
	s.Appearance = component.NewAppearance(s)
	s.Interaction = component.NewInteraction(s)
	s.Visibility = component.NewVisibility(s)
	s.Dirty = component.NewDirty()

	// Default styles
	s.SetStyle(style.Style{
		Background:  style.PColor(color.RGBA{220, 220, 220, 255}), // Track color
		BorderColor: style.PColor(color.RGBA{180, 180, 180, 255}), // Thumb color
	})

	s.SetSize(10, 100)
	return s, nil
}

// --- Interface implementations ---

func (s *ScrollBar) GetNode() *component.Node                         { return s.Node }
func (s *ScrollBar) GetLayoutProperties() *component.LayoutProperties { return s.LayoutProperties }
func (s *ScrollBar) Update()                                          {}
func (s *ScrollBar) Cleanup()                                         { s.SetParent(nil) }

// UPDATE: HasBeenLaidOutの実装をVisibilityコンポーネントへの委譲に変更しました。
func (s *ScrollBar) HasBeenLaidOut() bool { return s.Visibility.HasBeenLaidOut() }

func (s *ScrollBar) Draw(info component.DrawInfo) {
	// UPDATE: hasBeenLaidOutのチェックをHasBeenLaidOut()メソッド呼び出しに置き換えました。
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

func (s *ScrollBar) MarkDirty(relayout bool) {
	s.Dirty.MarkDirty(relayout)
	if relayout && !s.IsLayoutBoundary() {
		if parent := s.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (s *ScrollBar) SetPosition(x, y int) {
	// UPDATE: レイアウト済み状態の管理をVisibilityコンポーネントに委譲します。
	if !s.HasBeenLaidOut() {
		s.SetLaidOut(true)
	}
	if posX, posY := s.GetPosition(); posX != x || posY != y {
		s.Transform.SetPosition(x, y)
		s.MarkDirty(false)
	}
}

func (s *ScrollBar) SetSize(width, height int) {
	if w, h := s.GetSize(); w != width || h != height {
		s.Transform.SetSize(width, height)
		s.MarkDirty(true)
	}
}

func (s *ScrollBar) SetMinSize(width, height int) {}
func (s *ScrollBar) GetMinSize() (int, int) {
	return 0, 0
}

func (s *ScrollBar) HitTest(x, y int) component.Widget {
	// TODO: Implement hit testing for thumb dragging
	return nil
}

func (s *ScrollBar) HandleEvent(e *event.Event) {
	// UPDATE: イベント処理の責務をInteractionコンポーネントに委譲
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

// --- AbsolutePositioner Implementation ---
func (s *ScrollBar) SetRequestedPosition(x, y int) {
	s.Transform.SetRequestedPosition(x, y)
	s.MarkDirty(true)
}

func (s *ScrollBar) GetRequestedPosition() (int, int) {
	return s.Transform.GetRequestedPosition()
}

// --- ScrollBarBuilder ---
// UPDATE: 汎用のcomponent.Builderを埋め込むように変更
type ScrollBarBuilder struct {
	component.Builder[*ScrollBarBuilder, *ScrollBar]
}

// NewScrollBarBuilderは新しいScrollBarBuilderを生成します。
func NewScrollBarBuilder() *ScrollBarBuilder {
	sb, err := newScrollBar()
	b := &ScrollBarBuilder{}
	b.Init(b, sb)
	b.AddError(err)
	return b
}

// Build は、最終的なScrollBarを構築して返します。
func (b *ScrollBarBuilder) Build() (*ScrollBar, error) {
	// UPDATE: 処理を基底のBuilder.Buildに完全に委譲
	return b.Builder.Build()
}

// TrackColor はスクロールバーのトラック（背景）の色を設定します。
func (b *ScrollBarBuilder) TrackColor(c color.Color) *ScrollBarBuilder {
	// 汎用Builderのスタイルヘルパーを利用
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
