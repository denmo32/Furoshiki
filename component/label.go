package component

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
// 主にテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*TextWidget
}