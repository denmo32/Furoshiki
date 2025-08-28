package widget

import (
	"errors"
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/utils"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
)

// Labelはテキストを表示するためのシンプルなウィジェットです。
// 以前のTextWidgetを埋め込む代わりに、必要なコンポーネントを直接合成します。
type Label struct {
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Text
	*component.Visibility
	*component.Dirty

	hasBeenLaidOut bool
	minSize        component.Size
}

// --- インターフェース実装の検証 ---
var _ component.Widget = (*Label)(nil)
var _ component.NodeOwner = (*Label)(nil)
var _ component.AppearanceOwner = (*Label)(nil)
var _ component.TextOwner = (*Label)(nil)
var _ component.LayoutPropertiesOwner = (*Label)(nil)
var _ component.VisibilityOwner = (*Label)(nil)
var _ component.DirtyManager = (*Label)(nil)
var _ component.HeightForWider = (*Label)(nil)

var _ component.AbsolutePositioner = (*Label)(nil)

// newLabelは、新しいコンポーネントベースのLabelを生成します。
func newLabel(text string) (*Label, error) {
	l := &Label{}
	l.Node = component.NewNode(l)
	l.Transform = component.NewTransform()
	l.LayoutProperties = component.NewLayoutProperties()
	l.Appearance = component.NewAppearance(l)
	l.Text = component.NewText(l, text)
	l.Visibility = component.NewVisibility(l)
	l.Dirty = component.NewDirty()
	t := theme.GetCurrent()
	l.SetStyle(t.Label.Default)
	l.SetSize(100, 30)
	return l, nil
}

// --- インターフェース実装 ---

func (l *Label) GetNode() *component.Node                { return l.Node }
func (l *Label) GetLayoutProperties() *component.LayoutProperties { return l.LayoutProperties }
func (l *Label) Update()                                  {}
func (l *Label) Cleanup()                                 { l.SetParent(nil) }

func (l *Label) Draw(info component.DrawInfo) {
	if !l.IsVisible() || !l.hasBeenLaidOut {
		return
	}
	x, y := l.GetPosition()
	width, height := l.GetSize()
	finalX := x + info.OffsetX
	finalY := y + info.OffsetY
	styleToUse := l.GetStyle()
	component.DrawStyledBackground(info.Screen, finalX, finalY, width, height, styleToUse)
	finalRect := image.Rect(finalX, finalY, finalX+width, finalY+height)
	component.DrawAlignedText(info.Screen, l.Text.Text(), finalRect, styleToUse, l.WrapText())
}

