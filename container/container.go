package container

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/layout"
	"image"
	"log"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
)

// Scroller is an interface used by Container to get scroll offsets for clipping.
type Scroller interface {
	GetScrollOffset() (x, y int)
}

// Container is a component that holds child Widgets and manages their layout.
type Container struct {
	*component.WidgetCore
	*component.Appearance
	*component.Interaction

	// UPDATE: childrenフィールドを削除。子要素の管理はWidgetCore.Nodeに一本化されました。
	// これにより、データの一貫性が保証され、AddChildのバグが解消されます。

	layout         layout.Layout
	warned         bool
	clipsChildren  bool
	offscreenImage *ebiten.Image
}

// --- Interface implementation verification ---
var _ component.StandardWidget = (*Container)(nil)
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)
var _ event.EventTarget = (*Container)(nil)

// NewContainer creates a new Container instance without using a builder.
func NewContainer() (*Container, error) {
	c := &Container{
		// UPDATE: childrenフィールドがなくなったため、初期化も不要
	}

	c.WidgetCore = component.NewWidgetCore(c)
	c.Appearance = component.NewAppearance(c)
	c.Interaction = component.NewInteraction(c)

	c.layout = &layout.FlexLayout{} // Default to FlexLayout
	return c, nil
}

// --- Method implementations ---

func (c *Container) SetClipsChildren(clips bool) {
	if c.clipsChildren != clips {
		c.clipsChildren = clips
		c.MarkDirty(false)
	}
}

func (c *Container) GetLayout() layout.Layout {
	return c.layout
}

func (c *Container) SetLayout(layout layout.Layout) {
	c.layout = layout
	c.MarkDirty(true)
}

func (c *Container) Update() {
	if c.GetParent() == nil && !c.HasBeenLaidOut() {
		c.SetLaidOut(true)
		c.MarkDirty(true)
	}

	if !c.IsVisible() {
		return
	}

	c.checkSizeWarning()

	if c.NeedsRelayout() {
		if c.layout != nil {
			if err := c.layout.Layout(c); err != nil {
				log.Printf("Error during layout calculation: %v\n%s", err, debug.Stack())
			}
		}
	}

	// UPDATE: Nodeから子を取得するように変更
	for _, child := range c.GetChildren() {
		child.Update()
	}

	if c.IsDirty() {
		c.ClearDirty()
	}
}

func (c *Container) checkSizeWarning() {
	if c.warned {
		return
	}
	if c.GetFlex() == 0 {
		width, height := c.GetSize()
		if width == 0 && height == 0 && c.GetParent() == nil {
			log.Printf("Warning: Root container created with no flex and zero size. It may not be visible.\n")
			c.warned = true
		}
	}
}

func (c *Container) Draw(info component.DrawInfo) {
	if !c.IsVisible() || !c.HasBeenLaidOut() {
		return
	}

	if c.clipsChildren {
		c.drawWithClipping(info)
	} else {
		c.drawWithoutClipping(info)
	}
}

func (c *Container) drawWithoutClipping(info component.DrawInfo) {
	x, y := c.GetPosition()
	width, height := c.GetSize()
	finalX := x + info.OffsetX
	finalY := y + info.OffsetY
	component.DrawStyledBackground(info.Screen, finalX, finalY, width, height, c.ReadOnlyStyle())

	for _, child := range c.GetChildren() {
		child.Draw(info)
	}
}

