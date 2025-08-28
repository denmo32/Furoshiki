package component

// VisibilityOwnerは、Visibilityコンポーネントを所有するオブジェクトが
// 満たすべきインターフェースを定義します。
type VisibilityOwner interface {
	DirtyManager
}

// Visibilityはウィジェットの可視状態を管理します。
// これにより、可視性に関連するロジックとデータをカプセル化します。
type Visibility struct {
	owner     VisibilityOwner
	isVisible bool
}

// NewVisibilityは、新しいVisibilityコンポーネントを生成します。
func NewVisibility(owner VisibilityOwner) *Visibility {
	return &Visibility{
		owner:     owner,
		isVisible: true, // ウィジェットはデフォルトで可視です
	}
}

// SetVisibleはウィジェットの可視性を設定します。
func (v *Visibility) SetVisible(visible bool) {
	if v.isVisible != visible {
		v.isVisible = visible
		// 可視性の変更はレイアウトに影響するため、再レイアウトを要求します。
		v.owner.MarkDirty(true)
	}
}

// IsVisibleはウィジェットが現在可視であるかどうかを返します。
func (v *Visibility) IsVisible() bool {
	return v.isVisible
}
