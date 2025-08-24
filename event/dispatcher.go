package event

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
// [修正] イベントの Handled フラグを共有するため、イベントはポインタで渡されます。
func (d *Dispatcher) Dispatch(target EventTarget, cx, cy int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	currentTarget := target

	// 1. ホバー状態の更新 (MouseEnter, MouseLeave)
	if currentTarget != d.hoveredComponent {
		if d.hoveredComponent != nil {
			d.hoveredComponent.SetHovered(false)
			d.hoveredComponent.HandleEvent(&Event{Type: MouseLeave, Target: d.hoveredComponent, X: cx, Y: cy})
		}
		if currentTarget != nil {
			currentTarget.SetHovered(true)
			currentTarget.HandleEvent(&Event{Type: MouseEnter, Target: currentTarget, X: cx, Y: cy})
		}
		d.hoveredComponent = currentTarget
	}

	// 2. マウス移動イベント (MouseMove)
	if d.hoveredComponent != nil {
		d.hoveredComponent.HandleEvent(&Event{Type: MouseMove, Target: d.hoveredComponent, X: cx, Y: cy})
	}

	// 3. マウスボタン押下イベント (MouseDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if d.hoveredComponent != nil {
			d.pressedComponent = d.hoveredComponent // 押されたコンポーネントを記憶

			// --- 改善点 ---
			// ウィジェットの押下状態(Pressed)の管理をDispatcherが一元的に行うように変更。
			// これにより、状態遷移のロジックがDispatcherに集約され、
			// 各ウィジェットのHandleEvent実装がシンプルになります。
			d.pressedComponent.SetPressed(true)

			d.pressedComponent.HandleEvent(&Event{
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
		if d.pressedComponent != nil {
			// --- 改善点 ---
			// マウスボタンが解放されたタイミングで、押されていたウィジェットの
			// 押下状態(Pressed)を解除します。これにより状態管理の責務がDispatcherに集約されます。
			d.pressedComponent.SetPressed(false)

			// MouseUpイベントは、現在ホバーしているコンポーネントではなく、
			// 最初に「押された」コンポーネント(pressedComponent)に送る必要があります。
			// これにより、ボタンの外でマウスを離した場合でも、押されたボタンが確実に
			// MouseUpイベントを受け取れるようになります。
			mouseUpEvent := &Event{
				Type:        MouseUp,
				Target:      d.pressedComponent, // イベントのターゲットは押されていたコンポーネント
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft,
			}
			d.pressedComponent.HandleEvent(mouseUpEvent)

			// 次に、Clickイベントを発行するかどうかを決定します。
			// クリックが成立するのは、マウスを押したコンポーネントと離したコンポーネントが同じ場合のみです。
			if d.pressedComponent == d.hoveredComponent {
				d.pressedComponent.HandleEvent(&Event{
					Type:        EventClick,
					Target:      d.pressedComponent,
					X:           cx,
					Y:           cy,
					Timestamp:   time.Now().UnixNano(),
					MouseButton: ebiten.MouseButtonLeft,
				})
			}
		}

		// イベント処理が完了したら、押下状態をリセットします。
		d.pressedComponent = nil
	}

	// [新規追加] 5. マウスホイールイベント (MouseScroll)
	// Ebitenからホイールの移動量を取得します。
	wheelX, wheelY := ebiten.Wheel()
	// 移動量があり、かつカーソルが何らかのウィジェット上にある場合にイベントを発行します。
	if (wheelX != 0 || wheelY != 0) && d.hoveredComponent != nil {
		d.hoveredComponent.HandleEvent(&Event{
			Type:    MouseScroll,
			Target:  d.hoveredComponent,
			X:       cx,
			Y:       cy,
			ScrollX: wheelX,
			ScrollY: wheelY,
		})
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
