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

// --- Widget Interface ---
// Widgetは全てのUI要素の基本的な振る舞いを定義するインターフェースです。
// これには位置、サイズ、スタイル、イベント処理などが含まれます。
type Widget interface {
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
	SetParent(parent Container) // 親はContainer型である必要があります
	GetParent() Container
	HitTest(x, y int) Widget
	SetHovered(hovered bool)
	IsHovered() bool
	SetVisible(visible bool)
	IsVisible() bool
	SetRelayoutBoundary(isBoundary bool)
}

// --- Container Interface ---
// Containerは子Widgetを持つことができるWidgetです。
// UIの階層構造を構築するために使用されます。
type Container interface {
	Widget // ContainerはWidgetのすべての振る舞いを継承します
	AddChild(child Widget)
	RemoveChild(child Widget)
	GetChildren() []Widget
}

// --- LayoutableWidget ---
// LayoutableWidgetは、Widgetインターフェースの基本的な実装を提供する構造体です。
// 他の具体的なウィジェット（Button, Labelなど）は、この構造体を埋め込むことで基本的な機能を利用します。
type LayoutableWidget struct {
	x, y                int
	width, height       int
	minWidth, minHeight int
	flex                int
	style               style.Style
	dirty               bool
	eventHandlers       map[event.EventType]event.EventHandler
	parent              Container // 親への参照
	isHovered           bool
	isVisible           bool // 可視性フラグ
	relayoutBoundary    bool // レイアウトの境界フラグ
}

// NewLayoutableWidget は、デフォルト値で LayoutableWidget を初期化します。
func NewLayoutableWidget() *LayoutableWidget {
	return &LayoutableWidget{
		isVisible: true, // デフォルトで表示状態にする
	}
}

// --- LayoutableWidget Methods (Interface Implementation) ---

func (w *LayoutableWidget) Update() {
	// この基本実装では何もしませんが、子を持つコンテナなどでオーバーライドされます。
	// ダーティフラグのクリアは、レイアウト計算後に移動される可能性があります。
	w.ClearDirty()
}

func (w *LayoutableWidget) Draw(screen *ebiten.Image) {
	// 非表示のウィジェットは描画しない
	if !w.isVisible {
		return
	}
	// 背景と境界線の描画（フィールドへ直接アクセス）
	drawStyledBackground(screen, w.x, w.y, w.width, w.height, w.style)
}

func (w *LayoutableWidget) SetPosition(x, y int) {
	if w.x != x || w.y != y {
		w.x = x
		w.y = y
		w.MarkDirty()
	}
}

func (w *LayoutableWidget) GetPosition() (x, y int) {
	return w.x, w.y
}

func (w *LayoutableWidget) SetSize(width, height int) {
	if w.width != width || w.height != height {
		w.width = width
		w.height = height
		w.MarkDirty()
	}
}

func (w *LayoutableWidget) GetSize() (width, height int) {
	return w.width, w.height
}

func (w *LayoutableWidget) SetMinSize(width, height int) {
	if w.minWidth != width || w.minHeight != height {
		w.minWidth = width
		w.minHeight = height
		w.MarkDirty()
	}
}

func (w *LayoutableWidget) GetMinSize() (width, height int) {
	return w.minWidth, w.minHeight
}

func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	w.MarkDirty()
}

func (w *LayoutableWidget) GetStyle() *style.Style {
	return &w.style
}

func (w *LayoutableWidget) MarkDirty() {
	if w.dirty {
		return // 既にdirtyなら何もしない
	}
	w.dirty = true

	// 親が存在し、かつ自身がレイアウト境界でない場合にのみ、親へdirtyフラグを伝播させる
	if w.parent != nil && !w.relayoutBoundary {
		w.parent.MarkDirty()
	}
}

func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	w.relayoutBoundary = isBoundary
}

func (w *LayoutableWidget) IsDirty() bool {
	return w.dirty
}

func (w *LayoutableWidget) ClearDirty() {
	w.dirty = false
}

