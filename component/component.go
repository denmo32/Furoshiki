package component

import (
	"fmt"
	"furoshiki/event"
	"furoshiki/style"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Component Interface ---
// Componentは全てのUIコンポーネントの核となるインターフェースです。
type Component interface {
	Update()
	Draw(screen *ebiten.Image)
	SetPosition(x, y int)
	GetPosition() (x, y int)
	SetSize(width, height int)
	GetSize() (width, height int)
	SetMinSize(width, height int)
	GetMinSize() (width, height int)
	SetStyle(style style.Style)
	GetStyle() *style.Style
	MarkDirty()
	IsDirty() bool
	ClearDirty()
	AddEventHandler(eventType event.EventType, handler event.EventHandler)
	RemoveEventHandler(eventType event.EventType)
	HandleEvent(event event.Event)
	SetFlex(flex int)
	GetFlex() int
	AddChild(child Component)
	RemoveChild(child Component)
	GetChildren() []Component
	SetParent(parent Component)
	GetParent() Component
	HitTest(x, y int) Component
	SetHovered(hovered bool)
	IsHovered() bool
	SetVisible(visible bool) // 可視性を設定するメソッドを追加
	IsVisible() bool         // 可視性を取得するメソッドを追加
	SetRelayoutBoundary(isBoundary bool)
}

// --- LayoutableComponent ---
// LayoutableComponentはレイアウト可能なコンポーネントの基底構造体です。
type LayoutableComponent struct {
	x, y                int
	width, height       int
	minWidth, minHeight int
	flex                int
	style               style.Style
	dirty               bool
	eventHandlers       map[event.EventType]event.EventHandler
	children            []Component
	parent              Component
	isHovered           bool
	isVisible           bool // 可視性フラグを追加
	relayoutBoundary    bool // レイアウトの境界フラグ
}

// NewLayoutableComponent は、デフォルト値で LayoutableComponent を初期化します。
// 特に isVisible を true に設定するために重要です。
func NewLayoutableComponent() *LayoutableComponent {
	return &LayoutableComponent{
		isVisible: true,
	}
}

// --- LayoutableComponent Methods (Interface Implementation) ---

func (c *LayoutableComponent) Update() {
	// 非表示のコンポーネントとその子孫は更新しない
	if !c.isVisible {
		return
	}
	for _, child := range c.children {
		child.Update()
	}
	c.ClearDirty()
}

func (c *LayoutableComponent) Draw(screen *ebiten.Image) {
	// 非表示のコンポーネントとその子孫は描画しない
	if !c.isVisible {
		return
	}
	style := c.GetStyle()
	x, y := c.GetPosition()
	width, height := c.GetSize()
	drawStyledBackground(screen, x, y, width, height, *style)
	for _, child := range c.children {
		child.Draw(screen)
	}
}

func (c *LayoutableComponent) SetPosition(x, y int) {
	if c.x != x || c.y != y {
		c.x = x
		c.y = y
		c.MarkDirty()
	}
}

func (c *LayoutableComponent) GetPosition() (x, y int) {
	return c.x, c.y
}

func (c *LayoutableComponent) SetSize(width, height int) {
	if c.width != width || c.height != height {
		c.width = width
		c.height = height
		c.MarkDirty()
	}
}

func (c *LayoutableComponent) GetSize() (width, height int) {
	return c.width, c.height
}

func (c *LayoutableComponent) SetMinSize(width, height int) {
	if c.minWidth != width || c.minHeight != height {
		c.minWidth = width
		c.minHeight = height
		c.MarkDirty()
	}
}

func (c *LayoutableComponent) GetMinSize() (width, height int) {
	return c.minWidth, c.minHeight
}

func (c *LayoutableComponent) SetStyle(style style.Style) {
	c.style = style
	c.MarkDirty()
}

func (c *LayoutableComponent) GetStyle() *style.Style {
	return &c.style
}

func (c *LayoutableComponent) MarkDirty() {
	if c.dirty {
		return // 既にdirtyなら何もしない
	}
	c.dirty = true

	// 親が存在し、かつ自身がレイアウト境界でない場合にのみ、親へdirtyフラグを伝播させる
	if c.parent != nil && !c.relayoutBoundary {
		c.parent.MarkDirty()
	}
}

// SetRelayoutBoundary は、このコンポーネントがレイアウトの境界であるかを設定します。
// 境界に設定されたコンポーネントは、自身の変更を親に伝播させなくなり、
// レイアウトの再計算範囲を限定することができます。
func (c *LayoutableComponent) SetRelayoutBoundary(isBoundary bool) {
	c.relayoutBoundary = isBoundary
}

func (c *LayoutableComponent) IsDirty() bool {
	return c.dirty
}

func (c *LayoutableComponent) ClearDirty() {
	c.dirty = false
}

func (c *LayoutableComponent) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if c.eventHandlers == nil {
		c.eventHandlers = make(map[event.EventType]event.EventHandler)
	}
	c.eventHandlers[eventType] = handler
}

