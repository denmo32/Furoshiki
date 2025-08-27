package component

import "furoshiki/style"

// StyleManager はウィジェットのスタイルを状態ベースで管理する責務を担います。
// これにより、インタラクティブなウィジェットのスタイルロジックを共通化し、堅牢性を高めます。
// 内部でスタイルのマージ結果をキャッシュし、描画ループでのパフォーマンスを最適化します。
type StyleManager struct {
	baseStyle   style.Style
	stateStyles map[WidgetState]style.Style
	mergedCache map[WidgetState]style.Style
	owner       DirtyManager // オーナーウィジェットのダーティ状態を更新するための参照
}

// NewStyleManager は新しいStyleManagerインスタンスを生成します。
// オーナーウィジェットへの参照を受け取り、スタイル変更時に自動でダーティフラグを立てられるようにします。
func NewStyleManager(owner DirtyManager) *StyleManager {
	return &StyleManager{
		stateStyles: make(map[WidgetState]style.Style),
		mergedCache: make(map[WidgetState]style.Style),
		owner:       owner,
	}
}

// SetBaseStyle はウィジェットの基本スタイルを設定します。
// このスタイルは、他の全ての状態の基礎となります。
// 意図しない外部からの変更を防ぐため、内部にはスタイルのディープコピーを保持します。
func (sm *StyleManager) SetBaseStyle(s style.Style) {
	if sm.baseStyle.Equals(s) {
		return
	}
	sm.baseStyle = s.DeepCopy()
	sm.clearCache()
	sm.owner.MarkDirty(true) // スタイルの変更はレイアウトに影響する可能性がある
}

// GetBaseStyle はウィジェットの現在の基本スタイルの安全なコピーを返します。
func (sm *StyleManager) GetBaseStyle() style.Style {
	return sm.baseStyle.DeepCopy()
}

// ReadOnlyBaseStyle はウィジェットの現在の基本スタイルをコピーせずに返します。
// パフォーマンスが重要な描画ループなどでの使用を想定しており、返されたスタイルは変更してはいけません。
func (sm *StyleManager) ReadOnlyBaseStyle() style.Style {
	return sm.baseStyle
}

// SetStyleForState は、特定のインタラクティブ状態に対応するスタイルを、既存のスタイルにマージします。
// これにより、ビルダーパターンで .HoverStyle(...).Border(...) のような連続したスタイル設定が可能になります。
func (sm *StyleManager) SetStyleForState(state WidgetState, s style.Style) {
	existingStateStyle, _ := sm.stateStyles[state]
	mergedStyle := style.Merge(existingStateStyle, s)
	sm.stateStyles[state] = mergedStyle

	delete(sm.mergedCache, state) // 関連するキャッシュのみを破棄
	sm.owner.MarkDirty(true)
}

// GetStyleForState は、指定された状態に適用すべき最終的なスタイルを計算して返します。
// 最初に基本スタイルを適用し、その上に状態固有のスタイルをマージします。
// 計算結果はキャッシュされ、パフォーマンスを向上させます。
func (sm *StyleManager) GetStyleForState(state WidgetState) style.Style {
	if merged, ok := sm.mergedCache[state]; ok {
		return merged
	}

	finalStyle := sm.baseStyle.DeepCopy()
	if stateSpecificStyle, ok := sm.stateStyles[state]; ok {
		finalStyle = style.Merge(finalStyle, stateSpecificStyle)
	}

	sm.mergedCache[state] = finalStyle
	return finalStyle
}

// clearCache は全てのマージ済みスタイルキャッシュを破棄します。
// 基本スタイルが変更された際に呼び出されます。
func (sm *StyleManager) clearCache() {
	// NOTE: [FIX] マップの型を正しく `map[WidgetState]style.Style` に修正しました。
	sm.mergedCache = make(map[WidgetState]style.Style)
}