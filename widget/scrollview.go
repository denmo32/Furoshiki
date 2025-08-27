package widget

import (
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
)

// ScrollView は、コンテンツをスクロール表示するためのコンテナウィジェットです。
// このウィジェットは、基本的なウィジェット機能(LayoutableWidget)と、
// 子要素の管理とクリッピング描画を行う内部コンテナ(container.Container)を「合成」しています。
// この設計により、ScrollView独自の複雑な更新・描画ロジックと、汎用コンテナのロジックを
// 明確に分離し、コードの堅牢性を高めています。
type ScrollView struct {
	*component.LayoutableWidget
	container         *container.Container
	layout            layout.Layout
	contentContainer  component.Widget
	vScrollBar        component.ScrollBarWidget
	scrollY           float64
	contentHeight     int
	ScrollSensitivity float64
}

// コンパイル時にインターフェースの実装を検証します。
var _ component.Container = (*ScrollView)(nil)
var _ layout.ScrollViewer = (*ScrollView)(nil)
var _ container.Scroller = (*ScrollView)(nil)

// newScrollView はScrollViewのインスタンスを生成します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//       常にNewScrollViewBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newScrollView() (*ScrollView, error) {
	sv := &ScrollView{
		ScrollSensitivity: 20.0,
	}
	sv.LayoutableWidget = component.NewLayoutableWidget()

	// NOTE: ScrollViewの依存コンポーネントを、それぞれの公開コンストラクタ/ビルダー経由で生成します。
	internalContainer, err := container.NewContainer()
	if err != nil {
		return nil, err
	}
	sv.container = internalContainer
	sv.container.SetClipsChildren(true)

	if err := sv.Init(sv); err != nil {
		return nil, err
	}
	sv.container.SetParent(sv)
	sv.layout = &layout.ScrollViewLayout{}

	// NOTE: ScrollBarの生成にビルダーを使用することで、安全な初期化を保証します。
	vScrollBar, err := NewScrollBarBuilder().Build()
	if err != nil {
		return nil, err
	}
	sv.vScrollBar = vScrollBar
	sv.AddChild(vScrollBar)

	// 【提案1対応】HandleEventをオーバーライドする代わりに、専用のイベントハンドラを登録します。
	// これにより、イベント処理ロジックが一貫した方法で管理されます。
	sv.AddEventHandler(event.MouseScroll, sv.onMouseScroll)

	return sv, nil
}

// onMouseScroll は、MouseScrollイベントに応答してコンテンツをスクロールします。
// 【提案1対応】HandleEventのオーバーライドから移行した新しいイベントハンドラです。
func (sv *ScrollView) onMouseScroll(e *event.Event) event.Propagation {
	scrollAmount := e.ScrollY * sv.ScrollSensitivity
	sv.scrollY -= scrollAmount
	sv.MarkDirty(true)
	// ScrollViewがスクロールイベントを処理したので、親ウィジェットへの伝播を停止します。
	return event.StopPropagation
}

// SetContent は、スクロールさせるコンテンツコンテナを設定します。
func (sv *ScrollView) SetContent(content component.Widget) {
	if sv.contentContainer != nil {
		sv.RemoveChild(sv.contentContainer)
	}
	sv.contentContainer = content
	if content != nil {
		sv.AddChild(content)
		sv.AddChild(sv.vScrollBar) // スクロールバーが最前面に来るように再追加
	}
	sv.MarkDirty(true)
}

// SetStyle はScrollViewと、その描画を担当する内部コンテナの両方にスタイルを設定します。
// これにより、ビルダーで設定された枠線などが正しく描画されるようになります。
func (sv *ScrollView) SetStyle(s style.Style) {
	sv.LayoutableWidget.SetStyle(s)
	if sv.container != nil {
		sv.container.SetStyle(s)
	}
}

