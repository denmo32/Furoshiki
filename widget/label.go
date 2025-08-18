package widget

import (
	"image/color"

	"furoshiki/component"
	"furoshiki/style"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
// 主にテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*component.TextWidget
}

// HitTest は、指定された座標がラベルの領域内にあるかを判定します。
// [追加] component.LayoutableWidgetの基本的なテストを呼び出し、ヒットした場合は
// LayoutableWidgetではなく、具象型であるLabel自身を返します。
// これにより、イベントシステムが正しいウィジェットインスタンスを扱えるようになります。
func (l *Label) HitTest(x, y int) component.Widget {
	// 埋め込まれたLayoutableWidgetのHitTestを呼び出して、基本的な境界チェックを行います
	if l.LayoutableWidget.HitTest(x, y) != nil {
		// ヒットした場合、インターフェースを満たす具象型であるLabel自身(*l)を返します
		return l
	}
	return nil
}

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は、デフォルトのスタイルで初期化されたLabelBuilderを返します。
func NewLabelBuilder() *LabelBuilder {
	label := &Label{
		TextWidget: component.NewTextWidget(""),
	}
	label.SetSize(100, 30)
	defaultStyle := style.Style{
		Background: color.Transparent,
		TextColor:  color.Black,
		Padding:    style.Insets{Top: 2, Right: 5, Bottom: 2, Left: 5},
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