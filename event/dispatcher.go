package event

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// EventTargetは、Dispatcherがイベントをディスパッチするためにウィジェットが満たすべき最低限の振る舞いを定義するインターフェースです。
// このインターフェースをeventパッケージ内で定義することで、componentパッケージへの依存をなくし、インポートサイクルを回避します。
// component.Widgetは（構造的に）このインターフェースを満たします。
type EventTarget interface {
	SetHovered(bool)
	HandleEvent(Event)
}

// --- Global Event Dispatcher ---

// Dispatcherは、UIイベントを一元管理し、適切なコンポーネントにディスパッチします。
// シングルトンパターンで実装され、UIツリー全体で唯一のインスタンスを共有します。
type Dispatcher struct {
	// 現在マウスカーソルがホバーしているコンポーネントを追跡します。
	hoveredComponent EventTarget
	mutex            sync.Mutex
}

var (
	instance *Dispatcher
	once     sync.Once
)

// GetDispatcher は、Dispatcherのシングルトンインスタンスをスレッドセーフに返します。
func GetDispatcher() *Dispatcher {
	once.Do(func() {
		instance = &Dispatcher{}
	})
	return instance
}

// Dispatch は、マウスイベントを処理し、適切なイベントをコンポーネントに発行します。
// このメソッドは、アプリケーションのメインUpdateループから毎フレーム呼び出されることを想定しています。
// target引数には、HitTestによって特定された、現在マウスカーソル下にあるウィジェットを渡します。
func (d *Dispatcher) Dispatch(target EventTarget, cx, cy int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	currentTarget := target

	// ホバー状態の更新
	if currentTarget != d.hoveredComponent {
		// 1. 以前ホバーされていたコンポーネントからマウスが離れた場合 (MouseLeave)
		if d.hoveredComponent != nil {
			d.hoveredComponent.SetHovered(false)
			d.hoveredComponent.HandleEvent(Event{
				Type:   MouseLeave,
				Target: d.hoveredComponent,
			})
		}
		// 2. 新しいコンポーネントにマウスが入った場合 (MouseEnter)
		if currentTarget != nil {
			currentTarget.SetHovered(true)
			currentTarget.HandleEvent(Event{
				Type:   MouseEnter,
				Target: currentTarget,
			})
		}
		// 3. ホバー中のコンポーネントを更新
		d.hoveredComponent = currentTarget
	}

	// マウス移動イベント (ホバー中のコンポーネントに対してのみ発行)
	if d.hoveredComponent != nil {
		d.hoveredComponent.HandleEvent(Event{
			Type:   MouseMove,
			Target: d.hoveredComponent,
			X:      cx,
			Y:      cy,
		})
	}

	// マウスクリックイベント (ホバー中のコンポーネントに対してのみ発行)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if d.hoveredComponent != nil {
			d.hoveredComponent.HandleEvent(Event{
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

// Reset は、ディスパッチャの内部状態をリセットします。
// これは、UIツリーが完全に再構築される場合や、モーダルウィンドウを閉じた後などに役立ちます。
func (d *Dispatcher) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.hoveredComponent = nil
}