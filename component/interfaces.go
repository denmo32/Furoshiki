package component

import (
    "furoshiki/event"
    "furoshiki/style"

    "github.com/hajimehoshi/ebiten/v2"
)

// --- Widget Interface ---
// Widgetは全てのUI要素の基本的な振る舞いを定義するインターフェースです。
type Widget interface {
    // --- ライフサイクルメソッド ---
    // Updateはウィジェットの状態を更新します。通常、毎フレーム呼び出されます。
    Update()
    // Drawはウィジェットをスクリーンに描画します。
    Draw(screen *ebiten.Image)
    // Cleanupはウィジェットが不要になった際にリソースを解放します。
    Cleanup()

    // --- 位置とサイズ関連メソッド ---
    PositionSetter
    SizeSetter
    MinSizeSetter

    // --- スタイル関連メソッド ---
    StyleGetterSetter

    // --- 状態管理メソッド ---
    DirtyManager
    InteractiveState

    // --- イベント処理メソッド ---
    EventHandler
    HitTester

    // --- レイアウト関連メソッド ---
    LayoutProperties
    HierarchyManager
}

// PositionSetter はウィジェットの位置を設定・取得するためのインターフェースです
type PositionSetter interface {
    SetPosition(x, y int)
    GetPosition() (x, y int)
}

// SizeSetter はウィジェットのサイズを設定・取得するためのインターフェースです
type SizeSetter interface {
    SetSize(width, height int)
    GetSize() (width, height int)
}

// MinSizeSetter はウィジェットの最小サイズを設定・取得するためのインターフェースです
type MinSizeSetter interface {
    SetMinSize(width, height int)
    GetMinSize() (width, height int)
}

// StyleGetterSetter はウィジェットのスタイルを設定・取得するためのインターフェースです
type StyleGetterSetter interface {
    SetStyle(style style.Style)
    GetStyle() style.Style
}

// DirtyManager はウィジェットのダーティ状態を管理するためのインターフェースです
type DirtyManager interface {
    MarkDirty(relayout bool)
    IsDirty() bool
    ClearDirty()
}

// InteractiveState はウィジェットの対話状態を管理するためのインターフェースです
type InteractiveState interface {
    SetHovered(hovered bool)
    IsHovered() bool
    SetPressed(pressed bool)
    IsPressed() bool
    SetVisible(visible bool)
    IsVisible() bool
    SetDisabled(disabled bool)
    IsDisabled() bool
}

// EventHandler はイベント処理のためのインターフェースです
type EventHandler interface {
    AddEventHandler(eventType event.EventType, handler event.EventHandler)
    RemoveEventHandler(eventType event.EventType)
    HandleEvent(e event.Event)
}

// HitTester はヒットテストのためのインターフェースです
type HitTester interface {
    HitTest(x, y int) Widget
}

// LayoutProperties はレイアウトプロパティを管理するためのインターフェースです
type LayoutProperties interface {
    SetFlex(flex int)
    GetFlex() int
    SetRelayoutBoundary(isBoundary bool)
}

// HierarchyManager は階層構造を管理するためのインターフェースです
type HierarchyManager interface {
    SetParent(parent Container)
    GetParent() Container
}

// --- Container Interface ---
// Containerは子Widgetを持つことができるWidgetです。
type Container interface {
    Widget // ContainerはWidgetのすべての振る舞いを継承します
    
    // --- 子要素管理メソッド ---
    AddChild(child Widget)
    RemoveChild(child Widget)
    GetChildren() []Widget
}