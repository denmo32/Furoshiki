package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Label component ---
// LabelはTextWidgetを直接埋め込みます。Label固有のロジックは今のところありません。
// 主にテキストを表示するためのシンプルなウィジェットです。
type Label struct {
	*component.TextWidget
}

// [追加] Label自身のDrawメソッドを明示的に実装します。
// これにより、埋め込まれたTextWidgetのDrawメソッドへの暗黙的な依存がなくなり、
// 描画の振る舞いが明確かつ自己完結し、意図しないバグを防ぎます。
func (l *Label) Draw(screen *ebiten.Image) {
	if !l.IsVisible() {
		return
	}
	// 自身のプロパティを取得
	x, y := l.GetPosition()
	width, height := l.GetSize()
	style := l.GetStyle()
	text := l.Text()

	// 取得したスタイルを使って、背景とテキストをヘルパー関数で描画します。
	component.DrawStyledBackground(screen, x, y, width, height, style)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), style)
}

// [削除] HitTestメソッドは、component.LayoutableWidgetの汎用的な実装で十分なため、削除します。
// LayoutableWidgetは初期化時に具象ウィジェット(self)への参照を受け取り、
// HitTestが成功した際にその参照を返すため、具象型でのオーバーライドは不要です。

// --- LabelBuilder ---
// LabelBuilder は、Labelを安全かつ流れるように構築するためのビルダーです。
type LabelBuilder struct {
	Builder[*LabelBuilder, *Label]
}

// NewLabelBuilder は、デフォルトのスタイルで初期化されたLabelBuilderを返します。
// [修正] 初期化をself参照パターンに合わせ、スタイル設定をポインタ対応にします。
func NewLabelBuilder() *LabelBuilder {
	// まずラベルインスタンスを作成
	label := &Label{}
	// 次に、ラベル自身をselfとして渡してTextWidgetを初期化
	label.TextWidget = component.NewTextWidget(label, "")

	label.SetSize(100, 30)

	// [修正] 具象型の値からcolor.Color型の変数を作成し、そのアドレスを渡す
	bgColor := color.Color(color.Transparent)
	textColor := color.Color(color.Black)

	defaultStyle := style.Style{
		Background: &bgColor,
		TextColor:  &textColor,
		Padding: &style.Insets{
			Top: 2, Right: 5, Bottom: 2, Left: 5,
		},
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