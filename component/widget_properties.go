package component

import (
	"furoshiki/style"
)

// 【改善】このファイル内で頻繁に使用されるため、ヘルパー関数として定義します。
// これにより、Go 1.21未満の環境でも動作し、コードの意図が明確になります。
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SetPosition はウィジェットの絶対座標を設定します。
// 座標が変更された場合、再描画を要求します（レイアウトの再計算は不要）。
// このメソッドはレイアウトシステムによって呼び出されるため、ここで初めてウィジェットが
// 「レイアウト済み」であるとマークします。
func (w *LayoutableWidget) SetPosition(x, y int) {
	// 初めて位置が設定される（＝最初のレイアウトが完了した）際に、レイアウト済みフラグを立てます。
	if !w.state.hasBeenLaidOut {
		w.state.hasBeenLaidOut = true
	}

	if w.position.x != x || w.position.y != y {
		w.position.x = x
		w.position.y = y
		w.MarkDirty(false) // 位置変更は再描画のみが必要
	}
}

// GetPosition はウィジェットの絶対座標を返します。
func (w *LayoutableWidget) GetPosition() (x, y int) {
	return w.position.x, w.position.y
}

// SetSize はウィジェットのサイズを設定します。
// サイズが変更された場合、親コンテナにレイアウトの再計算を要求します。
func (w *LayoutableWidget) SetSize(width, height int) {
	// サイズが負の値の場合は処理しません。
	if width < 0 || height < 0 {
		return
	}

	if w.size.width != width || w.size.height != height {
		w.size.width = width
		w.size.height = height
		w.MarkDirty(true) // サイズ変更は再レイアウトが必要
	}
}

// GetSize はウィジェットのサイズを返します。
func (w *LayoutableWidget) GetSize() (width, height int) {
	return w.size.width, w.size.height
}

// SetMinSize はウィジェットの最小サイズを設定します。
// 最小サイズが変更された場合、親コンテナにレイアウトの再計算を要求します。
func (w *LayoutableWidget) SetMinSize(width, height int) {
	// 最小サイズが負の値の場合は処理しません。
	if width < 0 || height < 0 {
		return
	}

	if w.minSize.width != width || w.minSize.height != height {
		w.minSize.width = width
		w.minSize.height = height
		w.MarkDirty(true) // 最小サイズ変更は再レイアウトが必要
	}
}

// GetMinSize はウィジェットの最小サイズを返します。
// ユーザーが明示的に設定した最小サイズ(.MinSize())と、ウィジェットの
// コンテンツ（テキストなど）から計算される最小サイズのうち、大きい方を返します。
// コンテンツの最小サイズは、contentMinSizeFuncを通じて取得されます。
func (w *LayoutableWidget) GetMinSize() (width, height int) {
	// ユーザーが .MinSize() で明示的に設定した最小サイズを取得します。
	userMinWidth, userMinHeight := w.minSize.width, w.minSize.height

	// 【改善】コンテンツサイズを計算する関数が設定されていれば、それを呼び出して
	// ユーザー設定値と比較し、大きい方を最終的な最小サイズとします。
	if w.contentMinSizeFunc != nil {
		contentMinWidth, contentMinHeight := w.contentMinSizeFunc()
		// 【改善】ローカルで定義したmax関数を使い、各次元で大きい方の値を採用します。
		finalMinWidth := max(contentMinWidth, userMinWidth)
		finalMinHeight := max(contentMinHeight, userMinHeight)
		return finalMinWidth, finalMinHeight
	}

	// コンテンツサイズ計算関数がない場合は、ユーザー設定値のみを返します。
	return userMinWidth, userMinHeight
}

// SetRequestedPosition は、レイアウトに対する希望の相対位置を設定します。
//
// 重要: このメソッドは、親コンテナが `AbsoluteLayout` (主に `ui.ZStack` で作成) を
// 使用している場合にのみ有効です。`FlexLayout` (`VStack`, `HStack`) や `GridLayout` の
// 中にあるウィジェットに対してこのメソッドを使用しても、設定はレイアウトシステムによって
// 無視されるため効果はありません。これは意図された挙動です。
func (w *LayoutableWidget) SetRequestedPosition(x, y int) {
	if w.requestedPos.x != x || w.requestedPos.y != y {
		w.requestedPos.x = x
		w.requestedPos.y = y
		// 希望位置の変更は再レイアウトをトリガーすべきです。
		w.MarkDirty(true)
	}
}

// GetRequestedPosition は、レイアウトに対する希望の相対位置を返します。
func (w *LayoutableWidget) GetRequestedPosition() (int, int) {
	return w.requestedPos.x, w.requestedPos.y
}

// SetStyle はウィジェットのスタイルを設定します。
// スタイルの変更はレイアウトに影響する可能性があるため、再レイアウトを要求します。
func (w *LayoutableWidget) SetStyle(style style.Style) {
	w.style = style
	// スタイルの変更はパディングやフォントサイズに影響し、レイアウトが変わる可能性があるため、
	// 安全策として再レイアウトを要求します。
	w.MarkDirty(true)
}

// GetStyle はウィジェットの現在のスタイルの安全なコピーを返します。
// 内部の値もコピーされるため、この戻り値を変更しても元のウィジェットには影響しません。
func (w *LayoutableWidget) GetStyle() style.Style {
	return w.style.DeepCopy()
}

// SetFlex はFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (w *LayoutableWidget) SetFlex(flex int) {
	if flex < 0 {
		flex = 0 // flex値は0以上である必要があります。
	}
	if w.layout.flex != flex {
		w.layout.flex = flex
		w.MarkDirty(true)
	}
}

// GetFlex はウィジェットの伸縮係数を返します。
func (w *LayoutableWidget) GetFlex() int {
	return w.layout.flex
}

// 【改善】SetLayoutBoundary は、このウィジェットをレイアウト計算の境界とするか設定します。
// メソッド名を SetRelayoutBoundary から SetLayoutBoundary に変更して直感性を向上させました。
func (w *LayoutableWidget) SetLayoutBoundary(isBoundary bool) {
	if w.layout.relayoutBoundary != isBoundary {
		w.layout.relayoutBoundary = isBoundary
		// 境界設定が変更された場合は、再レイアウトが必要です。
		w.MarkDirty(true)
	}
}

// 【改善】SetRelayoutBoundary は、このウィジェットをレイアウト計算の境界とするか設定します。
// 後方互換性のために残します。内部的には SetLayoutBoundary を呼び出します。
func (w *LayoutableWidget) SetRelayoutBoundary(isBoundary bool) {
	// 新しいメソッドに委譲
	w.SetLayoutBoundary(isBoundary)
}

// SetParent はウィジェットの親コンテナを設定します。
func (w *LayoutableWidget) SetParent(parent Container) {
	w.hierarchy.parent = parent
}

// GetParent はウィジェットの親コンテナを返します。
func (w *LayoutableWidget) GetParent() Container {
	return w.hierarchy.parent
}
