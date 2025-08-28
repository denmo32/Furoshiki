package component

// WidgetCoreは、ほとんどのウィジェットに共通するコンポーネントとロジックを集約する
// 基底構造体です。これを埋め込むことで、コードの重複をなくし、メンテナンス性を向上させます。
// 「継承より合成」の原則に基づき、振る舞いの実装を共通化する役割を担います。
type WidgetCore struct {
	*Node
	*Transform
	*LayoutProperties
	*Visibility
	*Dirty

	// minSizeは、ユーザーによって明示的に設定された最小サイズを保持します。
	// GetMinSizeの実装は、この値とコンテンツ固有の最小サイズを比較して最終的な値を返します。
	minSize Size
}

// NewWidgetCoreは、新しいWidgetCoreインスタンスを生成します。
// 各ウィジェットのコンストラクタから呼び出されることを想定しています。
func NewWidgetCore(self NodeOwner) *WidgetCore {
	// self (具象ウィジェット自身) を各コンポーネントに渡すことで、
	// 階層構造やダーティフラグの伝播を正しく機能させます。
	wc := &WidgetCore{
		Node:             NewNode(self),
		Transform:        NewTransform(),
		LayoutProperties: NewLayoutProperties(),
	}
	// DirtyとVisibilityはオーナーとしてDirtyManagerを要求するため、wc自身を渡します。
	wc.Dirty = NewDirty()
	wc.Visibility = NewVisibility(wc)
	return wc
}

// --- 共通インターフェースの実装 ---

// GetNodeはNodeOwnerインターフェースの実装です。
func (c *WidgetCore) GetNode() *Node { return c.Node }

// GetLayoutPropertiesはLayoutPropertiesOwnerインターフェースの実装です。
func (c *WidgetCore) GetLayoutProperties() *LayoutProperties { return c.LayoutProperties }

// MarkDirtyはDirtyManagerインターフェースの実装です。
// ダーティフラグを設定し、親ウィジェットへの伝播ロジックをカプセル化します。
func (c *WidgetCore) MarkDirty(relayout bool) {
	c.Dirty.MarkDirty(relayout)
	// レイアウトに影響する変更で、かつ自身がレイアウト境界でない場合、
	// 親コンテナにも再レイアウトが必要であることを伝播させます。
	if relayout && !c.IsLayoutBoundary() {
		if parent := c.GetParent(); parent != nil {
			if dm, ok := parent.(DirtyManager); ok {
				dm.MarkDirty(true)
			}
		}
	}
}

// SetPositionはPositionSetterインターフェースの実装です。
// 初めて位置が設定された際に、ウィジェットがレイアウト済みであることを記録します。
func (c *WidgetCore) SetPosition(x, y int) {
	if !c.HasBeenLaidOut() {
		c.SetLaidOut(true)
	}
	if currX, currY := c.GetPosition(); currX != x || currY != y {
		c.Transform.SetPosition(x, y)
		// 位置の変更は通常、再描画のみを要求します。
		c.MarkDirty(false)
	}
}

// SetSizeはSizeSetterインターフェースの実装です。
// サイズの変更はレイアウトに影響するため、再レイアウトを要求します。
func (c *WidgetCore) SetSize(width, height int) {
	if w, h := c.GetSize(); w != width || h != height {
		c.Transform.SetSize(width, height)
		c.MarkDirty(true)
	}
}

// SetMinSizeはMinSizeSetterインターフェースの実装です。
// ユーザー指定の最小サイズを設定します。
func (c *WidgetCore) SetMinSize(width, height int) {
	c.minSize.Width = width
	c.minSize.Height = height
	c.MarkDirty(true)
}

// GetMinSizeはMinSizeSetterインターフェースの実装です。
// ユーザー指定の最小サイズを返します。具象ウィジェットはこれをオーバーライドし、
// コンテンツサイズとこの値を比較して最終的な最小サイズを決定します。
func (c *WidgetCore) GetMinSize() (int, int) {
	return c.minSize.Width, c.minSize.Height
}

// SetRequestedPositionはAbsolutePositionerインターフェースの実装です。
func (c *WidgetCore) SetRequestedPosition(x, y int) {
	c.Transform.SetRequestedPosition(x, y)
	c.MarkDirty(true)
}

// GetRequestedPositionはAbsolutePositionerインターフェースの実装です。
func (c *WidgetCore) GetRequestedPosition() (int, int) {
	return c.Transform.GetRequestedPosition()
}