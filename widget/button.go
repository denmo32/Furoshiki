package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/theme"

	"github.com/hajimehoshi/ebiten/v2"
)

// Button は、クリック可能なUI要素です。
type Button struct {
	*component.TextWidget
	stateStyles map[component.WidgetState]style.Style
}

// NewButtonは、ボタンウィジェットの新しいインスタンスを生成し、初期化します。
func NewButton(text string) *Button {
	button := &Button{
		stateStyles: make(map[component.WidgetState]style.Style),
	}
	button.TextWidget = component.NewTextWidget(text)
	button.Init(button) // LayoutableWidgetの初期化

	t := theme.GetCurrent()
	button.stateStyles[component.StateNormal] = t.Button.Normal.DeepCopy()
	button.stateStyles[component.StateHovered] = t.Button.Hovered.DeepCopy()
	button.stateStyles[component.StatePressed] = t.Button.Pressed.DeepCopy()
	button.stateStyles[component.StateDisabled] = t.Button.Disabled.DeepCopy()

	button.SetStyle(t.Button.Normal)
	button.SetSize(100, 40)

	return button
}

// Draw はButtonを描画します。現在の状態に応じたスタイルを選択し、描画を委譲します。
func (b *Button) Draw(screen *ebiten.Image) {
	currentState := b.LayoutableWidget.CurrentState()
	styleToUse := b.stateStyles[currentState]
	b.TextWidget.DrawWithStyle(screen, styleToUse)
}

// --- ButtonBuilder ---
type ButtonBuilder struct {
	Builder[*ButtonBuilder, *Button]
}

// NewButtonBuilder は新しいButtonBuilderを生成します。
func NewButtonBuilder() *ButtonBuilder {
	button := NewButton("")
	b := &ButtonBuilder{}
	b.Init(b, button)
	return b
}

// SetStyleForState は、指定された単一の状態のスタイルを設定します。
func (b *ButtonBuilder) SetStyleForState(state component.WidgetState, s style.Style) *ButtonBuilder {
	baseStyle := b.Widget.stateStyles[state]
	b.Widget.stateStyles[state] = style.Merge(baseStyle, s)

	// Normal状態のスタイルはレイアウト計算の基準となるため、ウィジェットの基本スタイルも更新します。
	if state == component.StateNormal {
		b.Widget.SetStyle(b.Widget.stateStyles[component.StateNormal])
	}
	return b
}

// Style は、ボタンの通常時（Normal状態）のスタイルを設定します。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateNormal, s)
}

// StyleAllStates は、ボタンの全てのインタラクティブな状態に共通のスタイル変更を適用します。
func (b *ButtonBuilder) StyleAllStates(s style.Style) *ButtonBuilder {
	for state, baseStyle := range b.Widget.stateStyles {
		b.Widget.stateStyles[state] = style.Merge(baseStyle, s)
	}
	b.Widget.SetStyle(b.Widget.stateStyles[component.StateNormal])
	return b
}

// HoverStyle は、ホバー時のスタイルを個別に設定します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateHovered, s)
}

// PressedStyle は、押下時のスタイルを個別に設定します。
func (b *ButtonBuilder) PressedStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StatePressed, s)
}

// DisabledStyle は、無効時のスタイルを個別に設定します。
func (b *ButtonBuilder) DisabledStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateDisabled, s)
}

// Build は、最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	return b.Builder.Build()
}
