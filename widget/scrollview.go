package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
)

// ScrollView は、コンテンツをスクロール表示するためのコンテナウィジェットです。
// 汎用的な container.Container を埋め込むことで、基本的なコンテナ機能を継承します。
type ScrollView struct {
	*container.Container

	// スクロールビュー固有のプロパティ
	contentContainer component.Widget
	vScrollBar       component.ScrollBarWidget
	scrollY          float64
	contentHeight    int
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)
var _ container.Scroller = (*ScrollView)(nil)

// NewScrollView はScrollViewのインスタンスを生成します。
func NewScrollView() *ScrollView {
	sv := &ScrollView{}
	// 【エラー修正】 container.Containerの非公開フィールド(children)を構造体リテラルで
	// 初期化しようとしていたため、コンパイルエラーが発生していました。
	// このフィールドへの代入を削除します。childrenスライスはnilとして初期化されますが、
	// AddChildメソッド内のappend関数はnilスライスに対しても安全に動作するため問題ありません。
	c := &container.Container{}

	// 【改善】LayoutableWidgetを初期化し、`self`としてsv (*ScrollView) を渡してInitを呼び出します。
	// これにより、ヒットテストやイベント処理が正しくScrollViewのメソッドを呼び出すようになります。
	c.LayoutableWidget = component.NewLayoutableWidget()
	c.Init(sv)
	sv.Container = c
	sv.SetLayout(&layout.ScrollViewLayout{})

	vScrollBar, _ := NewScrollBarBuilder().Build()
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar)
	sv.SetClipsChildren(true)

	return sv
}

// 【重要修正】AddChild をオーバーライドします。
//
// このメソッドは、埋め込まれた container.Container の AddChild の振る舞いを拡張します。
// container.Container の AddChild は、子の親を「コンテナ自身」に設定します。
// しかし ScrollView の場合、子の親は「ScrollView インスタンス (sv)」であるべきです。
//
// このオーバーライドにより、子の親ポインタが正しく設定され、イベントのバブリングが
// この ScrollView のカスタム HandleEvent メソッドを正しく呼び出すようになります。
func (sv *ScrollView) AddChild(child component.Widget) {
	// 埋め込まれたコンテナのロジックを呼び出して、子リストへの追加や
	// 古い親からのデタッチ処理を行います。
	sv.Container.AddChild(child)

	// 親ポインタが ScrollView 自身を指しているか確認し、異なっていれば修正します。
	if child != nil && child.GetParent() != sv {
		child.SetParent(sv)
	}
}

// SetContent は、スクロールさせるコンテンツコンテナを設定します。
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		sv.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content

	if content != nil {
		// 上記でオーバーライドした AddChild が呼ばれます。
		sv.AddChild(content)
		sv.AddChild(sv.vScrollBar)
	}
	sv.MarkDirty(true)
}

// Update は、埋め込み先の Update を呼び出す前に、
// 自身のレイアウト処理を正しい引数で実行するようオーバーライドします。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}
	if sv.IsDirty() {
		if sv.NeedsRelayout() {
			if l := sv.GetLayout(); l != nil {
				l.Layout(sv) // 引数として `sv` 自身を渡す
			}
		}
		sv.ClearDirty()
	}
	for _, child := range sv.GetChildren() {
		child.Update()
	}
}

// HandleEvent は、ScrollView 固有のスクロール処理と、
// 埋め込みコンテナの標準的なイベント処理（バブリング等）を両立させます。
func (sv *ScrollView) HandleEvent(e *event.Event) {
	// Step 1: このウィジェット固有のイベント処理を先に行います。
	if !e.Handled && e.Type == event.MouseScroll {
		scrollAmount := e.ScrollY * 20 // スクロール感度
		sv.scrollY -= scrollAmount
		sv.MarkDirty(true) // 再レイアウトを要求
		e.Handled = true   // イベントを処理済みとしてマーク
	}

	// Step 2: 埋め込まれたコンテナの標準的なイベント処理を呼び出します。
	// これにより、カスタムハンドラが実行されたり、このウィジェットが
	// 処理しなかったイベントが親へ正しくバブリングされたりします。
	sv.Container.HandleEvent(e)
}

// HitTest は、埋め込み先のメソッドをそのまま呼び出します。
func (sv *ScrollView) HitTest(x, y int) component.Widget {
	return sv.Container.HitTest(x, y)
}

// --- layout.ScrollViewer interface implementation ---
func (sv *ScrollView) GetContentContainer() component.Widget {
	return sv.contentContainer
}
func (sv *ScrollView) GetVScrollBar() component.ScrollBarWidget {
	return sv.vScrollBar
}
func (sv *ScrollView) GetScrollY() float64 {
	return sv.scrollY
}
func (sv *ScrollView) SetScrollY(y float64) {
	if sv.scrollY != y {
		sv.scrollY = y
		sv.MarkDirty(false)
	}
}
func (sv *ScrollView) SetContentHeight(h int) {
	sv.contentHeight = h
}

// --- container.Scroller interface implementation ---
func (sv *ScrollView) GetScrollOffset() (x, y int) {
	return 0, -int(sv.scrollY)
}