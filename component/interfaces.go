package component

import (
	"furoshiki/event"
	"furoshiki/style"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Widget Interface ---
// Widgetは全てのUI要素の基本的な振る舞いを定義するインターフェースです。
// 多くのメソッドを含んでいますが、責務ごとに小さなインターフェースに分割されています。
type Widget interface {
	// --- ライフサイクル ---
	Update()
	Draw(screen *ebiten.Image)
	Cleanup()

	// --- 階層構造 ---
	HierarchyManager

	// --- 位置とサイズ ---
	PositionSetter
	SizeSetter
	MinSizeSetter

	// --- スタイル ---
	StyleGetterSetter

	// --- レイアウト ---
	LayoutProperties
	// SetLayoutData は、このウィジェットにレイアウト固有のデータを設定します。
	// 親コンテナのレイアウトシステム（例: AdvancedGridLayout）がこれを使用して、
	// ウィジェットごとの配置情報（行、列、スパンなど）を管理します。
	SetLayoutData(data any)
	// GetLayoutData は、このウィジェットに設定されたレイアウト固有のデータを返します。
	GetLayoutData() any

	// --- 状態管理 ---
	DirtyManager
	InteractiveState

	// --- イベント処理 ---
	EventHandler
	HitTester
}

// HeightForWider は、ウィジェットが特定の幅を与えられた場合に
// 必要となる高さを計算できることを示すインターフェースです。
// テキストの折り返しなど、コンテンツの高さが幅に依存するウィジェットによって実装されます。
type HeightForWider interface {
	GetHeightForWidth(width int) int
}

// ScrollBarWidget は、ScrollBarが実装すべきメソッドを定義します。
// これにより、他のパッケージが具体的なScrollBar型に依存することなく、
// このインターフェースを通じてScrollBarを操作できます。
type ScrollBarWidget interface {
	Widget // Widgetの基本機能を継承
	SetRatios(contentRatio, scrollRatio float64)
}

// HierarchyManager は階層構造を管理するためのインターフェースです
type HierarchyManager interface {
	SetParent(parent Container)
	GetParent() Container
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
	// NOTE: パフォーマンスが重要な読み取り専用の場面のために、
	// スタイルのディープコピーを生成しないメソッドを追加しました。
	// 返されたスタイルは変更してはいけません。
	ReadOnlyStyle() style.Style
}

// LayoutProperties はレイアウトプロパティを管理するためのインターフェースです
type LayoutProperties interface {
	SetFlex(flex int)
	GetFlex() int
	SetLayoutBoundary(isBoundary bool)
}

// DirtyManager はウィジェットのダーティ状態を管理するためのインターフェースです
type DirtyManager interface {
	MarkDirty(relayout bool)
	IsDirty() bool
	NeedsRelayout() bool // ウィジェットがレイアウトの再計算を必要とするかを返します。
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
	HasBeenLaidOut() bool // ウィジェットが一度でもレイアウトされたかを返します
	CurrentState() WidgetState
}

// EventHandler はイベント処理のためのインターフェースです
type EventHandler interface {
	AddEventHandler(eventType event.EventType, handler event.EventHandler)
	RemoveEventHandler(eventType event.EventType)
	HandleEvent(e *event.Event)
}

// AbsolutePositioner は、AbsoluteLayout内で希望の相対位置を指定できるウィジェットが実装するインターフェースです。
type AbsolutePositioner interface {
	SetRequestedPosition(x, y int)
	GetRequestedPosition() (x, y int)
}

// HitTester はヒットテストのためのインターフェースです
type HitTester interface {
	HitTest(x, y int) Widget
}

// --- Container Interface ---
// Containerは子Widgetを持つことができるWidgetです。
type Container interface {
	Widget // ContainerはWidgetのすべての振る舞いを継承します

	// --- 子要素管理 ---
	AddChild(child Widget)
	RemoveChild(child Widget)
	GetChildren() []Widget
}