func (c *LayoutableComponent) RemoveEventHandler(eventType event.EventType) {
	delete(c.eventHandlers, eventType)
}

func (c *LayoutableComponent) HandleEvent(event event.Event) {
	if handler, exists := c.eventHandlers[event.Type]; exists {
		handler(event)
	}
}

func (c *LayoutableComponent) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
	}
	if c.flex != flex {
		c.flex = flex
		c.MarkDirty()
	}
}

func (c *LayoutableComponent) GetFlex() int {
	return c.flex
}

func (c *LayoutableComponent) AddChild(child Component) {
	if child == nil {
		return
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.MarkDirty()
}

func (c *LayoutableComponent) RemoveChild(child Component) {
	if child == nil {
		return
	}
	for i, currentChild := range c.children {
		if currentChild == child {
			c.children = append(c.children[:i], c.children[i+1:]...)
			child.SetParent(nil)
			c.MarkDirty()
			return
		}
	}
}

func (c *LayoutableComponent) GetChildren() []Component {
	return c.children
}

func (c *LayoutableComponent) SetParent(parent Component) {
	c.parent = parent
}

func (c *LayoutableComponent) GetParent() Component {
	return c.parent
}

func (c *LayoutableComponent) HitTest(x, y int) Component {
	// 非表示のコンポーネントはヒットしない
	if !c.isVisible {
		return nil
	}
	cx, cy := c.GetPosition()
	cw, ch := c.GetSize()
	if !(image.Point{X: x, Y: y}.In(image.Rect(cx, cy, cx+cw, cy+ch))) {
		return nil
	}
	for i := len(c.children) - 1; i >= 0; i-- {
		child := c.children[i]
		if target := child.HitTest(x, y); target != nil {
			return target
		}
	}
	return c
}

func (c *LayoutableComponent) SetHovered(hovered bool) {
	if c.isHovered != hovered {
		c.isHovered = hovered
		// Hover状態の変更は通常レイアウトに影響しないため、MarkDirty()は呼び出さない。
		// 再描画はEbitengineの毎フレームのDrawループに任せる。
	}
}

func (c *LayoutableComponent) IsHovered() bool {
	return c.isHovered
}

// SetVisible はコンポーネントの可視性を設定します。
func (c *LayoutableComponent) SetVisible(visible bool) {
	if c.isVisible != visible {
		c.isVisible = visible
		// 可視性の変更はレイアウトに直接影響するため、再計算を要求する
		c.MarkDirty()
	}
}

// IsVisible はコンポーネントが可視状態であるかを返します。
func (c *LayoutableComponent) IsVisible() bool {
	return c.isVisible
}

// --- Drawing Helper ---

// drawStyledBackground は、指定されたスタイルでコンポーネントの背景と境界線を描画します。
func drawStyledBackground(dst *ebiten.Image, x, y, width, height int, s style.Style) {
	if width <= 0 || height <= 0 {
		return
	}
	// 背景色の描画
	if s.Background != nil && s.Background != color.Transparent {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), s.Background, false)
	}
	// 境界線の描画
	if s.BorderColor != nil && s.BorderWidth > 0 {
		vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), s.BorderWidth, s.BorderColor, false)
	}
}

