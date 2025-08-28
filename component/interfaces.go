package component

import (
	"furoshiki/event"
	"furoshiki/style"

	"github.com/hajimehoshi/ebiten/v2"
)

// NodeOwnerは、Node構造体を所有し、階層構造の一部となることができる
// オブジェクトが実装すべきインターフェースです。
// これにより、具象ウィジェットの型に依存することなく、UIツリーを操作できます。
type NodeOwner interface {
	// GetNodeは、オブジェクトが所有するNodeへのポインタを返します。
	GetNode() *Node
}

// DrawInfoは、ウィジェットの描画に必要なすべてのコンテキストを保持します。
// これを導入することで、Drawメソッドに状態変更の副作用（例: SetPositionの呼び出し）を
// 持ち込む必要がなくなり、描画ロジックが純粋で予測可能になります。
type DrawInfo struct {
	Screen *ebiten.Image
	// 親から渡される描画オフセット。
	// ウィジェットは自身の絶対座標にこのオフセットを加算して描画します。
	OffsetX, OffsetY int
}

// --- Widget Interface ---
// Widgetは全てのUI要素の基本的な振る舞いを定義するインターフェースです。
// UIツリーの構成要素として最低限必要な機能（更新、描画、階層管理、ダーティ管理、ヒットテスト）を提供します。
// 【提案1】インターフェースのスリム化:
// 以前は全ての機能を一つの巨大なインターフェースにまとめていましたが、Goの思想に則り、
// より小さく、責務が明確なインターフェースに分割しました。
// レイアウトやイベント処理など、特定の機能は型アサーションを通じて利用されます。
type Widget interface {
	// --- ライフサイクル ---
	Update()
	// Drawメソッドは、描画コンテキスト(DrawInfo)を受け取り、ウィジェットを描画します。
	// これにより、クリッピング描画などで必要になるオフセット情報を安全に渡せます。
	Draw(info DrawInfo)
	Cleanup()

	// --- 階層構造 ---
	NodeOwner        // GetNode()を提供します
	HierarchyManager // SetParent()/GetParent()を提供します

	// --- 状態管理 ---
	DirtyManager

	// --- イベント処理の起点 ---
	HitTester
}

// UPDATE: 標準的なウィジェットが実装すべき能力を定義する複合インターフェースを追加。
// これにより、各ウィジェットの責務がコード上で明確になります。
// StandardWidgetは、レイアウト可能で、スタイルを持ち、インタラクティブな振る舞いを
// 持つことができる、標準的なウィジェットの契約を定義します。
type StandardWidget interface {
	Widget
	PositionSetter
	SizeSetter
	MinSizeSetter
	StyleGetterSetter
	LayoutPropertiesOwner
	AbsolutePositioner
	EventProcessor
	InteractiveState
}

// UPDATE: テキストを持つウィジェットのための複合インターフェースを追加。
type TextBasedWidget interface {
	StandardWidget
	TextOwner      // component.Text を所有
	HeightForWider // 幅に応じた高さを計算
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
// 【提案1】インターフェースのスリム化に伴い、Widgetインターフェースの他に
// 必要な振る舞いを明示的に埋め込みます。
type ScrollBarWidget interface {
	Widget
	InteractiveState // SetVisibleのため
	SizeSetter       // GetSize, SetSizeのため
	PositionSetter   // SetPositionのため
	SetRatios(contentRatio, scrollRatio float64)
}

// HierarchyManager は階層構造を管理するためのインターフェースです
type HierarchyManager interface {
	SetParent(parent NodeOwner)
	GetParent() NodeOwner
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
	// UPDATE: HasBeenLaidOutはVisibilityコンポーネントの責務となったため、このインターフェースから削除されました。
	// HasBeenLaidOut() bool

	// HasBeenLaidOut はウィジェットが一度でもレイアウトされたかを返します。
	// NOTE: このメソッドはVisibilityコンポーネントによって提供されるため、
	// 以前のようにこのインターフェースに含める必要は理論上ありませんが、
	// IsVisibleと密接に関連するため、利便性のためにここに残すことも検討できます。
	// 今回は、責務の分離を徹底するため、ここからは削除し、
	// 必要な場面では (component.VisibilityOwner) でアサーションすることを推奨します。
	// しかし、多くのウィジェットで必要となるため、利便性のために残します。
	HasBeenLaidOut() bool
	CurrentState() WidgetState
}

// EventProcessor はイベント処理のためのインターフェースです
// NOTE: 以前の EventHandler から名称変更。
type EventProcessor interface {
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
