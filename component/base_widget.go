package component

import (
	"errors"
	"furoshiki/event"
	"furoshiki/style"
)

// LayoutableWidgetは、Widgetインターフェースの基本的な実装を提供する構造体です。
// 他の具体的なウィジェット（Button, Labelなど）は、この構造体を埋め込むことで、
// 位置、サイズ、スタイル、イベント処理などの共通機能を利用します。
type LayoutableWidget struct {
	// --- Position & Size ---
	position     position
	size         size
	minSize      size
	requestedPos position
	// contentMinSizeFunc は、コンテンツ（テキストなど）が要求する最小サイズを計算する関数です。
	// これにより、最小サイズ決定ロ-ジック（コンテンツ固有サイズ vs ユーザー設定サイズ）を
	// LayoutableWidgetに集約し、TextWidgetのような具象ウィジェットでのコード重複を避けます。
	contentMinSizeFunc func() (width, height int)

	// --- Layout & Style ---
	layout       layoutProperties
	// NOTE: 内部実装を隠蔽し、コンポーネントのカプセル化を強化するために非公開に戻しました。
	//       スタイル操作は `SetStyle`, `SetStyleForState` などの公開メソッド経由で行います。
	styleManager *StyleManager // NOTE: フィールドを非公開に (StyleManager -> styleManager)

	// --- State ---
	state widgetState

	// --- Hierarchy & Events ---
	hierarchy hierarchy
	// NOTE: イベントハンドラを複数登録できるよう、型をハンドラのスライスに変更しました。
	eventHandlers map[event.EventType][]event.EventHandler
	// self は、このLayoutableWidgetを埋め込んでいる具象ウィジェット自身への参照です。
	// これにより、HitTestのようなメソッドが、具体的な型（*Button, *Labelなど）を返すことができます。
	self Widget
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
	// layoutData は、特定のレイアウトシステムが必要とする追加情報を格納するための汎用フィールドです。
	// 例えば、AdvancedGridLayoutはここにウィジェットの行、列、スパン情報を格納します。
	layoutData any
}

// dirtyLevel はウィジェットのダーティ状態のレベルを示します。
// これにより、再描画のみが必要か、レイアウトの再計算まで必要かを効率的に管理します。
type dirtyLevel int

const (
	// levelClean はウィジェットがダーティでないことを示します。
	levelClean dirtyLevel = iota
	// levelRedrawDirty はウィジェットの再描画のみが必要なことを示します（例: ホバー状態の変化）。
	levelRedrawDirty
	// levelRelayoutDirty はウィジェットのレイアウト再計算と再描画が必要なことを示します（例: サイズの変更）。
	levelRelayoutDirty
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

// NewLayoutableWidget は、LayoutableWidget を初期化します。
// この時点ではself参照は未設定です。
func NewLayoutableWidget() *LayoutableWidget {
	// NOTE: コンストラクタ内でStyleManagerを初期化し、自身への参照を渡します。
	// これにより、StyleManagerはスタイル変更時にこのウィジェットのMarkDirtyを呼び出せます。
	w := &LayoutableWidget{
		state:         widgetState{isVisible: true, hasBeenLaidOut: false, dirtyLevel: levelClean},
		eventHandlers: make(map[event.EventType][]event.EventHandler),
	}
	// NOTE: 非公開になったstyleManagerフィールドに設定します。
	w.styleManager = NewStyleManager(w)
	return w
}

// Initは、LayoutableWidgetが埋め込まれる具象ウィジェットへの参照(self)を
// 安全に設定するためのメソッドです。Goではコンストラクタ内で自分自身へのポインタを
// 取得することが難しいため、この2段階の初期化プロセスを採用しています。
// これにより、コンストラクタのシグネチャがシンプルになり、ウィジェットの初期化手順が統一されます。
// NOTE: 以前はpanicを使用していましたが、Goのエラーハンドリング慣習に従い、
//
//	errorを返すように変更しました。これにより、呼び出し側が適切に
//	エラーを処理できるようになります。
func (w *LayoutableWidget) Init(self Widget) error {
	if self == nil {
		// selfがnilの場合、プログラムが予期せぬ動作をする可能性があるため、エラーを返します。
		// これは、ウィジェットのコンストラクタが常に正しいインスタンスを渡すことを強制する設計上の決定です。
		return errors.New("LayoutableWidget.Init: self cannot be nil")
	}
	if w.self != nil {
		// 既に初期化されている場合に再度Initを呼び出すのは、意図しない使われ方である可能性が高いため、
		// 安全のためにエラーを返します。
		return errors.New("LayoutableWidget.Init: widget has already been initialized")
	}
	w.self = self
	return nil
}

// SetStyleForState は、特定のインタラクティブ状態に対応するスタイルを設定します。
// 具象ウィジェット（例: Button）が、テーマやユーザー設定に基づいて状態ごとのスタイルを
// 登録するために使用します。
func (w *LayoutableWidget) SetStyleForState(state WidgetState, s style.Style) {
	w.styleManager.SetStyleForState(state, s)
}

// GetStyleForState は、指定された状態に適用すべき最終的なスタイルを計算して返します。
// このメソッドは内部でキャッシュを利用するため、描画ループ内で効率的にスタイルを取得できます。
func (w *LayoutableWidget) GetStyleForState(state WidgetState) style.Style {
	return w.styleManager.GetStyleForState(state)
}