package event

import "github.com/hajimehoshi/ebiten/v2"

// --- Event System ---

// EventType はUIイベントの種類（クリック、マウスオーバーなど）を定義します。
type EventType int

const (
	EventClick EventType = iota
	MouseEnter
	MouseLeave
	MouseMove
)

// Event は、UIコンポーネント間でやり取りされるイベント情報を保持します。
type Event struct {
	Type EventType
	// [改善] Targetの型を interface{} から eventTarget に変更し、型安全性を向上させます。
	Target      eventTarget // イベントが発生したコンポーネント
	X, Y        int         // マウスイベントの場合の座標
	Timestamp   int64
	MouseButton ebiten.MouseButton
}

// EventHandler は、特定のイベントタイプに応答するための関数シグネチャです。
type EventHandler func(e Event)