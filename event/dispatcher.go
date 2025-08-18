package event

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// eventTarget は、Dispatcherがイベントをディスパッチするためにウィジェットが満たすべき最低限の振る舞いを定義するインターフェースです。
// このインターフェースをeventパッケージ内で定義することで、componentパッケージへの依存をなくし、インポートサイクルを回避します。
// component.Widgetは（構造的に）このインターフェースを満たします。
type eventTarget interface {
	SetHovered(bool)
	HandleEvent(Event)
}

// --- Global Event Dispatcher ---

// Dispatcherは、UIイベントを一元管理し、適切なコンポーネントにディスパッチします。
// シングルトンパターンで実装され、UIツリー全体で唯一のインスタンスを共有します。
type Dispatcher struct {
	hoveredComponent eventTarget // 型を具体的なウィジェットからインターフェースに変更
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
// target引数には、ヒットテストの結果（通常はcomponent.Widget）がinterface{}型として渡されます。
func (d *Dispatcher) Dispatch(target interface{}, cx, cy int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 渡されたtargetをeventTargetインターフェースに変換できるか試みる。
	// これにより、イベント処理に必要なメソッドを持たないオブジェクトを安全に無視できる。
	var currentTarget eventTarget
	if target != nil {
		if ct, ok := target.(eventTarget); ok {
			currentTarget = ct
		}
	}

	// ホバー状態の更新
	if currentTarget != d.hoveredComponent {
		// 以前ホバーされていたコンポーネントからマウスが離れた
		if d.hoveredComponent != nil {
			d.hoveredComponent.SetHovered(false)
			d.hoveredComponent.HandleEvent(Event{
				Type:   MouseLeave,
				Target: d.hoveredComponent,
			})
		}
		// 新しいコンポーネントにマウスが入った
		if currentTarget != nil {
			currentTarget.SetHovered(true)
			currentTarget.HandleEvent(Event{
				Type:   MouseEnter,
				Target: currentTarget,
			})
		}
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
// これは、UIツリーが完全に再構築される場合などに役立ちます。
func (d *Dispatcher) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.hoveredComponent = nil
}