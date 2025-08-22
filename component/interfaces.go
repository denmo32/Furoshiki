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
	// Updateはウィジェットの状態を更新します。通常、毎フレーム呼び出されます。
	Update()
	// Drawはウィジェットをスクリーンに描画します。
	Draw(screen *ebiten.Image)
	// SetPositionはウィジェットの絶対スクリーン座標を設定します。
	SetPosition(x, y int)
	// GetPositionはウィジェットの絶対スクリーン座標を取得します。
	GetPosition() (x, y int)
	// SetSizeはウィジェットのサイズを設定します。
	SetSize(width, height int)
	// GetSizeはウィジェットのサイズを取得します。
	GetSize() (width, height int)
	// SetMinSizeはウィジェットの最小サイズを明示的に設定します。レイアウト計算時に考慮されます。
	SetMinSize(width, height int)
	// GetMinSizeはウィジェットが取るべき最小サイズを取得します。
	// この値は、ウィジェットのコンテンツ（テキストなど）と、SetMinSizeでユーザーが
	// 明示的に設定した値の両方を考慮して決定されます。
	GetMinSize() (width, height int)
	// SetStyleはウィジェットのスタイルを設定します。
	SetStyle(style style.Style)
	// GetStyle はウィジェットの現在のスタイルの安全なコピーを返します。
	// このメソッドが返すStyleオブジェクトはディープコピーされているため、
	// 変更しても元のウィジェットに影響はありません。
	// スタイルを変更したい場合は、常に `SetStyle()` を使用してください。
	GetStyle() style.Style
	// MarkDirtyはウィジェットの状態が変更されたことをマークし、再描画や再レイアウトを要求します。
	// relayoutがtrueの場合、ウィジェットのサイズや位置に関する変更があったことを示し、親コンテナにレイアウトの再計算を要求します。
	MarkDirty(relayout bool)
	// IsDirtyはウィジェットが再描画または再レイアウトを必要とするかどうかを返します。
	IsDirty() bool
	// ClearDirtyはダーティフラグをクリアします。通常、コンテナがレイアウトを完了した後に呼び出されます。
	ClearDirty()
	// AddEventHandlerは指定されたイベントタイプのハンドラを登録します。
	AddEventHandler(eventType event.EventType, handler event.EventHandler)
	// RemoveEventHandlerは指定されたイベントタイプのハンドラを削除します。
	RemoveEventHandler(eventType event.EventType)
	// HandleEventはイベントを処理します。通常、イベントディスパッチャによって呼び出されます。
	HandleEvent(e event.Event)
	// SetFlexはFlexLayoutにおけるウィジェットの伸縮係数を設定します。
	SetFlex(flex int)
	// GetFlexはウィジェットの伸縮係数を取得します。
	GetFlex() int
	// SetParentはウィジェットの親コンテナを設定します。
	SetParent(parent Container)
	// GetParentはウィジェットの親コンテナを取得します。
	GetParent() Container
	// HitTestは指定された座標がウィジェットの領域内にあるかを判定し、ヒットした場合はウィジェット自身を返します。
	HitTest(x, y int) Widget
	// SetHoveredはマウスカーソルがウィジェット上にあるかどうかの状態を設定します。
	SetHovered(hovered bool)
	// IsHoveredはマウスカーソルがウィジェット上にあるかどうかの状態を返します。
	IsHovered() bool
	// SetVisibleはウィジェットの可視性を設定します。非表示のウィジェットは更新、描画、レイアウト計算の対象外となります。
	SetVisible(visible bool)
	// IsVisibleはウィジェットが可視状態であるかを返します。
	IsVisible() bool
	// SetDisabledはウィジェットの有効・無効状態を設定します。無効なウィジェットはユーザー入力を受け付けません。
	SetDisabled(disabled bool)
	// IsDisabledはウィジェットが無効状態であるかを返します。
	IsDisabled() bool
	// SetRelayoutBoundaryは、このウィジェットをレイアウト計算の境界とするか設定します。
	// trueに設定すると、このウィジェット内部の変更が親コンテナの再レイアウトを引き起こさなくなり、パフォーマンスが向上します。
	SetRelayoutBoundary(isBoundary bool)
	// Cleanupはウィジェットが不要になった際に、イベントハンドラや親への参照などのリソースを解放します。
	Cleanup()

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
	// AddChildはコンテナに子ウィジェットを追加します。
	AddChild(child Widget)
	// RemoveChildはコンテナから子ウィジェットを削除します。
	RemoveChild(child Widget)
	// GetChildrenはコンテナが保持するすべての子ウィジェットのスライスを返します。
	GetChildren() []Widget
}