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
	MouseDown                  // マウスボタンが要素上で押されたイベント
	MouseUp                    // マウスボタンが要素上で解放されたイベント
	// [新規追加] MouseScrollはマウスホイールのスクロールイベントです。
	MouseScroll
)

// Event は、UIコンポーネント間でやり取りされるイベント情報を保持します。
type Event struct {
	// Type はイベントの種類を示します。
	Type EventType
	// Target はイベントが発生したコンポーネントへの参照です。
	Target EventTarget
	// X, Y はマウスイベントが発生したスクリーン座標です。
	X, Y int
	// Timestamp はイベントが発生した時刻です。
	Timestamp int64
	// MouseButton は押されたマウスのボタンを示します。
	MouseButton ebiten.MouseButton
	// [新規追加] ScrollX, ScrollY はマウスホイールのスクロール量です。
	// 正の値は通常、右または下へのスクロールを示します。
	ScrollX, ScrollY float64
	// [追加] Handledは、イベントがウィジェットによって処理されたことを示します。
	// これがtrueに設定されると、イベントの親ウィジェットへの伝播（バブリング）が停止します。
	Handled bool
}

// [修正] EventHandler は、特定のイベントタイプに応答するための関数シグネチャーです。
// イベントの伝播を停止できるよう、イベントをポインタで受け取ります。
type EventHandler func(e *Event)

// [修正] EventTargetは、Dispatcherがイベントをディスパッチするためにウィジェットが満たすべき最低限の振る舞いを定義するインターフェースです。
// このインターフェースをeventパッケージ内で定義することで、componentパッケージへの依存をなくし、インポートサイクルを回避します。
// component.Widgetは（構造的に）このインターフェースを満たします。
type EventTarget interface {
	SetHovered(bool)
	SetPressed(bool) // ウィジェットの押下状態を設定するメソッドを追加
	// [修正] イベントが処理されたことをディスパッチャに伝えるため、イベントをポインタで受け取ります。
	HandleEvent(e *Event)
}