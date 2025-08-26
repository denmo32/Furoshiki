package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/theme"

	"github.com/hajimehoshi/ebiten/v2"
)

// Button は、クリック可能なUI要素です。
// component.InteractiveMixin を埋め込むことで、状態に基づいたスタイル管理を共通化します。
type Button struct {
	*component.TextWidget
	component.InteractiveMixin
}

// NewButtonは、ボタンウィジェットの新しいインスタンスを生成し、初期化します。
func NewButton(text string) *Button {
	button := &Button{}
	button.TextWidget = component.NewTextWidget(text)
	button.Init(button) // LayoutableWidgetの初期化

	// テーマから各種状態のスタイルを取得し、InteractiveMixinを初期化します。
	t := theme.GetCurrent()
	styles := map[component.WidgetState]style.Style{
		component.StateNormal:   t.Button.Normal.DeepCopy(),
		component.StateHovered:  t.Button.Hovered.DeepCopy(),
		component.StatePressed:  t.Button.Pressed.DeepCopy(),
		component.StateDisabled: t.Button.Disabled.DeepCopy(),
	}
	button.InitStyles(styles)

	// 通常時のスタイルをウィジェットの基本スタイルとして設定します。
	// この時点でSetStyleを呼ぶと、オーバーライドされたSetStyleがStateStyles[StateNormal]も更新します。
	button.SetStyle(styles[component.StateNormal])
	button.SetSize(100, 40)

	return button
}

// SetStyle はウィジェットの基本スタイルを設定します。
// インタラクティブウィジェットとして、基本スタイル(w.style)と
// InteractiveMixinが保持するNormal状態のスタイル(StateStyles[StateNormal])の両方を更新します。
// これにより、ビルダー等からSetStyleが呼ばれた際にスタイルの一貫性が保たれます。
func (b *Button) SetStyle(s style.Style) {
	// 1. 基底ウィジェットのスタイルを設定します。
	//    基底のSetStyleは変更検知を行い、不要な場合はダーティフラグを立てません。
	b.LayoutableWidget.SetStyle(s)

	// 2. InteractiveMixinの状態マップ内のNormalスタイルも更新します。
	//    これにより、GetActiveStyleが常に最新のNormalスタイルを参照できるようになります。
	//    DeepCopyを行い、b.styleとb.InteractiveMixin.StateStyles[component.StateNormal]が
	//    異なるメモリ領域を参照するようにし、意図しない副作用を防ぎます。
	b.InteractiveMixin.StateStyles[component.StateNormal] = s.DeepCopy()
}

// Draw はButtonを描画します。現在の状態に応じたスタイルを選択し、描画を委譲します。
// InteractiveMixinのGetActiveStyleを利用して、現在の状態に最適なスタイルを取得します。
func (b *Button) Draw(screen *ebiten.Image) {
	// 現在の状態（Normal, Hoveredなど）を取得します。
	currentState := b.LayoutableWidget.CurrentState()
	// Mixinから、現在の状態と基本スタイルに基づいて適用すべきスタイルを取得します。
	styleToUse := b.GetActiveStyle(currentState, b.GetStyle())
	// 取得したスタイルでウィジェットを描画します。
	b.TextWidget.DrawWithStyle(screen, styleToUse)
}

// SetStyleForState は、指定された単一の状態のスタイルを設定します。
// Normal状態のスタイルが変更された場合、ウィジェットの基本スタイルも更新します。
func (b *Button) SetStyleForState(state component.WidgetState, s style.Style) {
	// Normal状態のスタイルが変更された場合にウィジェットの基本スタイルを更新するためのコールバック。
	// b.SetStyle()ではなくb.LayoutableWidget.SetStyle()を呼ぶことで、
	// b.SetStyle()内で実行されるStateStyles[StateNormal]の更新との冗長性や潜在的な循環呼び出しを避けます。
	normalSetter := func(normalStyle style.Style) {
		b.LayoutableWidget.SetStyle(normalStyle)
	}
	b.InteractiveMixin.SetStyleForState(state, s, normalSetter)
}

// StyleAllStates は、ボタンの全てのインタラクティブな状態に共通のスタイル変更を適用します。
func (b *Button) StyleAllStates(s style.Style) {
	// normalSetterコールバックのロジックはSetStyleForStateと同様です。
	normalSetter := func(normalStyle style.Style) {
		b.LayoutableWidget.SetStyle(normalStyle)
	}
	b.InteractiveMixin.SetAllStyles(s, normalSetter)
}

// --- ButtonBuilder ---

// ButtonBuilder は、汎用の InteractiveTextBuilder を利用してButtonを構築します。
// これにより、状態ごとのスタイル設定（HoverStyle, PressedStyleなど）のロジックを再利用します。
type ButtonBuilder struct {
	// InteractiveTextBuilderを埋め込むことで、状態管理機能を持つテキストベースのウィジェットの
	// ビルダー機能を継承します。
	InteractiveTextBuilder[*ButtonBuilder, *Button]
}

// NewButtonBuilder は新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	button := NewButton("")
	b := &ButtonBuilder{}
	// 自身(b)と構築対象のウィジェット(button)を渡して、埋め込んだビルダーを初期化します。
	b.Init(b, button)
	return b
}

// Build は、最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	// 埋め込んだ汎用ビルダーのBuildメソッドを呼び出します。
	return b.Builder.Build()
}