// drawAlignedText は、指定された領域内にテキストを中央揃えで描画するヘルパー関数です。
// Button や Label など、テキストを持つ複数のコンポーネントから利用されることを想定しています。
func drawAlignedText(screen *ebiten.Image, textContent string, area image.Rectangle, s style.Style) {
	// テキストやフォントが空の場合は描画しない
	if textContent == "" || s.Font == nil {
		return
	}

	// パディングを適用した描画領域を計算
	contentRect := image.Rect(
		area.Min.X+s.Padding.Left,
		area.Min.Y+s.Padding.Top,
		area.Max.X-s.Padding.Right,
		area.Max.Y-s.Padding.Bottom,
	)
	contentWidth := contentRect.Dx()
	contentHeight := contentRect.Dy()

	// 描画領域がなければ何もしない
	if contentWidth <= 0 || contentHeight <= 0 {
		return
	}

	// テキストの描画サイズを計算
	bounds := text.BoundString(s.Font, textContent)
	// テキストを中央に配置するための座標を計算
	textX := contentRect.Min.X + (contentWidth-bounds.Dx())/2
	textY := contentRect.Min.Y + (contentHeight+bounds.Dy())/2

	// テキストカラーを設定（指定がなければ黒）
	var textColor color.Color = color.Black
	if s.TextColor != nil {
		textColor = s.TextColor
	}

	// テキストを描画
	text.Draw(screen, textContent, s.Font, textX, textY, textColor)
}

// --- Button component ---
type Button struct {
	*LayoutableComponent
	text string
	// onClick フィールドを削除し、イベントハンドラに一本化
	hoverStyle *style.Style
}

func (b *Button) Update() {
	b.LayoutableComponent.Update()
}

func (b *Button) Draw(screen *ebiten.Image) {
	// isVisibleチェックは埋め込み元のDrawで既に行われている
	width, height := b.GetSize()
	x, y := b.GetPosition()
	if width <= 0 || height <= 0 {
		return
	}
	currentStyle := b.GetStyle()
	if b.IsHovered() && b.hoverStyle != nil {
		currentStyle = b.hoverStyle
	}
	// 背景と境界線を描画
	drawStyledBackground(screen, x, y, width, height, *currentStyle)
	// 新しいヘルパー関数を使用してテキストを描画
	drawAlignedText(screen, b.text, image.Rect(x, y, x+width, y+height), *currentStyle)
}

// --- ButtonBuilder ---
type ButtonBuilder struct {
	button *Button
	errors []error
}

func NewButtonBuilder() *ButtonBuilder {
	defaultStyle := style.Style{
		Background:  color.RGBA{R: 220, G: 220, B: 220, A: 255},
		TextColor:   color.Black,
		BorderColor: color.Gray{Y: 150},
		BorderWidth: 1,
		Padding:     style.Insets{Top: 5, Right: 10, Bottom: 5, Left: 10},
	}
	button := &Button{
		LayoutableComponent: NewLayoutableComponent(),
	}
	button.width = 100
	button.height = 40
	button.SetStyle(defaultStyle)

	return &ButtonBuilder{
		button: button,
	}
}

// calculateMinSizeInternal は、現在のテキストとスタイルに基づいて最小サイズを計算し、設定します。
// ユーザーが明示的に呼び出す必要がないように、内部的に使用されます。
func (b *ButtonBuilder) calculateMinSizeInternal() {
	style := b.button.GetStyle()
	if b.button.text != "" && style.Font != nil {
		bounds := text.BoundString(style.Font, b.button.text)
		minWidth := bounds.Dx() + style.Padding.Left + style.Padding.Right
		metrics := style.Font.Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + style.Padding.Top + style.Padding.Bottom
		b.button.SetMinSize(minWidth, minHeight)
	}
}

// CalculateMinSize は、手動で最小サイズの再計算をトリガーしたい場合に使用します。
// 通常はText()やStyle()の設定時に自動で計算されます。
func (b *ButtonBuilder) CalculateMinSize() *ButtonBuilder {
	b.calculateMinSizeInternal()
	return b
}

