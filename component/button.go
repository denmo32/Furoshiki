package component

import (
	"fmt"
	"furoshiki/event"
	"furoshiki/style"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
type Button struct {
	*TextWidget
	hoverStyle   *style.Style
	currentStyle *style.Style // 描画時に使用するスタイルをキャッシュ
}

// SetHovered はホバー状態を設定し、描画スタイルを更新します。
func (b *Button) SetHovered(hovered bool) {
	// LayoutableWidgetの基本実装を呼び出す
	b.LayoutableWidget.SetHovered(hovered)
	// 状態に応じてスタイルを切り替える
	if b.isHovered && b.hoverStyle != nil {
		b.currentStyle = b.hoverStyle
	} else {
		b.currentStyle = &b.style
	}
}

// Draw はButtonを描画します。キャッシュされたスタイルを使用します。
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.isVisible {
		return
	}
	// キャッシュされたスタイルで描画
	drawStyledBackground(screen, b.x, b.y, b.width, b.height, *b.currentStyle)
	drawAlignedText(screen, b.text, image.Rect(b.x, b.y, b.x+b.width, b.y+b.height), *b.currentStyle)
}

// --- ButtonBuilder ---
type ButtonBuilder struct {
	button *Button
	errors []error
}

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

func (b *ButtonBuilder) calculateMinSizeInternal() {
	minWidth, minHeight := b.button.calculateMinSize()
	b.button.SetMinSize(minWidth, minHeight)
}

func (b *ButtonBuilder) CalculateMinSize() *ButtonBuilder {
	// この呼び出しはBuild時に自動的に行われるため、通常は不要です。
	// 明示的に計算したい場合のために残しますが、内部処理はBuildに集約します。
	return b
}

func (b *ButtonBuilder) Text(text string) *ButtonBuilder {
	b.button.SetText(text)
	return b
}

func (b *ButtonBuilder) Size(width, height int) *ButtonBuilder {
	if width < 0 || height < 0 {
		b.errors = append(b.errors, fmt.Errorf("button size must be non-negative, got %dx%d", width, height))
		return b
	}
	b.button.SetSize(width, height)
	return b
}

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

func (b *ButtonBuilder) Style(s style.Style) *ButtonBuilder {
	existingStyle := b.button.GetStyle()
	mergedStyle := style.Merge(*existingStyle, s)
	b.button.SetStyle(mergedStyle)
	// スタイルが変更されたので、キャッシュも更新
	b.button.currentStyle = &b.button.style
	return b
}

func (b *ButtonBuilder) HoverStyle(s style.Style) *ButtonBuilder {
	mergedHoverStyle := style.Merge(*b.button.GetStyle(), s)
	b.button.hoverStyle = &mergedHoverStyle
	return b
}

func (b *ButtonBuilder) Flex(flex int) *ButtonBuilder {
	if flex < 0 {
		b.errors = append(b.errors, fmt.Errorf("button flex must be non-negative, got %d", flex))
		return b
	}
	b.button.SetFlex(flex)
	return b
}

func (b *ButtonBuilder) Build() (*Button, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("button build errors: %v", b.errors)
	}
	b.calculateMinSizeInternal()
	return b.button, nil
}
