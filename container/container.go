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
	// Component-based fields
	*component.Node
	*component.Transform
	*component.LayoutProperties
	*component.Appearance
	*component.Interaction
	*component.Visibility
	*component.Dirty

	// Container-specific fields
	children       []component.Widget
	layout         layout.Layout
	warned         bool
	clipsChildren  bool
	offscreenImage *ebiten.Image
	// UPDATE: hasBeenLaidOutフィールドはVisibilityコンポーネントに統合されたため削除されました。
	// hasBeenLaidOut bool
	minSize component.Size
}

// --- Interface implementation verification ---
var _ component.Widget = (*Container)(nil)
var _ component.Container = (*Container)(nil)
var _ layout.Container = (*Container)(nil)
var _ event.EventTarget = (*Container)(nil)
var _ component.EventProcessor = (*Container)(nil)

// NewContainer creates a new Container instance without using a builder.
func NewContainer() (*Container, error) {
	c := &Container{
		children: make([]component.Widget, 0),
	}

	// Initialize components
	c.Node = component.NewNode(c)
	c.Transform = component.NewTransform()
	c.LayoutProperties = component.NewLayoutProperties()
	c.Appearance = component.NewAppearance(c)
	c.Interaction = component.NewInteraction(c)
	c.Visibility = component.NewVisibility(c)
	c.Dirty = component.NewDirty()

	c.layout = &layout.FlexLayout{} // Default to FlexLayout
	return c, nil
}

// --- Method implementations ---

func (c *Container) GetNode() *component.Node                   { return c.Node }
func (c *Container) GetLayoutProperties() *component.LayoutProperties { return c.LayoutProperties }

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
	// UPDATE: hasBeenLaidOutのチェックをVisibilityコンポーネントのHasBeenLaidOut()に置き換えました。
	if c.GetParent() == nil && !c.HasBeenLaidOut() {
		// UPDATE: 初回レイアウトをトリガーするためにSetLaidOutを呼び出します。
		// これにより、Visibilityコンポーネント内でダーティフラグが設定され、初回描画が保証されます。
		c.SetLaidOut(true)
		c.MarkDirty(true) // Mark dirty to trigger initial layout
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

	for _, child := range c.children {
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
	// UPDATE: hasBeenLaidOutのチェックをVisibilityコンポーネントのHasBeenLaidOut()に置き換えました。
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

	for _, child := range c.children {
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

	for _, child := range c.children {
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

	for i := len(c.children) - 1; i >= 0; i-- {
		child := c.children[i]
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
	for _, child := range c.children {
		child.Cleanup()
	}
	c.children = nil

	if c.offscreenImage != nil {
		c.offscreenImage.Deallocate()
		c.offscreenImage = nil
	}

	c.SetParent(nil)
}

func (c *Container) detachChild(child component.Widget) bool {
	if child == nil {
		return false
	}
	for i, currentChild := range c.children {
		if currentChild == child {
			c.children = append(c.children[:i], c.children[i+1:]...)
			child.SetParent(nil)
			return true
		}
	}
	return false
}

func (c *Container) AddChild(child component.Widget) {
	if child == nil {
		return
	}
	if oldParent := child.GetParent(); oldParent != nil {
		if container, ok := oldParent.(*Container); ok {
			container.detachChild(child)
		}
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.MarkDirty(true)
}

func (c *Container) RemoveChild(child component.Widget) {
	if c.detachChild(child) {
		child.Cleanup()
		c.MarkDirty(true)
	}
}

func (c *Container) GetChildren() []component.Widget {
	return c.children
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

func (c *Container) MarkDirty(relayout bool) {
	c.Dirty.MarkDirty(relayout)
	if relayout && !c.IsLayoutBoundary() {
		if parent := c.GetParent(); parent != nil {
			if dm, ok := parent.(component.DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

func (c *Container) SetPosition(x, y int) {
	// UPDATE: レイアウト済み状態の管理をVisibilityコンポーネントに委譲します。
	if !c.HasBeenLaidOut() {
		c.SetLaidOut(true)
	}
	if currX, currY := c.GetPosition(); currX != x || currY != y {
		c.Transform.SetPosition(x, y)
		c.MarkDirty(false)
	}
}

func (c *Container) SetSize(width, height int) {
	if w, h := c.GetSize(); w != width || h != height {
		c.Transform.SetSize(width, height)
		c.MarkDirty(true)
	}
}

func (c *Container) SetMinSize(width, height int) {
	c.minSize.Width = width
	c.minSize.Height = height
	c.MarkDirty(true)
}

func (c *Container) GetMinSize() (int, int) {
	return c.minSize.Width, c.minSize.Height
}

func (c *Container) HandleEvent(e *event.Event) {
	// UPDATE: イベントハンドラの安全な実行をInteractionコンポーネントに委譲
	c.Interaction.TriggerHandlers(e)

	// 親ウィジェットへのイベント伝播
	if e != nil && !e.Handled && c.GetParent() != nil {
		if processor, ok := c.GetParent().(component.EventProcessor); ok {
			processor.HandleEvent(e)
		}
	}
}