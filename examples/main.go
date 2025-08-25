package main

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/ui"
	"furoshiki/widget"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
)

// Game はEbitenのゲーム構造体を保持します。
type Game struct {
	root             component.Widget
	detailTitleLabel *widget.Label
	detailInfoLabel  *widget.Label
}

// NewGame は新しいGameインスタンスを作成し、UIを構築します。
func NewGame() *Game {
	g := &Game{}

	appTheme := theme.GetCurrent()
	appTheme.SetDefaultFont(basicfont.Face7x13)
	theme.SetCurrent(appTheme)

	root, err := ui.HStack(func(b *ui.FlexBuilder) {
		b.Size(800, 600).
			BackgroundColor(appTheme.BackgroundColor).
			Padding(10).
			Gap(10)

		// 左ペイン：スクロールビュー
		b.ScrollView(func(sv *widget.ScrollViewBuilder) {
			sv.Size(250, 580).
				Border(1, color.Gray{Y: 200})

			content, _ := ui.VStack(func(b *ui.FlexBuilder) {
				b.Padding(8).Gap(5)

				for i := 1; i <= 50; i++ {
					itemNumber := i
					b.Button(func(btn *widget.ButtonBuilder) {
						btn.Text(fmt.Sprintf("Item %d", itemNumber)).
							Size(220, 30).
							OnClick(func(e *event.Event) {
								log.Printf("Clicked: Item %d", itemNumber)
								if g.detailTitleLabel != nil {
									g.detailTitleLabel.SetText(fmt.Sprintf("Details for Item %d", itemNumber))
								}
								if g.detailInfoLabel != nil {
									g.detailInfoLabel.SetText(fmt.Sprintf("Here you would see more detailed information about item number %d. This text is updated dynamically.", itemNumber))
								}
							})
					})
				}
			}).Build()

			sv.Content(content)
		})

		// 右ペイン：詳細表示エリア
		b.VStack(func(b *ui.FlexBuilder) {
			b.Flex(1).
				Padding(10).
				Gap(10).
				Border(1, color.Gray{Y: 200}).
				AlignItems(layout.AlignStretch)

			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Details").
					Size(0, 30).
					TextColor(color.White).
					BackgroundColor(appTheme.PrimaryColor).
					AssignTo(&g.detailTitleLabel)
			})

			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Please select an item from the list on the left.").
					Flex(1).
					TextAlign(style.TextAlignLeft).
					VerticalAlign(style.VerticalAlignTop).
					AssignTo(&g.detailInfoLabel)
			})

			b.Spacer()

			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Furoshiki UI Demo").
					Size(0, 20).
					TextAlign(style.TextAlignRight).
					TextColor(color.Gray{Y: 128})
			})
		})

	}).Build()

	if err != nil {
		log.Fatalf("UI build failed: %v", err)
	}

	g.root = root
	return g
}

// Update はゲームの状態を更新します。
func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	dispatcher := event.GetDispatcher()
	target := g.root.HitTest(cx, cy)
	dispatcher.Dispatch(target, cx, cy)
	g.root.Update()
	return nil
}

// Draw はゲームを描画します。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 255})
	g.root.Draw(screen)
}

// Layout はEbitenにゲームの画面サイズを伝えます。
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Furoshiki UI - ScrollView Example")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
