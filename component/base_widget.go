package component

import (
	"furoshiki/event"
	"furoshiki/style"
)

// LayoutableWidgetは、Widgetインターフェースの基本的な実装を提供する構造体です。
// 他の具体的なウィジェット（Button, Labelなど）は、この構造体を埋め込むことで、
// 位置、サイズ、スタイル、イベント処理などの共通機能を利用します。
type LayoutableWidget struct {
	// --- Position & Size ---
	position           position
	size               size
	minSize            size
	requestedPos       position
	// 【改善】コンテンツの最小サイズを計算する関数。
	// これにより、最小サイズ決定ロジック（コンテンツサイズとユーザー設定サイズの比較）を
	// LayoutableWidgetに集約し、TextWidgetのような具象ウィジェットでのコードの重複を避けます。
	contentMinSizeFunc func() (width, height int)

	// --- Layout & Style ---
	layout layoutProperties
	style  style.Style

	// --- State ---
	state widgetState

	// --- Hierarchy & Events ---
	hierarchy     hierarchy
	eventHandlers map[event.EventType]event.EventHandler
	self          Widget
}

// position はウィジェットの位置情報を保持します
type position struct {
	x, y int
}

// size はウィジェットのサイズ情報を保持します
type size struct {
	width, height int
}

// layoutProperties はレイアウト関連のプロパティを保持します
type layoutProperties struct {
	flex             int
	relayoutBoundary bool
}

// dirtyLevel はウィジェットのダーティ状態のレベルを示します。
// これにより、再描画のみが必要か、レイアウトの再計算まで必要かを効率的に管理します。
type dirtyLevel int

const (
	// clean はウィジェットがダーティでないことを示します。
	clean dirtyLevel = iota
	// redrawDirty はウィジェットの再描画のみが必要なことを示します（例: ホバー状態の変化）。
	redrawDirty
	// relayoutDirty はウィジェットのレイアウト再計算と再描画が必要なことを示します（例: サイズの変更）。
	relayoutDirty
)

// widgetState はウィジェットの状態を保持します
type widgetState struct {
	dirtyLevel     dirtyLevel // dirty と relayoutDirty を置き換える新しいフィールド
	isHovered      bool
	isPressed      bool
	isVisible      bool
	isDisabled     bool
	hasBeenLaidOut bool // レイアウトが一度でも実行されたかを追跡するフラグ
}

// hierarchy はウィジェットの階層構造情報を保持します
type hierarchy struct {
	parent Container
}

// 【改善】NewLayoutableWidget は、self引数を取らずに LayoutableWidget を初期化します。
// self参照は、後からInitメソッドを呼び出して設定します。
func NewLayoutableWidget() *LayoutableWidget {
	return &LayoutableWidget{
		// isVisibleはデフォルトでtrue、hasBeenLaidOutはレイアウト計算が行われるまでfalseで初期化します。
		// dirtyLevelはデフォルトでcleanです。
		state:         widgetState{isVisible: true, hasBeenLaidOut: false, dirtyLevel: clean},
		eventHandlers: make(map[event.EventType]event.EventHandler),
	}
}

// 【改善】Initは、LayoutableWidgetが埋め込まれる具象ウィジェットへの参照(self)を
// 安全に設定するためのメソッドです。これにより、コンストラクタのシグネチャがシンプルになり、
// ウィジェットの初期化手順が統一されます。
func (w *LayoutableWidget) Init(self Widget) {
	if self == nil {
		// selfがnilの場合、プログラムが予期せぬ動作をする可能性があるため、パニックを発生させます。
		// これは、ウィジェットのコンストラクタが常に正しいインスタンスを渡すことを強制する設計上の決定です。
		panic("LayoutableWidget.Init: self cannot be nil")
	}
	if w.self != nil {
		// 既に初期化されている場合に再度Initを呼び出すのは、意図しない使われ方である可能性が高いため、
		// 安全のためにパニックさせます。
		panic("LayoutableWidget.Init: widget has already been initialized")
	}
	w.self = self
}