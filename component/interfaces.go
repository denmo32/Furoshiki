package component

import (
	"furoshiki/event"
	"furoshiki/style"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Widget Interface ---
// Widgetは全てのUI要素の基本的な振る舞いを定義するインターフェースです。
// これには位置、サイズ、スタイル、イベント処理などが含まれます。
type Widget interface {
	Update()
	Draw(screen *ebiten.Image)
	SetPosition(x, y int)
	GetPosition() (x, y int)
	SetSize(width, height int)
	GetSize() (width, height int)
	SetMinSize(width, height int)
	GetMinSize() (width, height int)
	SetStyle(style style.Style)
	// GetStyle はウィジェットの現在のスタイルを返します。
	// [修正] 戻り値をポインタ型(*style.Style)から値型(style.Style)に変更します。
	// これにより、GetStyle()経由でスタイル構造体自体を変更されることを防ぎ、SetStyle()の利用を強制します。
	// 注意: このメソッドはスタイルの浅いコピーを返します。返されたStyle構造体のポインタフィールド
	// (例: Background, Padding) が指す先の値を変更すると、元のウィジェットのスタイルに影響が及ぶ
	// 可能性があります。スタイル全体を安全に変更するにはSetStyle()を使用してください。
	GetStyle() style.Style
	MarkDirty(relayout bool)
	IsDirty() bool
	ClearDirty()
	AddEventHandler(eventType event.EventType, handler event.EventHandler)
	RemoveEventHandler(eventType event.EventType)
	HandleEvent(event event.Event)
	SetFlex(flex int)
	GetFlex() int
	SetParent(parent Container) // 親はContainer型である必要があります
	GetParent() Container
	HitTest(x, y int) Widget
	SetHovered(hovered bool)
	IsHovered() bool
	SetVisible(visible bool)
	IsVisible() bool
	SetRelayoutBoundary(isBoundary bool) // レイアウト境界フラグを設定
	Cleanup()                            // コンポーネントのクリーンアップ処理

	// Widgetはイベントディスパッチャが要求するevent.EventTargetインターフェースを
	// 構造的に満たす必要があります。
	// SetHovered(bool)
	// HandleEvent(event.Event)
}

// --- Container Interface ---
// Containerは子Widgetを持つことができるWidgetです。
// UIの階層構造を構築するために使用されます。
type Container interface {
	Widget // ContainerはWidgetのすべての振る舞いを継承します
	AddChild(child Widget)
	RemoveChild(child Widget)
	GetChildren() []Widget
}