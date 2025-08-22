package widget

import (
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/style"
	"furoshiki/theme"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
// Button は、クリック可能なUI要素です。TextWidgetを拡張し、状態に基づいたスタイル管理機能を追加します。
type Button struct {
	*component.TextWidget
	stateStyles  map[component.WidgetState]style.Style
	currentState component.WidgetState
	isPressed    bool // マウスボタンが現在このウィジェット上で押されているか
}

// Update はボタンの状態を更新します。LayoutableWidgetのUpdateをオーバーライドします。
func (b *Button) Update() {
	if b.IsDisabled() {
		if b.currentState != component.StateDisabled {
			b.currentState = component.StateDisabled
			b.MarkDirty(false)
		}
		return
	}

	newState := component.StateNormal
	if b.IsHovered() {
		newState = component.StateHovered
	}
	if b.isPressed {
		newState = component.StatePressed
	}

	if b.currentState != newState {
		b.currentState = newState
		b.MarkDirty(false) // スタイル変更は再描画のみ要求
	}
}

// HandleEvent はイベントを処理します。LayoutableWidgetのHandleEventをオーバーライドします。
func (b *Button) HandleEvent(e event.Event) {
	switch e.Type {
	case event.MouseDown:
		b.isPressed = true
	case event.MouseUp:
		b.isPressed = false
	}

	// 元のイベントハンドラ呼び出しロジックも実行
	b.TextWidget.HandleEvent(e)
}

// Draw はButtonを描画します。現在の状態に応じたスタイルを適用します。
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	x, y := b.GetPosition()
	width, height := b.GetSize()
	text := b.Text()

	// 現在の状態に最適なスタイルを選択
	// Pressed -> Hovered -> Normal の優先順位でフォールバック
	styleToUse, ok := b.stateStyles[b.currentState]
	if !ok {
		styleToUse, ok = b.stateStyles[component.StateHovered]
		if !ok {
			styleToUse = b.stateStyles[component.StateNormal]
		}
	}

	// 選択したスタイルで背景とテキストを描画
	component.DrawStyledBackground(screen, x, y, width, height, styleToUse)
	component.DrawAlignedText(screen, text, image.Rect(x, y, x+width, y+height), styleToUse)
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

	// --- テーマからスタイルを取得 ---
	t := theme.GetCurrent()
	button.stateStyles[component.StateNormal] = t.Button.Normal
	button.stateStyles[component.StateHovered] = t.Button.Hovered
	button.stateStyles[component.StatePressed] = t.Button.Pressed
	button.stateStyles[component.StateDisabled] = t.Button.Disabled

	// デフォルトのスタイルとしてNormalを適用
	button.SetStyle(t.Button.Normal)
	button.SetSize(100, 40) // TODO: Consider moving size to theme

	b := &ButtonBuilder{}
	b.Builder.Init(b, button)
	return b
}

// SetStyleForState は、指定された状態のスタイルを設定します。
// これが新しい、主要なスタイル設定メソッドとなります。
func (b *ButtonBuilder) SetStyleForState(state component.WidgetState, s style.Style) *ButtonBuilder {
	// 既存のスタイルとマージして設定
	baseStyle, ok := b.Widget.stateStyles[state]
	if !ok {
		// ベーススタイルが存在しない場合はNormalからコピー
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

// Style はボタンの基本（Normal状態）のスタイルを設定します。
// 下位互換性のために残されていますが、内部ではSetStyleForStateを呼び出します。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	// Normal状態のスタイルを上書き
	b.Widget.stateStyles[component.StateNormal] = style.Merge(b.Widget.stateStyles[component.StateNormal], s)
	// 変更を伝播させる
	b.Widget.MarkDirty(true)
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるときのスタイルを設定します。
// 下位互換性のために残されていますが、内部ではSetStyleForStateを呼び出します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateHovered, s)
}

// PressedStyle は、マウスボタンが押されている最中のスタイルを設定します。
func (b *ButtonBuilder) PressedStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StatePressed, s)
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	// Build時に、設定されていない状態のスタイルをフォールバックで埋めます。
	normalStyle, ok := b.Widget.stateStyles[component.StateNormal]
	if !ok {
		// テーマが空の場合など、万が一Normalが設定されていなかった場合の安全策
		normalStyle = theme.GetCurrent().Button.Normal
		b.Widget.stateStyles[component.StateNormal] = normalStyle
	}

	// HoveredがなければNormalをコピー
	if _, ok := b.Widget.stateStyles[component.StateHovered]; !ok {
		b.Widget.stateStyles[component.StateHovered] = normalStyle
	}
	// PressedがなければHoveredをコピー
	if _, ok := b.Widget.stateStyles[component.StatePressed]; !ok {
		b.Widget.stateStyles[component.StatePressed] = b.Widget.stateStyles[component.StateHovered]
	}

	// 汎用ビルダーのBuildを呼び出して、最終的なエラーチェックなどを行います。
	return b.Builder.Build("Button")
}