// Update はScrollViewの状態を更新します。
// 汎用コンテナのUpdateとは異なり、ScrollView専用のレイアウト計算を制御します。
// 【提案2対応】毎フレーム実行していた冗長な同期処理を削除しました。
// 同期はSetPositionとSetSizeが責務を担うように変更されています。
func (sv *ScrollView) Update() {
	if !sv.IsVisible() {
		return
	}

	// ScrollView自身が再レイアウトを要求されている場合のみ、専用のレイアウトを実行します。
	if sv.NeedsRelayout() {
		if sv.layout != nil {
			if err := sv.layout.Layout(sv); err != nil {
				// TODO: エラーハンドリング
			}
		}
	}

	// 内部コンテナのレイアウト処理は行わず、子要素のUpdateのみを再帰的に呼び出します。
	for _, child := range sv.container.GetChildren() {
		child.Update()
	}

	if sv.IsDirty() {
		sv.ClearDirty()
	}
}

// UPDATE: DrawメソッドのシグネチャをDrawInfoを受け取るように変更
// Draw はScrollViewを描画します。描画はクリッピング機能を持つ内部コンテナに完全に委譲します。
func (sv *ScrollView) Draw(info component.DrawInfo) {
	if !sv.IsVisible() {
		return
	}
	// 内部コンテナの描画にもDrawInfoを渡します。
	sv.container.Draw(info)
}

// SetPosition は自身の位置を設定し、その変更を内部コンテナにも伝播させます。
// 【提案2対応】Updateメソッドから状態同期ロジックを分離し、責務を明確化するために
// このオーバーライドメソッドが追加されました。
func (sv *ScrollView) SetPosition(x, y int) {
	// 1. 基底ウィジェット(自身)の位置を設定
	sv.LayoutableWidget.SetPosition(x, y)

	// 2. 内部コンテナの位置も同期
	if sv.container != nil {
		sv.container.SetPosition(x, y)
	}
}

// SetSize は自身のサイズを設定し、その変更を内部コンテナにも伝播させます。
// 【提案2対応】Updateメソッドから状態同期ロジックを分離し、責務を明確化するために
// このオーバーライドメソッドが追加されました。
func (sv *ScrollView) SetSize(width, height int) {
	// 1. 基底ウィジェット(自身)のサイズを設定
	sv.LayoutableWidget.SetSize(width, height)

	// 2. 内部コンテナのサイズも同期
	if sv.container != nil {
		sv.container.SetSize(width, height)
	}
}

// 【提案1対応】HandleEventのオーバーライドを削除。
// イベント処理はLayoutableWidgetのHandleEventに完全に委譲され、
// スクロール処理はコンストラクタで登録された専用ハンドラ(onMouseScroll)が担います。

// HitTest は、指定された座標がヒットするウィジェットを探します。
func (sv *ScrollView) HitTest(x, y int) component.Widget {
	if sv.LayoutableWidget.HitTest(x, y) != nil {
		// 自身がヒット範囲内であれば、次に内部コンテナの子要素をテストします。
		if target := sv.container.HitTest(x, y); target != nil {
			return target
		}
		return sv // 子にヒットしなければScrollView自身を返す
	}
	return nil
}

// --- メソッドの委譲 ---
func (sv *ScrollView) AddChild(child component.Widget)    { sv.container.AddChild(child) }
func (sv *ScrollView) RemoveChild(child component.Widget) { sv.container.RemoveChild(child) }
func (sv *ScrollView) GetChildren() []component.Widget    { return sv.container.GetChildren() }
func (sv *ScrollView) GetLayout() layout.Layout           { return sv.layout }
func (sv *ScrollView) SetLayout(l layout.Layout)          { sv.layout = l }
func (sv *ScrollView) GetPadding() layout.Insets          { return sv.container.GetPadding() }

// --- layout.ScrollViewer interface ---
func (sv *ScrollView) GetContentContainer() component.Widget    { return sv.contentContainer }
func (sv *ScrollView) GetVScrollBar() component.ScrollBarWidget { return sv.vScrollBar }
func (sv *ScrollView) GetScrollY() float64                      { return sv.scrollY }
func (sv *ScrollView) SetContentHeight(h int)                   { sv.contentHeight = h }
func (sv *ScrollView) SetScrollY(y float64) {
	if sv.scrollY != y {
		sv.scrollY = y
		sv.MarkDirty(false)
	}
}

// --- container.Scroller interface ---
func (sv *ScrollView) GetScrollOffset() (x, y int) { return 0, -int(sv.scrollY) }