func (c *Container) drawWithClipping(info component.DrawInfo) {
	containerX, containerY := c.GetPosition()
	containerWidth, containerHeight := c.GetSize()

	if containerWidth <= 0 || containerHeight <= 0 {
		return
	}

	if c.offscreenImage == nil || c.offscreenImage.Bounds().Dx() != containerWidth || c.offscreenImage.Bounds().Dy() != containerHeight {
		if c.offscreenImage != nil {
			c.offscreenImage.Deallocate()
		}
		c.offscreenImage = ebiten.NewImage(containerWidth, containerHeight)
	}
	c.offscreenImage.Clear()

	component.DrawStyledBackground(c.offscreenImage, 0, 0, containerWidth, containerHeight, c.ReadOnlyStyle())

	var scrollOffsetX, scrollOffsetY int
	if scroller, ok := any(c).(Scroller); ok {
		scrollOffsetX, scrollOffsetY = scroller.GetScrollOffset()
	}

	childDrawInfo := component.DrawInfo{
		Screen:  c.offscreenImage,
		OffsetX: -(containerX - scrollOffsetX),
		OffsetY: -(containerY - scrollOffsetY),
	}

	for _, child := range c.GetChildren() {
		child.Draw(childDrawInfo)
	}

	opts := &ebiten.DrawImageOptions{}
	finalX := float64(containerX + info.OffsetX)
	finalY := float64(containerY + info.OffsetY)
	opts.GeoM.Translate(finalX, finalY)
	info.Screen.DrawImage(c.offscreenImage, opts)
}

func (c *Container) HitTest(x, y int) component.Widget {
	if !c.IsVisible() || c.IsDisabled() {
		return nil
	}

	children := c.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		isVisible := true
		if is, ok := child.(component.InteractiveState); ok {
			isVisible = is.IsVisible()
		}
		if !isVisible {
			continue
		}
		if target := child.HitTest(x, y); target != nil {
			return target
		}
	}

	wx, wy := c.GetPosition()
	wwidth, wheight := c.GetSize()
	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if !rect.Empty() && (image.Point{X: x, Y: y}.In(rect)) {
		return c
	}

	return nil
}

func (c *Container) Cleanup() {
	// 1. 再帰的にすべての子のCleanupを呼び出す
	for _, child := range c.GetChildren() {
		child.Cleanup()
	}
	// 2. Nodeの子リストをクリアする
	c.GetNode().ClearChildren()

	// 3. このコンテナ自身のリソースを解放
	if c.offscreenImage != nil {
		c.offscreenImage.Deallocate()
		c.offscreenImage = nil
	}

	// 4. 親への参照を断つ
	c.SetParent(nil)
}

// UPDATE: 子要素の管理をWidgetCore.Nodeに委譲
func (c *Container) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	// NodeのAddChildは、古い親からのデタッチ処理も自動的に行うため、
	// ここでのロジックがシンプルかつ安全になります。
	c.WidgetCore.Node.AddChild(child)
	c.MarkDirty(true)
}

// UPDATE: 子要素の管理をWidgetCore.Nodeに委譲
func (c *Container) RemoveChild(child component.Widget) {
	// NodeのRemoveChildは子の親ポインタをnilにするだけです。
	// ウィジェットの完全な破棄（リソース解放など）を行うのはコンテナの責務なので、
	// ここでCleanupを呼び出します。
	c.WidgetCore.Node.RemoveChild(child)
	child.Cleanup()
	c.MarkDirty(true)
}

// UPDATE: 子要素の管理をWidgetCore.Nodeに委譲
func (c *Container) GetChildren() []component.Widget {
	// Nodeは[]NodeOwnerを返すため、Containerインターフェースが要求する[]Widgetに変換します。
	nodeOwners := c.WidgetCore.Node.GetChildren()
	widgets := make([]component.Widget, len(nodeOwners))
	for i, owner := range nodeOwners {
		// AddChildでWidget型しか受け付けないため、この型アサーションは常に成功するはずです。
		widgets[i] = owner.(component.Widget)
	}
	return widgets
}

func (c *Container) GetPadding() layout.Insets {
	style := c.ReadOnlyStyle()
	if style.Padding != nil {
		return layout.Insets{
			Top:    style.Padding.Top,
			Right:  style.Padding.Right,
			Bottom: style.Padding.Bottom,
			Left:   style.Padding.Left,
		}
	}
	return layout.Insets{}
}

func (c *Container) HandleEvent(e *event.Event) {
	c.Interaction.TriggerHandlers(e)

	if e != nil && !e.Handled && c.GetParent() != nil {
		if processor, ok := c.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}