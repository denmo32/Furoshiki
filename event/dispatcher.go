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
	// 現在マウスボタンが押下されているコンポーネントを追跡します。
	pressedComponent EventTarget
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

	// 1. ホバー状態の更新 (MouseEnter, MouseLeave)
	if currentTarget != d.hoveredComponent {
		if d.hoveredComponent != nil {
			d.hoveredComponent.SetHovered(false)
			d.hoveredComponent.HandleEvent(Event{Type: MouseLeave, Target: d.hoveredComponent, X: cx, Y: cy})
		}
		if currentTarget != nil {
			currentTarget.SetHovered(true)
			currentTarget.HandleEvent(Event{Type: MouseEnter, Target: currentTarget, X: cx, Y: cy})
		}
		d.hoveredComponent = currentTarget
	}

	// 2. マウス移動イベント (MouseMove)
	if d.hoveredComponent != nil {
		d.hoveredComponent.HandleEvent(Event{Type: MouseMove, Target: d.hoveredComponent, X: cx, Y: cy})
	}

	// 3. マウスボタン押下イベント (MouseDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if d.hoveredComponent != nil {
			d.pressedComponent = d.hoveredComponent // 押されたコンポーネントを記憶
			d.pressedComponent.HandleEvent(Event{
				Type:        MouseDown,
				Target:      d.pressedComponent,
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft,
			})
		}
	}

	// 4. マウスボタン解放イベント (MouseUp and Click)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// MouseUpイベントは、ボタンが離された時点のホバー要素に発行
		if d.hoveredComponent != nil {
			d.hoveredComponent.HandleEvent(Event{
				Type:        MouseUp,
				Target:      d.hoveredComponent,
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft,
			})
		}

		// Clickイベントは、押したコンポーネントと離したコンポーネントが同じ場合に発行
		if d.pressedComponent != nil && d.pressedComponent == d.hoveredComponent {
			d.pressedComponent.HandleEvent(Event{
				Type:        EventClick,
				Target:      d.pressedComponent,
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft,
			})
		}
		// 押下状態をリセット
		d.pressedComponent = nil
	}
}

// Reset は、ディスパッチャの内部状態をリセットします。
// これは、UIツリーが完全に再構築される場合や、モーダルウィンドウを閉じた後などに役立ちます。
func (d *Dispatcher) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.hoveredComponent = nil
	d.pressedComponent = nil
}