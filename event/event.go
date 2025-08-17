package event

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	Type        EventType
	Target      interface{} // イベントが発生したコンポーネント
	X, Y        int         // マウスイベントの場合の座標
	Timestamp   int64
	MouseButton ebiten.MouseButton
}

// EventHandler は、特定のイベントタイプに応答するための関数シグネチャです。
type EventHandler func(e Event)

// --- Global Event Dispatcher ---

// Dispatcherは、UIイベントを一元管理し、適切なコンポーネントにディスパッチします。
// シングルトンパターンで実装され、UIツリー全体で唯一のインスタンスを共有します。
type Dispatcher struct {
	hoveredComponent interface{}
	mutex            sync.Mutex
}

var (
	instance *Dispatcher
	once     sync.Once
)

// GetDispatcher は、Dispatcherのシングルトンインスタンスを返します。
func GetDispatcher() *Dispatcher {
	once.Do(func() {
		instance = &Dispatcher{}
	})
	return instance
}

// Dispatch は、マウスイベントを処理し、適切なイベントをコンポーネントに発行します。
// このメソッドは、アプリケーションのメインUpdateループから毎フレーム呼び出されることを想定しています。
func (d *Dispatcher) Dispatch(target interface{}, cx, cy int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// ホバー状態の更新
	if target != d.hoveredComponent {
		// 以前ホバーされていたコンポーネントからマウスが離れた
		if d.hoveredComponent != nil {
			if comp, ok := d.hoveredComponent.(interface {
				SetHovered(bool)
				HandleEvent(Event)
			}); ok {
				comp.SetHovered(false)
				comp.HandleEvent(Event{
					Type:   MouseLeave,
					Target: d.hoveredComponent,
				})
			}
		}
		// 新しいコンポーネントにマウスが入った
		if target != nil {
			if comp, ok := target.(interface {
				SetHovered(bool)
				HandleEvent(Event)
			}); ok {
				comp.SetHovered(true)
				comp.HandleEvent(Event{
					Type:   MouseEnter,
					Target: target,
				})
			}
		}
		d.hoveredComponent = target
	}

	// マウス移動イベント
	if d.hoveredComponent != nil {
		if comp, ok := d.hoveredComponent.(interface{ HandleEvent(Event) }); ok {
			comp.HandleEvent(Event{
				Type:   MouseMove,
				Target: d.hoveredComponent,
				X:      cx,
				Y:      cy,
			})
		}
	}

	// マウスクリックイベント
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if d.hoveredComponent != nil {
			if comp, ok := d.hoveredComponent.(interface{ HandleEvent(Event) }); ok {
				comp.HandleEvent(Event{
					Type:        EventClick,
					Target:      d.hoveredComponent,
					X:           cx,
					Y:           cy,
					Timestamp:   time.Now().UnixNano(),
					MouseButton: ebiten.MouseButtonLeft,
				})
			}
		}
	}
}
