package component

import (
	"fmt"
	"furoshiki/event"
	"furoshiki/style"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	relayoutDirty       bool // 再レイアウトが必要かどうかのフラグ
	eventHandlers       map[event.EventType]event.EventHandler
	parent              Container // 親への参照
	isHovered           bool
	isVisible           bool // 可視性フラグ
	relayoutBoundary    bool // レイアウトの境界フラグ
}

// NewLayoutableWidget は、デフォルト値で LayoutableWidget を初期化します。
func NewLayoutableWidget() *LayoutableWidget {
	return &LayoutableWidget{
		isVisible:     true, // デフォルトで表示状態にする
		eventHandlers: make(map[event.EventType]event.EventHandler),
	}
}

// --- LayoutableWidget Methods (Interface Implementation) ---

func (w *LayoutableWidget) Update() {
	// この基本実装は、具象ウィジェット（Button, Labelなど）で利用されます。
	// Container型は自身のUpdateメソッドでこれをオーバーライドします。
	// ダーティフラグのクリアは、レイアウト計算を行うContainerの責務であるため、
	// ここでは何もしません。
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
		w.MarkDirty(false) // 位置変更は再描画のみが必要
	}
}

func (w *LayoutableWidget) GetPosition() (x, y int) {
	return w.x, w.y
}

func (w *LayoutableWidget) SetSize(width, height int) {
	// サイズが負の値の場合は処理しない
	if width < 0 || height < 0 {
		return
	}

	if w.width != width || w.height != height {
		w.width = width
		w.height = height
		w.MarkDirty(true) // サイズ変更は再レイアウトが必要
	}
}

func (w *LayoutableWidget) GetSize() (width, height int) {
	return w.width, w.height
}

func (w *LayoutableWidget) SetMinSize(width, height int) {
	// 最小サイズが負の値の場合は処理しない
	if width < 0 || height < 0 {
		return
	}

	if w.minWidth != width || w.minHeight != height {
		w.minWidth = width
		w.minHeight = height
		w.MarkDirty(true) // 最小サイズ変更は再レイアウトが必要
	}
}

func (w *LayoutableWidget) GetMinSize() (width, height int) {
	return w.minWidth, w.minHeight
}

func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更は必ずしも再レイアウトを必要としないかもしれないが、
	// Paddingなどが変わる可能性があるため、安全策としてtrueにする
	w.MarkDirty(true)
}

func (w *LayoutableWidget) GetStyle() *style.Style {
	return &w.style
}

func (w *LayoutableWidget) MarkDirty(relayout bool) {
	// すでにダーティで、かつ再レイアウトフラグが既に立っている場合は何もしない
	if w.dirty && (!relayout || w.relayoutDirty) {
		return
	}

	w.dirty = true
	if relayout {
		w.relayoutDirty = true
	}

	// 親が存在し、かつ自身がレイアウト境界でなく、再レイアウトが必要な場合のみ伝播
	if w.parent != nil && !w.relayoutBoundary && relayout {
		// インターフェース経由で直接呼び出す - 型アサーションは不要
		w.parent.MarkDirty(true)
	}
}

func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	if w.relayoutBoundary != isBoundary {
		w.relayoutBoundary = isBoundary
		// 境界設定が変更された場合は、再レイアウトが必要
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) IsDirty() bool {
	return w.dirty
}

func (w *LayoutableWidget) ClearDirty() {
	w.dirty = false
	w.relayoutDirty = false
}

func (w *LayoutableWidget) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if w.eventHandlers == nil {
		w.eventHandlers = make(map[event.EventType]event.EventHandler)
	}
	w.eventHandlers[eventType] = handler
}

func (w *LayoutableWidget) RemoveEventHandler(eventType event.EventType) {
	if w.eventHandlers != nil {
		delete(w.eventHandlers, eventType)
	}
}

func (w *LayoutableWidget) HandleEvent(event event.Event) {
	if handler, exists := w.eventHandlers[event.Type]; exists {
		// イベントハンドラの実行中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic in event handler: %v\n", r)
			}
		}()
		handler(event)
	}
}

func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
	}
	if w.flex != flex {
		w.flex = flex
		w.MarkDirty(true)
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
	// 境界チェックをより明確に実装
	rect := image.Rect(w.x, w.y, w.x+w.width, w.y+w.height)
	if rect.Empty() {
		return nil
	}
	if !(image.Point{X: x, Y: y}.In(rect)) {
		return nil
	}
	return w
}

func (w *LayoutableWidget) SetHovered(hovered bool) {
	if w.isHovered != hovered {
		w.isHovered = hovered
		// 再描画のみを要求するため、relayoutはfalse
		w.MarkDirty(false)
	}
}

func (w *LayoutableWidget) IsHovered() bool {
	return w.isHovered
}

func (w *LayoutableWidget) SetVisible(visible bool) {
	if w.isVisible != visible {
		w.isVisible = visible
		// 表示状態の変更はレイアウトに影響する可能性があるためtrue
		w.MarkDirty(true)
	}
}

func (w *LayoutableWidget) IsVisible() bool {
	return w.isVisible
}

// Cleanup は、コンポーネントが不要になったときにリソースを解放するためのメソッドです。
func (w *LayoutableWidget) Cleanup() {
	// イベントハンドラをクリア
	w.eventHandlers = nil
	// 親からの参照を解除
	w.parent = nil
}
