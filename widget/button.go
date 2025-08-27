package widget

import (
	"furoshiki/component"
	"furoshiki/style"
	"furoshiki/theme"
)

// Button は、クリック可能なUI要素です。
// NOTE: component.InteractiveMixinの埋め込みは廃止され、状態ベースのスタイル管理は
// LayoutableWidgetが持つStyleManagerに完全に委譲されます。
type Button struct {
	*component.TextWidget
	// component.InteractiveMixin // 廃止
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

	// NOTE: テーマから各種状態のスタイルを取得し、StyleManagerに設定します。
	t := theme.GetCurrent()
	// 1. Normal状態のスタイルを、ウィジェットの「基本スタイル」として設定します。
	button.SetStyle(t.Button.Normal)
	// 2. 他の状態のスタイルを、状態固有のスタイルとして設定します。
	//    これらは描画時に基本スタイルとマージされます。
	// NOTE: [FIX] エクスポートされた StyleManager フィールドを使用します。
	button.StyleManager.SetStyleForState(component.StateHovered, t.Button.Hovered)
	button.StyleManager.SetStyleForState(component.StatePressed, t.Button.Pressed)
	button.StyleManager.SetStyleForState(component.StateDisabled, t.Button.Disabled)

	button.SetSize(100, 40)

	return button, nil
}

// SetStyle はウィジェットの基本スタイル(Normal状態の基礎)を設定します。
// このメソッドはLayoutableWidgetのStyleManagerのSetBaseStyleを呼び出します。
func (b *Button) SetStyle(s style.Style) {
	b.LayoutableWidget.SetStyle(s)
}

// UPDATE: DrawメソッドのシグネチャをDrawInfoを受け取るように変更
// Draw はButtonを描画します。現在の状態に応じたスタイルを選択し、描画を委譲します。
// NOTE: StyleManagerを利用して、現在の状態に最適なスタイルを効率的に取得します。
func (b *Button) Draw(info component.DrawInfo) {
	// 現在の状態（Normal, Hoveredなど）を取得します。
	currentState := b.LayoutableWidget.CurrentState()
	// StyleManagerから、現在の状態に基づいて適用すべきスタイルを取得します。
	// このメソッドは内部でキャッシュを利用するため、毎フレームの不要なスタイルコピーを回避できます。
	// NOTE: [FIX] エクスポートされた StyleManager フィールドを使用します。
	styleToUse := b.StyleManager.GetStyleForState(currentState)
	// 取得したスタイルでウィジェットを描画します。
	// UPDATE: DrawWithStyleにinfoを渡すように変更
	b.TextWidget.DrawWithStyle(info, styleToUse)
}

// SetStyleForState は、指定された単一の状態のスタイルを、既存のスタイルにマージします。
func (b *Button) SetStyleForState(state component.WidgetState, s style.Style) {
	// NOTE: [FIX] エクスポートされた StyleManager フィールドを使用します。
	b.StyleManager.SetStyleForState(state, s)
}

// StyleAllStates は、ボタンの全てのインタラクティブな状態に共通のスタイル変更を適用します。
func (b *Button) StyleAllStates(s style.Style) {
	// NOTE: このメソッドは、各状態スタイルにsをマージします。
	// 基本スタイルを変更したい場合はビルダーの .Style() を使用してください。
	// NOTE: [FIX] エクスポートされた StyleManager フィールドを使用します。
	b.StyleManager.SetStyleForState(component.StateNormal, s)
	b.StyleManager.SetStyleForState(component.StateHovered, s)
	b.StyleManager.SetStyleForState(component.StatePressed, s)
	b.StyleManager.SetStyleForState(component.StateDisabled, s)
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