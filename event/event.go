package event

import "github.com/hajimehoshi/ebiten/v2"

// EventTypeはイベントの種類を定義します。
type EventType int

const (
	EventClick EventType = iota
	MouseEnter           // マウスカーソルがコンポーネントの領域に入った
	MouseLeave           // マウスカーソルがコンポーネントの領域から出た
	MouseMove            // マウスカーソルがコンポーネントの領域内で移動した
	EventDrag
	EventFocus
	EventBlur
)

// EventHandlerはイベントを処理する関数です。
type EventHandler func(event Event)

// EventはUIイベントを表します。
type Event struct {
	Type EventType
	// eventパッケージがcomponentパッケージに依存しないようにするため、
	// ターゲットは汎用的なinterface{}を使用します。
	Target      interface{}
	Timestamp   int64
	X, Y        int                 // イベントが発生したカーソルの座標
	MouseButton ebiten.MouseButton  // クリックイベントで押されたマウスボタンを追加
	// Handledフラグは、イベントが処理済みであり、
	// 親コンポーネントへの伝播（バブリング）を停止すべきかを示します。
	Handled bool
	Data    interface{}
}