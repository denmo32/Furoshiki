package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"image/color"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。主にテキストを表示するためのシンプルなウィジェットです。
// 描画は、埋め込まれたcomponent.TextWidgetのDrawメソッドによって直接処理されます。
type Label struct {
	*component.TextWidget
}

// [削除] Label自身のDrawメソッドは、埋め込まれたTextWidgetのDrawメソッドと完全に同一のため削除しました。
// これによりコードの冗長性がなくなり、TextWidgetの描画ロジックが直接利用されます。

// [削除] HitTestメソッドは、component.LayoutableWidgetの汎用的な実装で十分なため、削除します。
// LayoutableWidgetは初期化時に具象ウィジェット(self)への参照を受け取り、
// HitTestが成功した際にその参照を返すため、具象型でのオーバーライドは不要です。

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

func NewLabelBuilder() *LabelBuilder {
	// まずラベルインスタンスを作成
	label := &Label{}
	// 次に、ラベル自身をselfとして渡してTextWidgetを初期化
	label.TextWidget = component.NewTextWidget(label, "")

	label.SetSize(100, 30)

	defaultStyle := style.Style{
		Background: style.PColor(color.Transparent),
		TextColor:  style.PColor(color.Black),
		Padding: style.PInsets(style.Insets{
			Top: 2, Right: 5, Bottom: 2, Left: 5,
		}),
	}
	label.SetStyle(defaultStyle)

	b := &LabelBuilder{}
	b.Builder.Init(b, label)
	return b
}

// Build は、設定に基づいて最終的なLabelを構築して返します。
func (b *LabelBuilder) Build() (*Label, error) {
	return b.Builder.Build("Label")
}