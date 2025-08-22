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
	// NewButtonBuilder()で全ての状態のスタイルが設定されることが保証されているため、
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
	// この時点で全てのインタラクティブな状態に対応するスタイルがマップに設定されるため、
	// 実行時にスタイルが見つからないという状況を防ぎます。
	t := theme.GetCurrent()
	button.stateStyles[component.StateNormal] = t.Button.Normal.DeepCopy()
	button.stateStyles[component.StateHovered] = t.Button.Hovered.DeepCopy()
	button.stateStyles[component.StatePressed] = t.Button.Pressed.DeepCopy()
	button.stateStyles[component.StateDisabled] = t.Button.Disabled.DeepCopy()

	// デフォルトのスタイルとしてNormalを適用
	// これにより、レイアウト計算などで使用される基本スタイルが設定されます。
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
		// 通常、コンストラクタで全状態が初期化されるため、このパスは通りません。
		// 安全策として、万が一ベースが存在しない場合はNormal状態のスタイルをコピーして使用します。
		baseStyle = b.Widget.stateStyles[component.StateNormal]
	}
	b.Widget.stateStyles[state] = style.Merge(baseStyle, s)

	// もし更新したのがNormal状態のスタイルであれば、ウィジェットの基本スタイルも更新します。
	// これにより、レイアウト計算が常に最新のNormal状態のスタイルプロパティ（パディング等）を
	// 使用することが保証されます。
	if state == component.StateNormal {
		b.Widget.SetStyle(b.Widget.stateStyles[component.StateNormal])
	}
	return b
}

// OnClick は、ボタンがクリックされたときに実行されるイベントハンドラを設定します。
func (b *ButtonBuilder) OnClick(handler event.EventHandler) *ButtonBuilder {
	if handler != nil {
		b.Widget.AddEventHandler(event.EventClick, handler)
	}
	return b
}

// Style は、ボタンの通常時（Normal状態）のスタイルを設定します。
// これは最も一般的に使用されるスタイル設定メソッドです。
// レイアウトに影響を与えるプロパティ（Padding, Marginなど）を変更する場合は、このメソッドを使用してください。
//
// 全ての状態に共通のスタイル変更を適用したい場合は `StyleAllStates()` を、
// 特定の状態のスタイルのみを変更したい場合は `HoverStyle()` や `PressedStyle()` を使用してください。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	// Normal状態のスタイルを更新します。内部でウィジェットの基本スタイルも更新されます。
	return b.SetStyleForState(component.StateNormal, s)
}

// StyleAllStates は、ボタンの全てのインタラクティブな状態（Normal, Hovered, Pressed, Disabled）に
// 共通のスタイル変更を適用します。
// 例えば、`StyleAllStates(style.Style{BorderRadius: style.PFloat32(10)})` を呼び出すと、
// 全ての状態の角丸が変更されますが、背景色など他のプロパティは各状態の既存の設定が維持されます。
func (b *ButtonBuilder) StyleAllStates(s style.Style) *ButtonBuilder {
	// 管理しているすべての状態スタイルに対して、渡されたスタイルをマージします。
	for state, baseStyle := range b.Widget.stateStyles {
		b.Widget.stateStyles[state] = style.Merge(baseStyle, s)
	}
	// Normal状態のスタイルが変更されたため、ウィジェットの基本スタイルも更新し、
	// レイアウトへの変更を反映させます。
	b.Widget.SetStyle(b.Widget.stateStyles[component.StateNormal])
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるとき（Hovered状態）のスタイルを個別に設定します。
// これは、Style()で設定されたスタイルをさらに上書きするために使用します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateHovered, s)
}

// PressedStyle は、マウスボタンが押されている最中（Pressed状態）のスタイルを個別に設定します。
// これは、Style()やHoverStyle()で設定されたスタイルをさらに上書きするために使用します。
func (b *ButtonBuilder) PressedStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StatePressed, s)
}

// DisabledStyle は、ボタンが無効化されているとき（Disabled状態）のスタイルを個別に設定します。
func (b *ButtonBuilder) DisabledStyle(s style.Style) *ButtonBuilder {
	return b.SetStyleForState(component.StateDisabled, s)
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
// Buttonのスタイルは、コンストラクタ(NewButtonBuilder)でテーマに基づいて
// 全てのインタラクティブな状態に対して初期化済みです。
// その後、Style(), HoverStyle(), PressedStyle() などのメソッドを通じて個別にカスタマイズできます。
func (b *ButtonBuilder) Build() (*Button, error) {
	// 以前のバージョンに存在した、各状態のスタイルに対するフォールバック処理は、
	// NewButtonBuilderで全ての状態がテーマから確実に初期化される設計となったため不要となり、削除されました。
	// これにより、ビルドプロセスが簡素化され、コードの意図がより明確になっています。

	// 汎用ビルダーのBuildを呼び出して、エラーチェックなどを行い、最終的なウィジェットを返します。
	return b.Builder.Build()
}