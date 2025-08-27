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

// 【提案2】イベントの伝播を制御するための型を定義します。
// これにより、ハンドラの戻り値でイベントをバブリングさせるか停止させるかを明示的に示せます。
type Propagation int

const (
	// Propagate はイベントの伝播を継続することを示します。
	Propagate Propagation = iota
	// StopPropagation はイベントの伝播を停止することを示します。
	StopPropagation
)

// Event は、UIコンポーネント間でやり取りされるイベント情報を保持します。
type Event struct {
	Type        EventType
	Target      EventTarget
	X, Y        int
	Timestamp   int64
	MouseButton ebiten.MouseButton
	ScrollX, ScrollY float64
	// Handledは、イベントがウィジェットによって処理されたことを示します。
	// これがtrueに設定されると、イベントの親ウィジェットへの伝播（バブリング）が停止します。
	// 【提案2】このフィールドは主に内部で使われ、ハンドラの戻り値によって制御されるようになります。
	Handled bool
}

// EventHandler は、特定のイベントタイプに応答するための関数シグネチャーです。
// 【提案2】戻り値としてPropagation型を返すようにシグネチャが変更されました。
// これにより、ハンドラは副作用なしにイベントの伝播を制御できます。
type EventHandler func(e *Event) Propagation

// EventTargetは、Dispatcherがイベントをディスパッチするためにウィジェットが満たすべき最低限の振る舞いを定義するインターフェースです。
// このインターフェースをeventパッケージ内で定義することで、componentパッケージへの依存をなくし、インポートサイクルを回避します。
// 【提案1】component.Widgetがスリム化されたため、インタラクティブなウィジェットは
// このEventTargetインターフェースを明示的に実装する(またはLayoutableWidgetを埋め込む)必要があります。
type EventTarget interface {
	SetHovered(bool)
	SetPressed(bool)
	HandleEvent(e *Event)
}