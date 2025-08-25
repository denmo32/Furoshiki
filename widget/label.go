package widget

import (
	"furoshiki/component"
	"furoshiki/theme"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。主にテキストを表示するためのシンプルなウィジェットです。
// 描画やヒットテストは、埋め込まれたcomponent.TextWidgetおよびcomponent.LayoutableWidgetの
// メソッドによって直接処理されるため、この型でメソッドをオーバーライドする必要はありません。
type Label struct {
	*component.TextWidget
}

// 【新規追加】NewLabelは、ラベルウィジェットの新しいインスタンスを生成し、初期化します。
// テキストとテーマに基づいたデフォルトスタイルが適用されます。
func NewLabel(text string) *Label {
	label := &Label{}
	label.TextWidget = component.NewTextWidget(text)
	label.Init(label) // LayoutableWidgetの初期化

	// テーマからデフォルトスタイルを取得して適用
	t := theme.GetCurrent()
	label.SetStyle(t.Label.Default)
	label.SetSize(100, 30) // TODO: Consider moving size to theme

	return label
}

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	// 【改善】新しいNewLabelコンストラクタを呼び出して、初期化ロジックを集約します。
	label := NewLabel("")

	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	return b
}

// Build は、設定に基づいて最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}