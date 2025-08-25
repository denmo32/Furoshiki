package theme

import (
	"furoshiki/style"
	"image/color"
	"sync"

	"golang.org/x/image/font"
)

// ButtonTheme はButtonウィジェットに関連するスタイルを定義します。
type ButtonTheme struct {
	Normal, Hovered, Pressed, Disabled style.Style
}

// LabelTheme はLabelウィジェットに関連するスタイルを定義します。
type LabelTheme struct {
	Default style.Style
}

// Theme はUI全体の視覚的スタイルを定義します。
type Theme struct {
	DefaultFont     font.Face
	BackgroundColor color.Color
	TextColor       color.Color
	PrimaryColor    color.Color
	SecondaryColor  color.Color
	Button          ButtonTheme
	Label           LabelTheme
}

// SetDefaultFont はテーマ内のすべてのウィジェットスタイルにデフォルトフォントを設定するヘルパーです。
func (t *Theme) SetDefaultFont(f font.Face) {
	if f == nil {
		return
	}
	t.DefaultFont = f
	t.Button.Normal.Font = style.PFont(f)
	t.Button.Hovered.Font = style.PFont(f)
	t.Button.Pressed.Font = style.PFont(f)
	t.Button.Disabled.Font = style.PFont(f)
	t.Label.Default.Font = style.PFont(f)
}

var (
	currentTheme *Theme
	mutex        sync.RWMutex
)

// SetCurrent はアプリケーション全体で使用されるテーマを設定します。
func SetCurrent(t *Theme) {
	mutex.Lock()
	defer mutex.Unlock()
	currentTheme = t
}

// GetCurrent は現在設定されているテーマを返します。
func GetCurrent() *Theme {
	mutex.RLock()
	if currentTheme != nil {
		defer mutex.RUnlock()
		return currentTheme
	}
	mutex.RUnlock()

	mutex.Lock()
	defer mutex.Unlock()
	if currentTheme == nil {
		currentTheme = newDefaultTheme()
	}
	return currentTheme
}

// newDefaultTheme はライブラリのデフォルトテーマを生成します。
func newDefaultTheme() *Theme {
	lightGray := color.RGBA{220, 220, 220, 255}
	darkGray := color.RGBA{105, 105, 105, 255}
	white := color.White
	black := color.Black

	btnNormal := style.Style{
		Background:    style.PColor(lightGray),
		TextColor:     style.PColor(black),
		BorderColor:   style.PColor(darkGray),
		BorderWidth:   style.PFloat32(1),
		Padding:       style.PInsets(style.Insets{Top: 5, Right: 10, Bottom: 5, Left: 10}),
		TextAlign:     style.PTextAlignType(style.TextAlignCenter),
		VerticalAlign: style.PVerticalAlignType(style.VerticalAlignMiddle),
	}
	btnHovered := style.Merge(btnNormal, style.Style{Opacity: style.PFloat64(0.9)})
	btnPressed := style.Merge(btnHovered, style.Style{Background: style.PColor(darkGray), TextColor: style.PColor(white), Opacity: style.PFloat64(1.0)})
	btnDisabled := style.Merge(btnNormal, style.Style{Opacity: style.PFloat64(0.5)})

	lblDefault := style.Style{
		Background: style.PColor(color.Transparent),
		TextColor:  style.PColor(black),
		Padding:    style.PInsets(style.Insets{Top: 2, Right: 5, Bottom: 2, Left: 5}),
	}

	return &Theme{
		BackgroundColor: color.RGBA{245, 245, 245, 255},
		TextColor:       black,
		PrimaryColor:    color.RGBA{70, 130, 180, 255}, // SteelBlue
		SecondaryColor:  lightGray,
		Button: ButtonTheme{
			Normal: btnNormal, Hovered: btnHovered, Pressed: btnPressed, Disabled: btnDisabled,
		},
		Label: LabelTheme{Default: lblDefault},
	}
}