func (l *Label) MarkDirty(relayout bool) {
	l.Dirty.MarkDirty(relayout)
	if relayout && !l.IsLayoutBoundary() {
		if parent := l.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (l *Label) SetPosition(x, y int) {
	l.Transform.SetPosition(x, y)
	l.hasBeenLaidOut = true
	l.MarkDirty(false)
}

func (l *Label) SetSize(width, height int) {
	if w, h := l.GetSize(); w != width || h != height {
		l.Transform.SetSize(width, height)
		l.MarkDirty(true)
	}
}

func (l *Label) SetMinSize(width, height int) {
	l.minSize.Width = width
	l.minSize.Height = height
}

func (l *Label) GetMinSize() (int, int) {
	contentMinWidth, contentMinHeight := l.calculateContentMinSize()
	return max(contentMinWidth, l.minSize.Width), max(contentMinHeight, l.minSize.Height)
}

func (l *Label) GetHeightForWidth(width int) int {
	if !l.WrapText() {
		_, h := l.calculateContentMinSize()
		return h
	}
	s := l.ReadOnlyStyle()
	if l.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	contentWidth := width - padding.Left - padding.Right
	if contentWidth <= 0 {
		_, h := l.calculateContentMinSize()
		return h
	}
	_, requiredHeight := component.CalculateWrappedText(*s.Font, l.Text.Text(), contentWidth)
	return requiredHeight + padding.Top + padding.Bottom
}

func (l *Label) calculateContentMinSize() (int, int) {
	s := l.ReadOnlyStyle()
	if l.Text.Text() == "" || s.Font == nil || *s.Font == nil {
		return 0, 0
	}
	padding := style.Insets{}
	if s.Padding != nil {
		padding = *s.Padding
	}
	metrics := (*s.Font).Metrics()
	contentMinHeight := (metrics.Ascent + metrics.Descent).Ceil() + padding.Top + padding.Bottom
	if l.WrapText() {
		longestWord := ""
		words := utils.SplitIntoWords(l.Text.Text())
		for _, word := range words {
			if len(word) > len(longestWord) {
				longestWord = word
			}
		}
		if longestWord == "" {
			longestWord = l.Text.Text()
		}
		bounds := text.BoundString(*s.Font, longestWord)
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	} else {
		bounds := text.BoundString(*s.Font, l.Text.Text())
		contentMinWidth := bounds.Dx() + padding.Left + padding.Right
		return contentMinWidth, contentMinHeight
	}
}

func (l *Label) HitTest(x, y int) component.Widget {
	return nil
}

// --- AbsolutePositioner Implementation ---

func (l *Label) SetRequestedPosition(x, y int) {
	l.Transform.SetRequestedPosition(x, y)
	l.MarkDirty(true)
}

func (l *Label) GetRequestedPosition() (int, int) {
	return l.Transform.GetRequestedPosition()
}

// --- LabelBuilder ---

// LabelBuilderは、新しいコンポーネントベースのLabelウィジェットを構築します。
// 以前の汎用ビルダーとは異なり、Label専用のビルダーとして実装されています。
type LabelBuilder struct {
	label  *Label
	errors []error
}

// NewLabelBuilderは新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	label, err := newLabel("")
	b := &LabelBuilder{label: label}
	if err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

// Buildは、最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}
	return b.label, nil
}

func (b *LabelBuilder) AddError(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

// --- Builder Methods ---

func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.label.SetText(text)
	return b
}

func (b *LabelBuilder) WrapText(wrap bool) *LabelBuilder {
	b.label.SetWrapText(wrap)
	return b
}

func (b *LabelBuilder) Size(width, height int) *LabelBuilder {
	b.label.SetSize(width, height)
	return b
}

func (b *LabelBuilder) MinSize(width, height int) *LabelBuilder {
	b.label.SetMinSize(width, height)
	return b
}

func (b *LabelBuilder) Flex(flex int) *LabelBuilder {
	b.label.SetFlex(flex)
	return b
}

func (b *LabelBuilder) AbsolutePosition(x, y int) *LabelBuilder {
	b.label.SetRequestedPosition(x, y)
	return b
}

func (b *LabelBuilder) AssignTo(target **Label) *LabelBuilder {
	if target == nil {
		b.errors = append(b.errors, errors.New("AssignTo target cannot be nil"))
		return b
	}
	*target = b.label
	return b
}

// --- Style Helper Methods ---

func (b *LabelBuilder) applyStyle(s style.Style) *LabelBuilder {
	existingStyle := b.label.GetStyle()
	b.label.SetStyle(style.Merge(existingStyle, s))
	return b
}

func (b *LabelBuilder) TextColor(c color.Color) *LabelBuilder {
	return b.applyStyle(style.Style{TextColor: style.PColor(c)})
}

func (b *LabelBuilder) BackgroundColor(c color.Color) *LabelBuilder {
	return b.applyStyle(style.Style{Background: style.PColor(c)})
}

func (b *LabelBuilder) Padding(p int) *LabelBuilder {
	return b.applyStyle(style.Style{Padding: style.PInsets(style.Insets{Top: p, Right: p, Bottom: p, Left: p})})
}

func (b *LabelBuilder) Border(width float32, c color.Color) *LabelBuilder {
	return b.applyStyle(style.Style{BorderWidth: style.PFloat32(width), BorderColor: style.PColor(c)})
}

func (b *LabelBuilder) TextAlign(align style.TextAlignType) *LabelBuilder {
	return b.applyStyle(style.Style{TextAlign: style.PTextAlignType(align)})
}

func (b *LabelBuilder) VerticalAlign(align style.VerticalAlignType) *LabelBuilder {
	return b.applyStyle(style.Style{VerticalAlign: style.PVerticalAlignType(align)})
}
