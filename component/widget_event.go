package component

import (
	"fmt"
	"furoshiki/event"
	"image"
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
				fmt.Printf("Recovered from panic in event handler: %v\n", r)
			}
		}()
		handler(event)
	}
}

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
	return w
}