func (w *LayoutableWidget) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if w.eventHandlers == nil {
		w.eventHandlers = make(map[event.EventType]event.EventHandler)
	}
	w.eventHandlers[eventType] = handler
}

func (w *LayoutableWidget) RemoveEventHandler(eventType event.EventType) {
	delete(w.eventHandlers, eventType)
}

func (w *LayoutableWidget) HandleEvent(event event.Event) {
	if handler, exists := w.eventHandlers[event.Type]; exists {
		handler(event)
	}
}

func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
	}
	if w.flex != flex {
		w.flex = flex
		w.MarkDirty()
	}
}

func (w *LayoutableWidget) GetFlex() int {
	return w.flex
}

func (w *LayoutableWidget) SetParent(parent Container) {
	w.parent = parent
}

func (w *LayoutableWidget) GetParent() Container {
	return w.parent
}

// HitTest は、指定された座標がウィジェットの範囲内にあるかをテストします。
// この基本実装では、子要素を持たないため、自分自身のみをチェックします。
func (w *LayoutableWidget) HitTest(x, y int) Widget {
	if !w.isVisible {
		return nil
	}
	if !(image.Point{X: x, Y: y}.In(image.Rect(w.x, w.y, w.x+w.width, w.y+w.height))) {
		return nil
	}
	return w
}

func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.isHovered != hovered {
		w.isHovered = hovered
	}
}

func (w *LayoutableWidget) IsHovered() bool {
	return w.isHovered
}

func (w *LayoutableWidget) SetVisible(visible bool) {
	if w.isVisible != visible {
		w.isVisible = visible
		w.MarkDirty()
	}
}

func (w *LayoutableWidget) IsVisible() bool {
	return w.isVisible
}

// --- Drawing Helper ---

func drawStyledBackground(dst *ebiten.Image, x, y, width, height int, s style.Style) {
	if width <= 0 || height <= 0 {
		return
	}
	if s.Background != nil && s.Background != color.Transparent {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(width), float32(height), s.Background, false)
	}
	if s.BorderColor != nil && s.BorderWidth > 0 {
		vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), s.BorderWidth, s.BorderColor, false)
	}
}

func drawAlignedText(screen *ebiten.Image, textContent string, area image.Rectangle, s style.Style) {
	if textContent == "" || s.Font == nil {
		return
	}
	contentRect := image.Rect(
		area.Min.X+s.Padding.Left,
		area.Min.Y+s.Padding.Top,
		area.Max.X-s.Padding.Right,
		area.Max.Y-s.Padding.Bottom,
	)
	if contentRect.Dx() <= 0 || contentRect.Dy() <= 0 {
		return
	}
	bounds := text.BoundString(s.Font, textContent)
	textX := contentRect.Min.X + (contentRect.Dx()-bounds.Dx())/2
	textY := contentRect.Min.Y + (contentRect.Dy()+bounds.Dy())/2
	var textColor color.Color = color.Black
	if s.TextColor != nil {
		textColor = s.TextColor
	}
	text.Draw(screen, textContent, s.Font, textX, textY, textColor)
}

// --- TextWidget ---
// TextWidgetは、テキスト表示に関連する共通の機能（テキスト内容、スタイル、最小サイズ計算）を提供します。
// ButtonやLabelなど、テキストを持つウィジェットはこれを埋め込みます。
type TextWidget struct {
	*LayoutableWidget
	text string
}

// NewTextWidget は新しいTextWidgetを生成します。
func NewTextWidget(text string) *TextWidget {
	return &TextWidget{
		LayoutableWidget: NewLayoutableWidget(),
		text:             text,
	}
}

// Text はウィジェットのテキストを取得します。
func (t *TextWidget) Text() string {
	return t.text
}

// SetText はウィジェットのテキストを設定し、ダーティフラグを立てます。
func (t *TextWidget) SetText(text string) {
	if t.text != text {
		t.text = text
		t.MarkDirty()
	}
}

