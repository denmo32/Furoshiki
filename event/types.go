package event

import "github.com/hajimehoshi/ebiten/v2"

// EventType はUIイベントの種類（クリック、マウスオーバーなど）を定義します。
type EventType int

const (
	EventClick EventType = iota
	MouseEnter
	MouseLeave
	MouseMove
	MouseDown
	MouseUp
	MouseScroll
)

// Event は、UIコンポーネント間でやり取りされるイベント情報を保持します。
type Event struct {
	Type             EventType
	Target           EventTarget
	X, Y             int
	Timestamp        int64
	MouseButton      ebiten.MouseButton
	ScrollX, ScrollY float64
	// Handledは、イベントがウィジェットによって処理されたことを示します。
	// これがtrueに設定されると、イベントの親ウィジェットへの伝播（バブリング）が停止します。
	Handled bool
}

// EventHandler は、特定のイベントタイプに応答するための関数シグネチャーです。
type EventHandler func(e *Event)

// EventTargetは、Dispatcherがイベントをディスパッチするためにウィジェットが満たすべき最低限の振る舞いを定義するインターフェースです。
// このインターフェースをeventパッケージ内で定義することで、componentパッケージへの依存をなくし、インポートサイクルを回避します。
type EventTarget interface {
	SetHovered(bool)
	SetPressed(bool)
	HandleEvent(e *Event)
}