func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.button.text = text
	b.calculateMinSizeInternal() // テキストが変更されたら自動で最小サイズを再計算
	return b
}

func (b *ButtonBuilder) Size(width, height int) *ButtonBuilder {
	b.button.SetSize(width, height)
	return b
}

func (b *ButtonBuilder) OnClick(onClick func()) *ButtonBuilder {
	if onClick != nil {
		b.button.AddEventHandler(event.EventClick, func(e event.Event) {
			// e.MouseButton == ebiten.MouseButtonLeft などをチェックすることも可能
			onClick()
		})
	}
	return b
}

func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	existingStyle := b.button.GetStyle()
	b.button.SetStyle(style.Merge(*existingStyle, s))
	b.calculateMinSizeInternal() // スタイル（フォントやパディング）が変更されたら再計算
	return b
}

func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	mergedHoverStyle := style.Merge(*b.button.GetStyle(), s)
	b.button.hoverStyle = &mergedHoverStyle
	return b
}

func (b *ButtonBuilder) Flex(flex int) *ButtonBuilder {
	b.button.SetFlex(flex)
	return b
}

func (b *ButtonBuilder) Build() (*Button, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("button build errors: %v", b.errors)
	}
	// Build時にも最終的な最小サイズを計算し、設定漏れを防ぐ
	b.calculateMinSizeInternal()
	return b.button, nil
}

// --- Label component ---
type Label struct {
	*LayoutableComponent
	text string
}

func (l *Label) Update() {
	l.LayoutableComponent.Update()
}

func (l *Label) Draw(screen *ebiten.Image) {
	// isVisibleチェックは埋め込み元のDrawで既に行われている
	style := l.GetStyle()
	width, height := l.GetSize()
	x, y := l.GetPosition()

	// 背景と境界線を描画
	drawStyledBackground(screen, x, y, width, height, *style)
	// 新しいヘルパー関数を使用してテキストを描画
	drawAlignedText(screen, l.text, image.Rect(x, y, x+width, y+height), *style)
}

// --- LabelBuilder ---
type LabelBuilder struct {
	label  *Label
	errors []error
}

func NewLabelBuilder() *LabelBuilder {
	label := &Label{
		LayoutableComponent: NewLayoutableComponent(),
	}
	label.width = 100
	label.height = 30
	defaultStyle := style.Style{
		Background: color.Transparent,
		TextColor:  color.Black,
		Padding:    style.Insets{Top: 2, Right: 5, Bottom: 2, Left: 5},
	}
	label.SetStyle(defaultStyle)

	return &LabelBuilder{
		label: label,
	}
}

// calculateMinSizeInternal は、現在のテキストとスタイルに基づいて最小サイズを計算し、設定します。
func (b *LabelBuilder) calculateMinSizeInternal() {
	style := b.label.GetStyle()
	if b.label.text != "" && style.Font != nil {
		bounds := text.BoundString(style.Font, b.label.text)
		minWidth := bounds.Dx() + style.Padding.Left + style.Padding.Right
		metrics := style.Font.Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + style.Padding.Top + style.Padding.Bottom
		b.label.SetMinSize(minWidth, minHeight)
	}
}

// CalculateMinSize は、手動で最小サイズの再計算をトリガーしたい場合に使用します。
func (b *LabelBuilder) CalculateMinSize() *LabelBuilder {
	b.calculateMinSizeInternal()
	return b
}

func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.label.text = text
	b.calculateMinSizeInternal()
	return b
}

func (b *LabelBuilder) Size(width, height int) *LabelBuilder {
	b.label.SetSize(width, height)
	return b
}

func (b *LabelBuilder) Style(s style.Style) *LabelBuilder {
	existingStyle := b.label.GetStyle()
	b.label.SetStyle(style.Merge(*existingStyle, s))
	b.calculateMinSizeInternal()
	return b
}

func (b *LabelBuilder) Flex(flex int) *LabelBuilder {
	b.label.SetFlex(flex)
	return b
}

func (b *LabelBuilder) Build() (*Label, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("label build errors: %v", b.errors)
	}
	b.calculateMinSizeInternal()
	return b.label, nil
}
