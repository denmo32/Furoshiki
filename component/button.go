package component

import (
	"furoshiki/style"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Button component ---
// Button は、クリック可能なUI要素です。TextWidgetを拡張し、ホバー状態のスタイル管理機能を追加します。
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