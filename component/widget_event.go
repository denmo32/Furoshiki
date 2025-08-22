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

// HandleEvent は、ディスパッチャから渡されたイベントを処理します。
// 対応するイベントタイプのハンドラが存在すれば、それを実行します。
func (w *LayoutableWidget) HandleEvent(e event.Event) {
    // 状態管理: マウスの押下状態を更新
    switch e.Type {
    case event.MouseDown:
        w.SetPressed(true)
    case event.MouseUp:
        w.SetPressed(false)
    }

    if handler, exists := w.eventHandlers[e.Type]; exists {
        // イベントハンドラの実行中にパニックが発生してもアプリケーション全体がクラッシュしないようにリカバリします。
        defer func() {
            if r := recover(); r != nil {
                // パニック発生時に、より詳細なデバッグ情報（スタックトレース）をログに出力します。
                log.Printf("Recovered from panic in event handler: %v\n%s", r, debug.Stack())
            }
        }()
        handler(e)
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