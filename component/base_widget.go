package component

import (
	"furoshiki/event"
	"furoshiki/style"
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