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

// newButtonは、ボタンウィジェットの新しいインスタンスを生成し、初期化します。
// NOTE: このコンストラクタは非公開になりました。ウィジェットの生成には
//       常にNewButtonBuilder()を使用してください。これにより、初期化漏れを防ぎます。
func newButton(text string) (*Button, error) {
	button := &Button{}
	button.TextWidget = component.NewTextWidget(text)
	if err := button.Init(button); err != nil {
		return nil, err
	}

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
	button.SetStyle(styles[component.StateNormal])
	button.SetSize(100, 40)

	return button, nil
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
	// Mixinから、現在の状態に基づいて適用すべきスタイルを取得します。
	// このメソッドは、該当する状態のスタイルがなければNormal状態のスタイルにフォールバックするため、
	// 毎フレームの不要なスタイルコピーを回避できます。
	styleToUse := b.GetActiveStyle(currentState)
	// 取得したスタイルでウィジェットを描画します。
	b.TextWidget.DrawWithStyle(screen, styleToUse)
}

// SetStyleForState は、指定された単一の状態のスタイルを設定します。
// InteractiveMixinの責務が状態マップの管理のみになったため、このメソッド内で
// スタイルのマージと、必要に応じた基本スタイルの更新を行います。
func (b *Button) SetStyleForState(state component.WidgetState, s style.Style) {
	// 1. Mixinに保存されている現在の状態のスタイルに、新しいスタイルをマージします。
	b.InteractiveMixin.SetStyleForState(state, s)

	// 2. もしNormal状態のスタイルが変更された場合は、ウィジェットの基本スタイルも同期させます。
	if state == component.StateNormal {
		// LayoutableWidget.SetStyleを直接呼ぶことで、このメソッド内で再度
		// StateStyles[StateNormal]を更新する冗長な処理を避けます。
		b.LayoutableWidget.SetStyle(b.InteractiveMixin.StateStyles[component.StateNormal])
	}
}

// StyleAllStates は、ボタンの全てのインタラクティブな状態に共通のスタイル変更を適用します。
// ここでも、Mixinの更新後に基本スタイルとの同期をこのメソッド内で行います。
func (b *Button) StyleAllStates(s style.Style) {
	// 1. Mixinが管理する全ての状態のスタイルを更新します。
	b.InteractiveMixin.SetAllStyles(s)
	// 2. 更新後の新しいNormalスタイルをウィジェットの基本スタイルとして設定します。
	b.LayoutableWidget.SetStyle(b.InteractiveMixin.StateStyles[component.StateNormal])
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
	button, err := newButton("")
	b := &ButtonBuilder{}
	// 自身(b)と構築対象のウィジェット(button)を渡して、埋め込んだビルダーを初期化します。
	b.Init(b, button)
	b.AddError(err)
	return b
}

// Build は、最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	// 埋め込んだ汎用ビルダーのBuildメソッドを呼び出します。
	return b.Builder.Build()
}