package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
)

var (
	mplusFont      font.Face
	mplusFontSmall font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusFontSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	root             component.Widget
	hoveredComponent component.Widget
	width, height    int
}

func NewGame() *Game {
	const initialWidth, initialHeight = 800, 600

	// --- ヘッダーの作成 ---
	header, _ := component.NewLabelBuilder().
		Text("Ebitengine UI Kit - 機能デモ").
		Size(0, 40).
		Style(style.Style{
			Font:       mplusFont,
			Background: color.RGBA{R: 80, G: 80, B: 90, A: 255},
			TextColor:  color.White,
			Padding:    style.Insets{Left: 15},
		}).
		Build()

	// --- サイドバーの作成 ---
	// ボタンビルダーは自動で最小サイズを計算するため、CalculateMinSize()の呼び出しは不要
	sideButton1, _ := component.NewButtonBuilder().Text("Dashboard").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Dashboard Clicked!") }).
		Build()

	sideButton2, _ := component.NewButtonBuilder().Text("Analytics").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Analytics Clicked!") }).
		Build()

	sideButton3, _ := component.NewButtonBuilder().Text("Settings").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Settings Clicked!") }).
		Build()

	sideBar, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionColumn,
			Justify:    layout.AlignStart,
			AlignItems: layout.AlignStretch,
			Gap:        5,
		}).
		Size(150, 0).
		Style(style.Style{
			Background: color.RGBA{R: 240, G: 240, B: 240, A: 255},
			Padding:    style.Insets{Top: 10, Right: 10, Bottom: 10, Left: 10},
		}).
		AddChildren(sideButton1, sideButton2, sideButton3).
		Build()
	// サイドバー自体にも最小幅を設定
	sideBar.SetMinSize(140, 0)

	// --- メインコンテンツエリアの作成 ---
	mainContentContainer, _ := createMainContent()

	// --- 中間レイアウト (サイドバー + メインコンテンツ) ---
	mainLayout, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionRow,
			AlignItems: layout.AlignStretch,
			Gap:        10,
		}).
		Flex(1).
		AddChildren(sideBar, mainContentContainer).
		Build()

	// --- フッター (ステータスバー) の作成 ---
	statusLabel, _ := component.NewLabelBuilder().Text("Status: OK").Style(style.Style{Font: mplusFontSmall}).Build()
	versionLabel, _ := component.NewLabelBuilder().Text("Version 1.0.0").Style(style.Style{Font: mplusFontSmall}).Build()
	spacer, _ := container.NewContainerBuilder().Flex(1).Build()
	spacer.SetMinSize(10, 0)

	footer, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionRow,
			AlignItems: layout.AlignCenter,
		}).
		Size(0, 25).
		Style(style.Style{
			Background:  color.RGBA{R: 220, G: 220, B: 220, A: 255},
			Padding:     style.Insets{Left: 10, Right: 10},
			BorderColor: color.Gray{Y: 180},
			BorderWidth: 1,
		}).
		AddChildren(statusLabel, spacer, versionLabel).
		Build()

	// --- ルートコンテナ (アプリケーション全体の親) ---
	rootContainer, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionColumn,
			AlignItems: layout.AlignStretch,
		}).
		Size(initialWidth, initialHeight).
		Style(style.Style{
			Background: color.RGBA{R: 250, G: 250, B: 250, A: 255},
		}).
		AddChildren(header, mainLayout, footer).
		Build()

	return &Game{
		root:   rootContainer,
		width:  initialWidth,
		height: initialHeight,
	}
}

func createMainContent() (*container.Container, error) {
	refreshButton, _ := component.NewButtonBuilder().Text("Refresh").OnClick(func() { fmt.Println("Refresh!") }).Style(style.Style{Font: mplusFont}).Build()
	exportButton, _ := component.NewButtonBuilder().Text("Export").OnClick(func() { fmt.Println("Export!") }).
		Style(style.Style{Font: mplusFont, Margin: style.Insets{Left: 20}}).
		Build()

	controlPanel, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, Justify: layout.AlignEnd, Gap: 10}).
		Size(0, 40).
		AddChildren(refreshButton, exportButton).
		Build()

	card1, _ := createCard("Total Revenue", "$45,231.89", "Increased by 20%")
	card2, _ := createCard("Subscriptions", "+2,350", "Gained 180 this week")
	card3, _ := createCard("Active Users", "9,876", "Currently online")
	cardArea, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, AlignItems: layout.AlignStart, Gap: 15}).
		Flex(1).
		AddChildren(card1, card2, card3).
		Build()
	cardArea.SetMinSize(300, 100) // カードエリアが極端に縮まないように設定

	return container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch, Gap: 10}).
		Flex(1).
		Style(style.Style{Padding: style.Insets{Top: 15, Right: 15, Bottom: 15, Left: 15}}).
		AddChildren(controlPanel, cardArea).
		Build()
}

func createCard(title, mainText, subText string) (component.Widget, error) {
	titleLabel, _ := component.NewLabelBuilder().Text(title).Style(style.Style{Font: mplusFontSmall, TextColor: color.Gray{Y: 100}}).Build()
	mainLabel, _ := component.NewLabelBuilder().Text(mainText).Style(style.Style{Font: mplusFont, TextColor: color.Black}).Size(0, 40).Build()
	subLabel, _ := component.NewLabelBuilder().Text(subText).Style(style.Style{Font: mplusFontSmall, TextColor: color.Gray{Y: 120}}).Build()
	spacer, _ := container.NewContainerBuilder().Flex(1).Build()

	card, err := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionColumn,
			AlignItems: layout.AlignStart,
			Gap:        5,
		}).
		Flex(1).
		Style(style.Style{
			Background:  color.White,
			BorderColor: color.Gray{Y: 200},
			BorderWidth: 1,
			Padding:     style.Insets{Top: 10, Right: 12, Bottom: 10, Left: 12},
		}).
		AddChildren(titleLabel, mainLabel, spacer, subLabel).
		Build()
	if err == nil {
		card.SetMinSize(150, 80) // カード自体の最小サイズも設定
	}
	return card, err
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	target := g.root.HitTest(cx, cy)

	if target != g.hoveredComponent {
		if g.hoveredComponent != nil {
			g.hoveredComponent.SetHovered(false)
			g.hoveredComponent.HandleEvent(event.Event{
				Type:   event.MouseLeave,
				Target: g.hoveredComponent,
			})
		}
		if target != nil {
			target.SetHovered(true)
			target.HandleEvent(event.Event{
				Type:   event.MouseEnter,
				Target: target,
			})
		}
		g.hoveredComponent = target
	}

	if g.hoveredComponent != nil {
		g.hoveredComponent.HandleEvent(event.Event{
			Type:   event.MouseMove,
			Target: g.hoveredComponent,
			X:      cx,
			Y:      cy,
		})
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.hoveredComponent != nil {
			g.hoveredComponent.HandleEvent(event.Event{
				Type:        event.EventClick,
				Target:      g.hoveredComponent,
				X:           cx,
				Y:           cy,
				Timestamp:   time.Now().UnixNano(),
				MouseButton: ebiten.MouseButtonLeft, // 押されたマウスボタンの情報をイベントに追加
			})
		}
	}

	g.root.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 250, G: 250, B: 250, A: 255})
	g.root.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.width != outsideWidth || g.height != outsideHeight {
		g.width = outsideWidth
		g.height = outsideHeight
		g.root.SetSize(outsideWidth, outsideHeight)
	}
	return g.width, g.height
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("UI Library Example (Improved)")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
