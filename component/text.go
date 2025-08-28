package component

// TextOwnerは、Textコンポーネントを所有するオブジェクトが
// 満たすべきインターフェースを定義します。
// これにより、Textは自身のオーナーのダーティ状態を更新できます。
type TextOwner interface {
	DirtyManager
}

// Textは、ウィジェットのテキスト内容と折り返しプロパティを管理します。
// これにより、テキストに関連するロジックとデータをカプセル化します。
type Text struct {
	owner    TextOwner
	text     string
	wrapText bool
}

// NewTextは、新しいTextコンポーネントを生成します。
func NewText(owner TextOwner, text string) *Text {
	return &Text{
		owner: owner,
		text:  text,
	}
}

// Textは現在のテキスト内容を返します。
func (t *Text) Text() string {
	return t.text
}

// SetTextはテキスト内容を設定し、必要であれば再レイアウトを要求します。
func (t *Text) SetText(text string) {
	if t.text != text {
		t.text = text
		t.owner.MarkDirty(true) // テキストの変更はレイアウトに影響する
	}
}

// WrapTextはテキストを折り返すかどうかのフラグを返します。
func (t *Text) WrapText() bool {
	return t.wrapText
}

// SetWrapTextはテキストの折り返しフラグを設定し、必要であれば再レイアウトを要求します。
func (t *Text) SetWrapText(wrap bool) {
	if t.wrapText != wrap {
		t.wrapText = wrap
		t.owner.MarkDirty(true) // 折り返し設定の変更はレイアウトに影響する
	}
}
