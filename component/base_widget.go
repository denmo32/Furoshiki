package component

import (
	"furoshiki/event"
	"furoshiki/style"
)

// --- LayoutableWidget ---
// LayoutableWidgetは、Widgetインターフェースの基本的な実装を提供する構造体です。
// 他の具体的なウィジェット（Button, Labelなど）は、この構造体を埋め込むことで、
// 位置、サイズ、スタイル、イベント処理などの共通機能を利用します。
type LayoutableWidget struct {
	// --- Position & Size ---
	x, y          int
	width, height int
	minWidth, minHeight int
	// requestedX, requestedY は、AbsoluteLayoutなどの特定のレイアウトに対して
	// ウィジェットが希望する相対位置を保持します。
	requestedX, requestedY int

	// --- Layout & Style ---
	flex             int
	style            style.Style
	relayoutBoundary bool // このウィジェットをレイアウト計算の境界とするかのフラグ

	// --- State ---
	dirty         bool // 再描画が必要かどうかのフラグ
	relayoutDirty bool // 再レイアウトが必要かどうかのフラグ
	isHovered     bool
	isVisible     bool // 可視性フラグ

	// --- Hierarchy & Events ---
	parent        Container // 親コンテナへの参照
	eventHandlers map[event.EventType]event.EventHandler
}

// NewLayoutableWidget は、デフォルト値で LayoutableWidget を初期化します。
func NewLayoutableWidget() *LayoutableWidget {
	return &LayoutableWidget{
		isVisible:     true, // デフォルトで表示状態にする
		eventHandlers: make(map[event.EventType]event.EventHandler),
	}
}