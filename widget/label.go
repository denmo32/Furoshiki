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

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は新しいLabelBuilderを生成します。
func NewLabelBuilder() *LabelBuilder {
	// まずラベルインスタンスを作成
	label := &Label{}
	// 次に、ラベル自身をselfとして渡してTextWidgetを初期化
	label.TextWidget = component.NewTextWidget(label, "")

	// --- テーマからスタイルを取得 ---
	t := theme.GetCurrent()
	label.SetStyle(t.Label.Default)

	label.SetSize(100, 30) // TODO: Consider moving size to theme

	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	return b
}

// Build は、設定に基づいて最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build()
}
