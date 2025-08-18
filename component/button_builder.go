package component

import (
	"fmt"
	"furoshiki/event"
	"furoshiki/style"
	"image/color"
)

// --- ButtonBuilder ---
// ButtonBuilder は、Buttonを安全かつ流れるように構築するためのビルダーです。
type ButtonBuilder struct {
	button *Button
	errors []error
}

// NewButtonBuilder は、デフォルトのスタイルで初期化されたButtonBuilderを返します。
func NewButtonBuilder() *ButtonBuilder {
	defaultStyle := style.Style{
		Background:  color.RGBA{R: 220, G: 220, B: 220, A: 255},
		TextColor:   color.Black,
		BorderColor: color.Gray{Y: 150},
		BorderWidth: 1,
		Padding:     style.Insets{Top: 5, Right: 10, Bottom: 5, Left: 10},
	}
	button := &Button{
		TextWidget: NewTextWidget(""),
	}
	button.width = 100
	button.height = 40
	button.SetStyle(defaultStyle)
	// 初期状態のスタイルをキャッシュ
	button.currentStyle = &button.style

	return &ButtonBuilder{
		button: button,
	}
}

// calculateMinSizeInternal は、ボタンのテキストとパディングに基づいて最小サイズを計算し、設定します。
func (b *ButtonBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.button.calculateMinSize()
	b.button.SetMinSize(minWidth, minHeight)
}

// CalculateMinSize は、ボタンの最小サイズを計算します。
// この呼び出しはBuild時に自動的に行われるため、通常はユーザーが呼び出す必要はありません。
func (b *ButtonBuilder) CalculateMinSize() *ButtonBuilder {
	// この呼び出しはBuild時に自動的に行われるため、通常は不要です。
	// 明示的に計算したい場合のために残しますが、内部処理はBuildに集約します。
	return b
}

// Text はボタンに表示されるテキストを設定します。
func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.button.SetText(text)
	return b
}

// Size はボタンのサイズを設定します。
func (b *ButtonBuilder) Size(width, height int) *ButtonBuilder {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("button size must be non-negative, got %dx%d", width, height))
		return b
	}
	b.button.SetSize(width, height)
	return b
}

// OnClick は、ボタンがクリックされたときに実行されるイベントハンドラを設定します。
func (b *ButtonBuilder) OnClick(onClick func()) *ButtonBuilder {
	if onClick != nil {
		b.button.AddEventHandler(event.EventClick, func(e event.Event) {
			// ハンドラの実行中にパニックが発生してもアプリケーション全体がクラッシュしないようにする
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in button click handler: %v\n", r)
				}
			}()
			onClick()
		})
	}
	return b
}

// Style はボタンの基本スタイルを設定します。
func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	existingStyle := b.button.GetStyle()
	mergedStyle := style.Merge(*existingStyle, s)
	b.button.SetStyle(mergedStyle)
	// スタイルが変更されたので、キャッシュも更新
	b.button.currentStyle = &b.button.style
	return b
}

// HoverStyle は、マウスカーソルがボタン上にあるときのスタイルを設定します。
func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	mergedHoverStyle := style.Merge(*b.button.GetStyle(), s)
	b.button.hoverStyle = &mergedHoverStyle
	return b
}

// Flex は、親がFlexLayoutの場合にボタンがどのように伸縮するかを設定します。
func (b *ButtonBuilder) Flex(flex int) *ButtonBuilder {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("button flex must be non-negative, got %d", flex))
		return b
	}
	b.button.SetFlex(flex)
	return b
}

// Build は、設定に基づいて最終的なButtonを構築して返します。
func (b *ButtonBuilder) Build() (*Button, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("button build errors: %v", b.errors)
	}
	b.calculateMinSizeInternal()
	return b.button, nil
}
