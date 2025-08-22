package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"furoshiki/theme"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
// Button は、クリック可能なUI要素です。TextWidgetを拡張し、状態に基づいたスタイル管理機能を追加します。
type Button struct {
	*component.TextWidget
	// stateStyles は、ボタンの各インタラクティブな状態に対応するスタイルを保持します。
	// Builderによってビルド時にすべての状態が設定されるため、Drawメソッドでのnilチェックは不要です。
	stateStyles map[component.WidgetState]style.Style
}

// Draw はButtonを描画します。
// ボタンの責務は、現在の状態に応じたスタイルを選択し、そのスタイルを使って
// 埋め込まれたTextWidgetの描画ロジックを呼び出すことです。
// これにより、描画コードの重複が排除され、関心の分離が促進されます。
func (b *Button) Draw(screen *ebiten.Image) {
	// 現在のインタラクティブな状態（Normal, Hovered, Pressed, Disabled）を取得します。
	currentState := b.LayoutableWidget.CurrentState()

	// 状態に対応するスタイルをマップから取得します。
	// Build()メソッドで全ての状態のスタイルが設定されることが保証されているため、
	// ここで存在チェックを行う必要はありません。
	styleToUse := b.stateStyles[currentState]

	// 埋め込まれたTextWidgetが提供する共通の描画メソッドを、
	// 選択したスタイルを渡して呼び出します。
	// IsVisibleやHasBeenLaidOutのチェックはこのメソッド内で行われます。
	b.TextWidget.DrawWithStyle(screen, styleToUse)
}

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
}

// NewButtonBuilder は新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	// ボタンインスタンスを作成
	button := &Button{
		stateStyles: make(map[component.WidgetState]style.Style),
	}
	// ボタン自身をselfとして渡してTextWidgetを初期化
	button.TextWidget = component.NewTextWidget(button, "")

	// --- テーマから各状態のデフォルトスタイルを取得し、設定します ---
	t := theme.GetCurrent()
	button.stateStyles[component.StateNormal] = t.Button.Normal.DeepCopy()
	button.stateStyles[component.StateHovered] = t.Button.Hovered.DeepCopy()
	button.stateStyles[component.StatePressed] = t.Button.Pressed.DeepCopy()
	button.stateStyles[component.StateDisabled] = t.Button.Disabled.DeepCopy()

	// デフォルトのスタイルとしてNormalを適用
	button.SetStyle(t.Button.Normal)
	button.SetSize(100, 40) // TODO: Consider moving size to theme

	b := &ButtonBuilder{}
	b.Init(b, button)
	return b
}

// SetStyleForState は、指定された単一の状態のスタイルを設定します。
// 既存のスタイルとマージして部分的な上書きを可能にします。
func (b *ButtonBuilder) SetStyleForState(state component.WidgetState, s style.Style) *ButtonBuilder {
	// 既存の状態スタイルをベースに、新しいスタイルをマージします。
	baseStyle, ok := b.Widget.stateStyles[state]
	if !ok {
		// 万が一ベーススタイルが存在しない場合は、Normal状態のスタイルをコピーして使用します。
		// 通常、コンストラクタで全状態が初期化されるため、このパスは通りません。
		baseStyle = b.Widget.stateStyles[component.StateNormal]
	}
	b.Widget.stateStyles[state] = style.Merge(baseStyle, s)
	return b
}

// OnClick は、ボタンがクリックされたときに実行されるイベントハンドラを設定します。
func (b *ButtonBuilder) OnClick(handler event.EventHandler) *ButtonBuilder {
	if handler != nil {
		b.Widget.AddEventHandler(event.EventClick, handler)
	}
	return b
}

// Style は、ボタンのすべてのインタラクティブな状態（Normal, Hovered, Pressed, Disabled）に
// 共通で適用される基本スタイルを設定します。
// 例えば、ここでTextColorを設定すると、すべての状態でテキストの色が変更されます。
// 特定の状態のスタイルを個別に変更したい場合は、HoverStyle()やPressedStyle()を使用してください。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	// 管理しているすべての状態スタイルに対して、渡されたスタイルをマージします。
	for state, baseStyle := range b.Widget.stateStyles {
		b.Widget.stateStyles[state] = style.Merge(baseStyle, s)
	}
	// 変更をウィジェットに反映させ、再レイアウト・再描画を要求します。
	b.Widget.MarkDirty(true)
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるとき（Hovered状態）のスタイルを個別に設定します。
// これは、Style()で設定された基本スタイルをさらに上書きするために使用します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateHovered, s)
}

// PressedStyle は、マウスボタンが押されている最中（Pressed状態）のスタイルを個別に設定します。
// これは、Style()やHoverStyle()で設定されたスタイルをさらに上書きするために使用します。
func (b *ButtonBuilder) PressedStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StatePressed, s)
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
// このメソッドは、各状態のスタイルが適切に設定されていることを保証します。
func (b *ButtonBuilder) Build() (*Button, error) {
	styles := b.Widget.stateStyles

	// ビルド時に、ユーザーによって明示的に設定されていない状態のスタイルを、
	// 合理的なデフォルト値で確実に埋めます。これにより、実行時のスタイル解決が不要になります。

	// Normalがなければテーマから取得します (通常はコンストラクタで設定済み)。
	normalStyle, ok := styles[component.StateNormal]
	if !ok {
		normalStyle = theme.GetCurrent().Button.Normal
		styles[component.StateNormal] = normalStyle
	}

	// Hoveredスタイルが設定されていなければ、Normalスタイルを継承します。
	if _, ok := styles[component.StateHovered]; !ok {
		styles[component.StateHovered] = normalStyle
	}

	// Pressedスタイルが設定されていなければ、Hoveredスタイルを継承します。
	// これにより、Hover -> Press の自然な視覚的変化が生まれます。
	if _, ok := styles[component.StatePressed]; !ok {
		styles[component.StatePressed] = styles[component.StateHovered]
	}

	// Disabledスタイルが設定されていなければ、Normalスタイルをベースに半透明にします。
	if _, ok := styles[component.StateDisabled]; !ok {
		// テーマのDisabledスタイルは既に半透明ですが、ユーザーがNormalスタイルのみを
		// カスタマイズした場合のフォールバックとして機能します。
		disabledStyle := style.Merge(normalStyle, style.Style{Opacity: style.PFloat64(0.5)})
		styles[component.StateDisabled] = disabledStyle
	}

	// 汎用ビルダーのBuildを呼び出して、最終的なエラーチェックなどを行います。
	return b.Builder.Build()
}
