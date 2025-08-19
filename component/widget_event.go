package component

import (
	"furoshiki/event"
	"image"
	"log"           // [追加] ログ出力のために追加
	"runtime/debug" // [追加] スタックトレース取得のために追加
)

func (w *LayoutableWidget) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if w.eventHandlers == nil {
		w.eventHandlers = make(map[event.EventType]event.EventHandler)
	}
	w.eventHandlers[eventType] = handler
}

func (w *LayoutableWidget) RemoveEventHandler(eventType event.EventType) {
	if w.eventHandlers != nil {
		delete(w.eventHandlers, eventType)
	}
}

func (w *LayoutableWidget) HandleEvent(event event.Event) {
	if handler, exists := w.eventHandlers[event.Type]; exists {
		// イベントハンドラの実行中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
		defer func() {
			if r := recover(); r != nil {
				// [改善] パニック発生時に、より詳細なデバッグ情報（スタックトレース）をログに出力します。
				log.Printf("Recovered from panic in event handler: %v\n%s", r, debug.Stack())
			}
		}()
		handler(event)
	}
}

// HitTest は、指定された座標がウィジェットの領域内にあるかを判定します。
// [改善] 戻り値として、初期化時に設定された具象ウィジェットへの参照(w.self)を返します。
// これにより、ButtonやLabelなどの具象ウィジェット側でこのメソッドをオーバーライドする必要がなくなります。
func (w *LayoutableWidget) HitTest(x, y int) Widget {
	if !w.isVisible {
		return nil
	}
	// 境界チェックをより明確に実装
	rect := image.Rect(w.x, w.y, w.x+w.width, w.y+w.height)
	if rect.Empty() {
		return nil
	}
	if !(image.Point{X: x, Y: y}.In(rect)) {
		return nil
	}
	// ヒットした場合、LayoutableWidget自身(w)ではなく、それを埋め込んでいる具象ウィジェット(w.self)を返します。
	return w.self
}