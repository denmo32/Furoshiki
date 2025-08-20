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
	// selfは、このLayoutableWidgetを埋め込んでいる具象ウィジェット(Button, Labelなど)への参照です。
	// これにより、HitTestのようなメソッドが、具象型を正しく返すことができます。
	self Widget
}

// NewLayoutableWidget は、デフォルト値で LayoutableWidget を初期化します。
// 引数には、この構造体を埋め込む具象ウィジェット自身のインスタンス(self)を渡す必要があります。
func NewLayoutableWidget(self Widget) *LayoutableWidget {
	if self == nil {
		// selfがnilの場合、プログラムが予期せぬ動作をする可能性があるため、パニックを発生させます。
		// これは、ウィジェットのコンストラクタが常に正しいインスタンスを渡すことを強制する設計上の決定です。
		panic("NewLayoutableWidget: self cannot be nil")
	}
	return &LayoutableWidget{
		self:          self,
		isVisible:     true, // デフォルトで表示状態にする
		eventHandlers: make(map[event.EventType]event.EventHandler),
	}
}