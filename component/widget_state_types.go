package component

// WidgetState は、ウィジェットが取りうるインタラクティブな状態を定義します。
type WidgetState int

const (
	// StateNormal は、ユーザー入力がないデフォルトの状態です。
	StateNormal WidgetState = iota
	// StateHovered は、マウスカーソルがウィジェット上にある状態です。
	StateHovered
	// StatePressed は、ウィジェットがクリックまたはタップされている最中の状態です。
	StatePressed
	// StateDisabled は、ウィジェットが無効化され、ユーザー入力を受け付けない状態です。
	StateDisabled
)