// Draw はTextWidgetを描画します。LayoutableWidgetのDrawをオーバーライドしてテキストを追加描画します。
func (t *TextWidget) Draw(screen *ebiten.Image) {
	if !t.isVisible {
		return
	}
	// 背景描画は基本のDrawを呼び出す
	t.LayoutableWidget.Draw(screen)
	// テキストを描画（フィールドへ直接アクセス）
	drawAlignedText(screen, t.text, image.Rect(t.x, t.y, t.x+t.width, t.y+t.height), t.style)
}

// calculateMinSize は、現在のテキストとスタイルに基づいて最小サイズを計算します。
func (t *TextWidget) calculateMinSize() (int, int) {
	style := t.GetStyle()
	if t.text != "" && style.Font != nil {
		bounds := text.BoundString(style.Font, t.text)
		minWidth := bounds.Dx() + style.Padding.Left + style.Padding.Right
		metrics := style.Font.Metrics()
		minHeight := (metrics.Ascent + metrics.Descent).Ceil() + style.Padding.Top + style.Padding.Bottom
		return minWidth, minHeight
	}
	return 0, 0
}

// --- Button component ---
type Button struct {
	*TextWidget
	hoverStyle *style.Style
}

// Draw はButtonを描画します。TextWidgetのDrawをオーバーライドしてホバー効果を追加します。
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.isVisible {
		return
	}
	// 現在適用すべきスタイルを選択（通常時 or ホバー時）
	currentStyle := &b.style
	if b.isHovered && b.hoverStyle != nil {
		currentStyle = b.hoverStyle
	}
	// 背景と境界線を描画（フィールドへ直接アクセス）
	drawStyledBackground(screen, b.x, b.y, b.width, b.height, *currentStyle)
	// テキストを描画（フィールドへ直接アクセス）
	drawAlignedText(screen, b.text, image.Rect(b.x, b.y, b.x+b.width, b.y+b.height), *currentStyle)
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
		TextWidget: NewTextWidget(""),
	}
	button.width = 100
	button.height = 40
	button.SetStyle(defaultStyle)

	return &ButtonBuilder{
		button: button,
	}
}

func (b *ButtonBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.button.calculateMinSize()
	b.button.SetMinSize(minWidth, minHeight)
}

func (b *ButtonBuilder) CalculateMinSize() *ButtonBuilder {
	b.calculateMinSizeInternal()
	return b
}

func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.button.SetText(text)
	b.calculateMinSizeInternal()
	return b
}

func (b *ButtonBuilder) Size(width, height int) *ButtonBuilder {
	b.button.SetSize(width, height)
	return b
}

func (b *ButtonBuilder) OnClick(onClick func()) *ButtonBuilder {
	if onClick != nil {
		b.button.AddEventHandler(event.EventClick, func(e event.Event) {
			onClick()
		})
	}
	return b
}

func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	existingStyle := b.button.GetStyle()
	b.button.SetStyle(style.Merge(*existingStyle, s))
	b.calculateMinSizeInternal()
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
	b.calculateMinSizeInternal()
	return b.button, nil
}

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
type Label struct {
	*TextWidget
}

// Drawは埋め込まれたTextWidgetのDrawメソッドをそのまま利用します。
// func (l *Label) Draw(screen *ebiten.Image) {
// 	l.TextWidget.Draw(screen)
// }

// --- LabelBuilder ---
type LabelBuilder struct {
	label  *Label
	errors []error
}

func NewLabelBuilder() *LabelBuilder {
	label := &Label{
		TextWidget: NewTextWidget(""),
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

func (b *LabelBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.label.calculateMinSize()
	b.label.SetMinSize(minWidth, minHeight)
}

func (b *LabelBuilder) CalculateMinSize() *LabelBuilder {
	b.calculateMinSizeInternal()
	return b
}

func (b *LabelBuilder) Text(text string) *LabelBuilder {
	b.label.SetText(text)
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
