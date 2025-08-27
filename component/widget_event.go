package component

import (
	"furoshiki/event"
	"image"
	"log"
	"runtime/debug"
)

// AddEventHandler は、指定されたイベントタイプに対応するイベントハンドラを登録します。
// NOTE: 複数のハンドラを登録できるように、内部実装がスライスベースに変更されました。
// 同じイベントタイプに対して複数回呼び出すと、ハンドラが順に追加されます。
func (w *LayoutableWidget) AddEventHandler(eventType event.EventType, handler event.EventHandler) {
	if w.eventHandlers == nil {
		// NOTE: eventHandlersの型が map[event.EventType]event.EventHandler から
		//       map[event.EventType][]event.EventHandler に変更されました。
		w.eventHandlers = make(map[event.EventType][]event.EventHandler)
	}
	// 指定されたイベントタイプのハンドラスライスに、新しいハンドラを追加します。
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], handler)
}

// RemoveEventHandler は、指定されたイベントタイプのイベントハンドラをすべて削除します。
func (w *LayoutableWidget) RemoveEventHandler(eventType event.EventType) {
	if w.eventHandlers != nil {
		delete(w.eventHandlers, eventType)
	}
}

// HandleEvent は、ディスパッチャから渡されたイベントを処理します。
// このメソッドはイベントバブリングを実装します。まず自身のハンドラを呼び出し、
// イベントがまだ処理されていない（e.Handled == false）場合、親ウィジェットの
// HandleEventメソッドを再帰的に呼び出します。
func (w *LayoutableWidget) HandleEvent(e *event.Event) {
	// NOTE: 複数のハンドラを順に実行するようにロジックが更新されました。
	if handlers, exists := w.eventHandlers[e.Type]; exists {
		// 登録されているすべてのハンドラをループ処理します。
		for _, handler := range handlers {
			// イベントが既に処理済みの場合、後続のハンドラの実行をスキップします。
			if e.Handled {
				break
			}
			// イベントハンドラ内でパニックが発生してもアプリケーションがクラッシュしないように保護します。
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in event handler: %v\n%s", r, debug.Stack())
					}
				}()
				// NOTE: このハンドラ呼び出しを個別に保護することで、特定のハンドラがパニックを起こしても、
				//       同じイベントに登録された他のハンドラの実行が継続されます。
				//       これは、UIの堅牢性を高めるための意図的な設計です。
				handler(e)
			}()
		}
	}

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
	if !w.IsVisible() || w.IsDisabled() {
		return nil
	}

	wx, wy := w.GetPosition()
	wwidth, wheight := w.GetSize()

	rect := image.Rect(wx, wy, wx+wwidth, wy+wheight)
	if rect.Empty() {
		return nil // サイズが0のウィジェットはヒットしません。
	}

	if !(image.Point{X: x, Y: y}.In(rect)) {
		return nil
	}

	// ヒットした場合、LayoutableWidget自身(w)ではなく、それを埋め込んでいる具象ウィジェット(w.self)を返します。
	return w.self
}