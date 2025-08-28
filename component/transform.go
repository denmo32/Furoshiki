package component

// Pointは2D座標を表します。
type Point struct {
	X, Y int
}

// Sizeは2Dの寸法を表します。
type Size struct {
	Width, Height int
}

// Transformは、ウィジェットの位置とサイズを管理します。
// データ保持に専念し、ロジックは最小限に留めます。
type Transform struct {
	position     Point
	size         Size
	requestedPos Point
}

// NewTransformは新しいTransformインスタンスを作成します。
func NewTransform() *Transform {
	return &Transform{}
}

// SetPositionは、位置を設定します。
func (t *Transform) SetPosition(x, y int) {
	t.position.X = x
	t.position.Y = y
}

// GetPositionは、現在の位置を返します。
func (t *Transform) GetPosition() (int, int) {
	return t.position.X, t.position.Y
}

// SetSizeは、サイズを設定します。
func (t *Transform) SetSize(width, height int) {
	t.size.Width = width
	t.size.Height = height
}

// GetSizeは、現在のサイズを返します。
func (t *Transform) GetSize() (int, int) {
	return t.size.Width, t.size.Height
}

// SetRequestedPosition は、レイアウトに対する希望の相対位置を設定します。
func (t *Transform) SetRequestedPosition(x, y int) {
	t.requestedPos.X = x
	t.requestedPos.Y = y
}

// GetRequestedPosition は、レイアウトに対する希望の相対位置を返します。
func (t *Transform) GetRequestedPosition() (int, int) {
	return t.requestedPos.X, t.requestedPos.Y
}
