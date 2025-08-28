package component

import "furoshiki/style"

// AppearanceOwnerは、Appearanceコンポーネントを所有するオブジェクトが
// 満たすべきインターフェースを定義します。
// これにより、Appearanceは自身のオーナーのダーティ状態を更新できます。
type AppearanceOwner interface {
	DirtyManager
}

// Appearanceは、StyleManagerを利用してウィジェットの視覚スタイルを管理します。
// これにより、スタイリングに関連するロジックとデータをカプセル化します。
type Appearance struct {
	styleManager *StyleManager
}

// NewAppearanceは、新しいAppearanceコンポーネントを生成します。
// オーナーへの参照を受け取り、スタイル変更時にダーティフラグを立てられるようにします。
func NewAppearance(owner AppearanceOwner) *Appearance {
	return &Appearance{
		styleManager: NewStyleManager(owner),
	}
}

// SetStyleはウィジェットの基本スタイルを設定します。
func (a *Appearance) SetStyle(s style.Style) {
	a.styleManager.SetBaseStyle(s)
}

// GetStyleはウィジェットの現在の基本スタイルの安全なコピー（ディープコピー）を返します。
func (a *Appearance) GetStyle() style.Style {
	return a.styleManager.GetBaseStyle()
}

// ReadOnlyStyleは、パフォーマンスが重要な場面のために、ウィジェットの現在の基本スタイルを
// コピーせずに返します。返されたスタイルは変更してはいけません。
func (a *Appearance) ReadOnlyStyle() style.Style {
	return a.styleManager.ReadOnlyBaseStyle()
}

// SetStyleForStateは、特定のインタラクティブ状態に対応するスタイルを、既存のスタイルにマージします。
func (a *Appearance) SetStyleForState(state WidgetState, s style.Style) {
	a.styleManager.SetStyleForState(state, s)
}

// GetStyleForStateは、指定された状態に適用すべき最終的なスタイルを計算して返します。
// 内部でキャッシュを利用するため、効率的にスタイルを取得できます。
func (a *Appearance) GetStyleForState(state WidgetState) style.Style {
	return a.styleManager.GetStyleForState(state)
}
