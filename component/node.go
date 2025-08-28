package component

// Nodeは、UIツリー内の階層構造（親子関係）を管理します。
// データ保持に専念し、ロジックは最小限に留めます。
type Node struct {
	parent   NodeOwner
	children []NodeOwner
	self     NodeOwner // 自分自身（NodeOwnerを実装した具象ウィジェット）への参照
}

// NewNodeは新しいNodeインスタンスを作成します。
func NewNode(self NodeOwner) *Node {
	if self == nil {
		// selfがnilの場合、プログラムが予期せぬ動作をする可能性があるため、
		// パニックを発生させて早期に問題を検出します。
		panic("NewNode: self cannot be nil")
	}
	return &Node{
		children: make([]NodeOwner, 0),
		self:     self,
	}
}

// GetNodeは、NodeOwnerインターフェースの実装です。
func (n *Node) GetNode() *Node {
	return n
}

// SetParentは、このノードの親を設定します。
func (n *Node) SetParent(parent NodeOwner) {
	n.parent = parent
}

// GetParentは、このノードの親を返します。
func (n *Node) GetParent() NodeOwner {
	return n.parent
}

// GetChildrenは、このノードの子のスライスを返します。
func (n *Node) GetChildren() []NodeOwner {
	return n.children
}

// AddChildは、子ノードを追加します。
// 既存の親から子をデタッチするロジックもここに含まれます。
func (n *Node) AddChild(child NodeOwner) {
	if child == nil || child.GetNode() == nil {
		return
	}

	if oldParent := child.GetNode().GetParent(); oldParent != nil {
		// UPDATE: 非公開メソッドの代わりに公開されたDetachメソッドを使用します。
		// これにより、子を移動させる際に、古い親から安全にデタッチできます。
		oldParent.GetNode().Detach(child)
	}

	child.GetNode().SetParent(n.self)
	n.children = append(n.children, child)
}

// RemoveChildは、子ノードを削除し、親ポインタをnilに設定します。
func (n *Node) RemoveChild(child NodeOwner) {
	if n.removeChild(child) {
		child.GetNode().SetParent(nil)
	}
}

// UPDATE: Detachは、子を親の管理下から外しますが、子の親ポインタは変更しません。
// これは、ウィジェットをツリーの別の場所に「移動」する際に使用されます。
// AddChildが内部的に呼び出します。
func (n *Node) Detach(child NodeOwner) bool {
	return n.removeChild(child)
}

// UPDATE: ClearChildrenは、すべての子要素をリストから削除します。
// これは、コンテナがクリーンアップされる際に、安全に子リストを空にするために使用されます。
func (n *Node) ClearChildren() {
	// 参照を解放するためにスライスをnilにします。
	n.children = nil
}

// removeChildは、子リストから指定された子を削除する内部ヘルパーです。
func (n *Node) removeChild(child NodeOwner) bool {
	if child == nil || child.GetNode() == nil {
		return false
	}
	for i, currentChild := range n.children {
		if currentChild == child {
			n.children = append(n.children[:i], n.children[i+1:]...)
			return true
		}
	}
	return false
}