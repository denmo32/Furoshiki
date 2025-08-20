package event

import "github.com/hajimehoshi/ebiten/v2"

// --- Event System ---

// EventType はUIイベントの種類（クリック、マウスオーバーなど）を定義します。
type EventType int

const (
	EventClick EventType = iota // マウスクリックイベント
	MouseEnter                 // マウスカーソルが要素に入ったイベント
	MouseLeave                 // マウスカーソルが要素から離れたイベント
	MouseMove                  // マウスカーソルが要素上で移動したイベント
)

// Event は、UIコンポーネント間でやり取りされるイベント情報を保持します。
type Event struct {
	// Type はイベントの種類を示します。
	Type EventType
	// Target はイベントが発生したコンポーネントへの参照です。
	// 型安全性を向上させるため、公開された EventTarget インターフェースを使用します。
	Target EventTarget
	// X, Y はマウスイベントが発生したスクリーン座標です。
	X, Y int
	// Timestamp はイベントが発生した時刻です。
	Timestamp int64
	// MouseButton は押されたマウスのボタンを示します。
	MouseButton ebiten.MouseButton
}

// EventHandler は、特定のイベントタイプに応答するための関数シグネチャです。
type EventHandler func(e Event)