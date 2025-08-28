package component

// LayoutPropertiesOwnerは、LayoutProperties構造体を所有し、レイアウトシステムから
// プロパティを読み書きされるオブジェクトが実装すべきインターフェースです。
type LayoutPropertiesOwner interface {
	GetLayoutProperties() *LayoutProperties
}

// LayoutPropertiesは、レイアウトシステムがウィジェットを配置するために必要な
// プロパティを管理します。
type LayoutProperties struct {
	flex             int
	relayoutBoundary bool
	// layoutDataは、特定のレイアウトシステムが必要とする追加情報を格納するための汎用フィールドです。
	// 例えば、AdvancedGridLayoutはここにウィジェットの行、列、スパン情報を格納します。
	layoutData any
}

// NewLayoutPropertiesは新しいLayoutPropertiesインスタンスを作成します。
func NewLayoutProperties() *LayoutProperties {
	return &LayoutProperties{}
}

// GetLayoutPropertiesは、LayoutPropertiesOwnerインターフェースの実装です。
func (l *LayoutProperties) GetLayoutProperties() *LayoutProperties {
	return l
}

// SetFlexはFlexLayoutにおけるウィジェットの伸縮係数を設定します。
func (l *LayoutProperties) SetFlex(flex int) {
	if flex < 0 {
		flex = 0
	}
	l.flex = flex
}

// GetFlexはウィジェットの伸縮係数を返します。
func (l *LayoutProperties) GetFlex() int {
	return l.flex
}

// SetLayoutBoundaryは、このウィジェットをレイアウト計算の境界とするか設定します。
func (l *LayoutProperties) SetLayoutBoundary(isBoundary bool) {
	l.relayoutBoundary = isBoundary
}

// IsLayoutBoundaryは、このウィジェットがレイアウト計算の境界であるかを返します。
func (l *LayoutProperties) IsLayoutBoundary() bool {
	return l.relayoutBoundary
}

// SetLayoutDataはウィジェットにレイアウト固有のデータを設定します。
func (l *LayoutProperties) SetLayoutData(data any) {
	l.layoutData = data
}

// GetLayoutDataはウィジェットからレイアウト固有のデータを取得します。
func (l *LayoutProperties) GetLayoutData() any {
	return l.layoutData
}
