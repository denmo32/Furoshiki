package component

import (
	"furoshiki/event"
	"image"
	"log"
	"runtime/debug"
)

// AddEventHandler は、指定されたイベントタイプに対応するイベントハンドラを登録します。
// 同じイベントタイプにハンドラが既に存在する場合、上書きされます。
func (w *LayoutableWidget) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if w.eventHandlers == nil {
		w.eventHandlers = make(map[event.EventType]event.EventHandler)
	}
	w.eventHandlers[eventType] = handler
}

// RemoveEventHandler は、指定されたイベントタイプのイベントハンドラを削除します。
func (w *LayoutableWidget) RemoveEventHandler(eventType event.EventType) {
	if w.eventHandlers != nil {
		delete(w.eventHandlers, eventType)
	}
}

// [修正] HandleEvent は、ディスパッチャから渡されたイベントを処理します。
// このメソッドはイベントバブリングを実装します。まず自身のハンドラを呼び出し、
// イベントがまだ処理されていない（e.Handled == false）場合、親ウィジェットの
// HandleEventメソッドを再帰的に呼び出します。
func (w *LayoutableWidget) HandleEvent(e *event.Event) {
	// 状態管理ロジックはDispatcherに集約されたため削除。
	// このメソッドは、登録されたカスタムイベントハンドラの実行に専念します。
	if handler, exists := w.eventHandlers[e.Type]; exists {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in event handler: %v\n%s", r, debug.Stack())
			}
		}()
		handler(e)
	}

	// [追加] イベントバブリングのロジック
	// イベントがこのウィジェットで処理されておらず（Handledがfalse）、
	// かつ親ウィジェットが存在する場合に、イベントを親に伝播させます。
	if e != nil && !e.Handled && w.hierarchy.parent != nil {
		w.hierarchy.parent.HandleEvent(e)
	}
}

// HitTest は、指定された座標がウィジェットの領域内にあるかを判定します。
// 戻り値として、初期化時に設定された具象ウィジェットへの参照(w.self)を返します。
// これにより、ButtonやLabelなどの具象ウィジェット側でこのメソッドをオーバーライドする必要がなくなります。
func (w *LayoutableWidget) HitTest(x, y int) Widget {
	// isVisible/isDisabledフィールドではなくIsVisible()/IsDisabled()メソッドを使用するように修正
	if !w.IsVisible() || w.IsDisabled() {
		return nil
	}

	// x, y, width, heightフィールドへの直接アクセスではなく、GetPosition()/GetSize()メソッドを使用するように修正
	wx, wy := w.GetPosition()
	wwidth, wheight := w.GetSize()

	// ウィジェットの矩形領域を定義します。
	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if rect.Empty() {
		return nil // サイズが0のウィジェットはヒットしません。
	}

	// 指定された座標が矩形内にあるかチェックします。
	if !(image.Point{X: x, Y: y}.In(rect)) {
		return nil
	}

	// ヒットした場合、LayoutableWidget自身(w)ではなく、それを埋め込んでいる具象ウィジェット(w.self)を返します。
	return w.self
}