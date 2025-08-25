package event

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Dispatcherは、UIイベントを一元管理し、適切なコンポーネントにディスパッチします。
type Dispatcher struct {
	hoveredComponent EventTarget
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
func (d *Dispatcher) Dispatch(target EventTarget, cx, cy int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 1. ホバー状態の更新 (MouseEnter, MouseLeave)
	if target != d.hoveredComponent {
		if d.hoveredComponent != nil {
			d.hoveredComponent.SetHovered(false)
			d.hoveredComponent.HandleEvent(&Event{Type: MouseLeave, Target: d.hoveredComponent, X: cx, Y: cy})
		}
		if target != nil {
			target.SetHovered(true)
			target.HandleEvent(&Event{Type: MouseEnter, Target: target, X: cx, Y: cy})
		}
		d.hoveredComponent = target
	}

	// 2. マウス移動イベント (MouseMove)
	if d.hoveredComponent != nil {
		d.hoveredComponent.HandleEvent(&Event{Type: MouseMove, Target: d.hoveredComponent, X: cx, Y: cy})
	}

	// 3. マウスボタン押下イベント (MouseDown)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if d.hoveredComponent != nil {
			d.pressedComponent = d.hoveredComponent
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
			d.pressedComponent.SetPressed(false)

			// MouseUpイベントは、最初に「押された」コンポーネントに送ります。
			d.pressedComponent.HandleEvent(&Event{
				Type:        MouseUp,
				Target:      d.pressedComponent,
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft,
			})

			// クリックが成立するのは、押したコンポーネントと離したコンポーネントが同じ場合のみです。
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
		d.pressedComponent = nil
	}

	// 5. マウスホイールイベント (MouseScroll)
	wheelX, wheelY := ebiten.Wheel()
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
func (d *Dispatcher) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.hoveredComponent = nil
	d.pressedComponent = nil
}
