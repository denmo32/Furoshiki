package component

import (
    "furoshiki/event"
    "furoshiki/style"
)

// LayoutableWidgetは、Widgetインターフェースの基本的な実装を提供する構造体です。
// 他の具体的なウィジェット（Button, Labelなど）は、この構造体を埋め込むことで、
// 位置、サイズ、スタイル、イベント処理などの共通機能を利用します。
type LayoutableWidget struct {
    // --- Position & Size ---
    position      position
    size          size
    minSize       size
    requestedPos  position

    // --- Layout & Style ---
    layout        layoutProperties
    style         style.Style

    // --- State ---
    state         widgetState

    // --- Hierarchy & Events ---
    hierarchy     hierarchy
    eventHandlers map[event.EventType]event.EventHandler
    self          Widget
}

// position はウィジェットの位置情報を保持します
type position struct {
    x, y int
}

// size はウィジェットのサイズ情報を保持します
type size struct {
    width, height int
}

// layoutProperties はレイアウト関連のプロパティを保持します
type layoutProperties struct {
    flex             int
    relayoutBoundary bool
}

// widgetState はウィジェットの状態を保持します
type widgetState struct {
    dirty         bool
    relayoutDirty bool
    isHovered     bool
    isPressed     bool
    isVisible     bool
    isDisabled    bool
}

// hierarchy はウィジェットの階層構造情報を保持します
type hierarchy struct {
    parent Container
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
        state:         widgetState{isVisible: true}, // デフォルトで表示状態にする
        eventHandlers: make(map[event.EventType]event.EventHandler),
    }
}