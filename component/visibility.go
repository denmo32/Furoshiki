package component

// VisibilityOwnerは、Visibilityコンポーネントを所有するオブジェクトが
// 満たすべきインターフェースを定義します。
type VisibilityOwner interface {
	DirtyManager
}

// Visibilityはウィジェットの可視状態とレイアウト済み状態を管理します。
// これにより、ウィジェットが描画可能かを判断するためのロジックとデータをカプセル化します。
type Visibility struct {
	owner          VisibilityOwner
	isVisible      bool
	hasBeenLaidOut bool // UPDATE: レイアウト済み状態を管理するフラグを追加
}

// NewVisibilityは、新しいVisibilityコンポーネントを生成します。
func NewVisibility(owner VisibilityOwner) *Visibility {
	return &Visibility{
		owner:          owner,
		isVisible:      true, // ウィジェットはデフォルトで可視です
		hasBeenLaidOut: false, // 初期状態は未レイアウト
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

// UPDATE: SetLaidOutはウィジェットがレイアウトされたことを記録します。
// 状態が未レイアウトからレイアウト済みに変更された場合、再描画を要求します。
// これにより、初回描画が確実に行われます。
func (v *Visibility) SetLaidOut(laidOut bool) {
	if v.hasBeenLaidOut != laidOut {
		v.hasBeenLaidOut = laidOut
		if laidOut {
			// レイアウトが完了したフレームで描画されるように、再描画を要求します。
			// ここではレイアウトの再計算は不要なため `false` を指定します。
			v.owner.MarkDirty(false)
		}
	}
}

// UPDATE: HasBeenLaidOutはウィジェットが一度でもレイアウトされたかを返します。
func (v *Visibility) HasBeenLaidOut() bool {
	return v.hasBeenLaidOut
}