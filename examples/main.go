package main

import (
	"fmt"
	"image/color"
	"log"

	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/widget"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	root          component.Widget
	width, height int
}

func NewGame() *Game {
	const initialWidth, initialHeight = 800, 600

	// --- ヘッダーの作成 ---
	header, err := widget.NewLabelBuilder().
		Text("Ebitengine UI Kit - 機能デモ").
		Size(0, 40).
		Style(style.Style{
			Font:       mplusFont,
			Background: color.RGBA{R: 80, G: 80, B: 90, A: 255},
			TextColor:  color.White,
			Padding:    style.Insets{Left: 15},
		}).
		Build()
	if err != nil {
		log.Printf("Error creating header: %v", err)
		header, _ = widget.NewLabelBuilder().Text("Header").Build()
	}

	// --- サイドバーの作成 ---
	// ボタンビルダーは自動で最小サイズを計算するため、CalculateMinSize()の呼び出しは不要
	sideButton1, err := widget.NewButtonBuilder().Text("Dashboard").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Dashboard Clicked!") }).
		Build()
	if err != nil {
		log.Printf("Error creating sideButton1: %v", err)
		sideButton1, _ = widget.NewButtonBuilder().Text("Button1").Build()
	}

	sideButton2, err := widget.NewButtonBuilder().Text("Analytics").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Analytics Clicked!") }).
		Build()
	if err != nil {
		log.Printf("Error creating sideButton2: %v", err)
		sideButton2, _ = widget.NewButtonBuilder().Text("Button2").Build()
	}

	sideButton3, err := widget.NewButtonBuilder().Text("Settings").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { fmt.Println("Settings Clicked!") }).
		Build()
	if err != nil {
		log.Printf("Error creating sideButton3: %v", err)
		sideButton3, _ = widget.NewButtonBuilder().Text("Button3").Build()
	}

	sideBar, err := container.NewContainerBuilder().
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
	if err != nil {
		log.Printf("Error creating sideBar: %v", err)
		sideBar, _ = container.NewContainerBuilder().Build()
	}
	// サイドバー自体にも最小幅を設定
	sideBar.SetMinSize(140, 0)

	// --- メインコンテンツエリアの作成 ---
	mainContentContainer, err := createMainContent()
	if err != nil {
		log.Printf("Error creating mainContentContainer: %v", err)
		mainContentContainer, _ = container.NewContainerBuilder().Build()
	}

	// --- 中間レイアウト (サイドバー + メインコンテンツ) ---
	mainLayout, err := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{
			Direction:  layout.DirectionRow,
			AlignItems: layout.AlignStretch,
			Gap:        10,
		}).
		Flex(1).
		AddChildren(sideBar, mainContentContainer).
		Build()
	if err != nil {
		log.Printf("Error creating mainLayout: %v", err)
		mainLayout, _ = container.NewContainerBuilder().Build()
	}

	// --- フッター (ステータスバー) の作成 ---
	statusLabel, err := widget.NewLabelBuilder().Text("Status: OK").Style(style.Style{Font: mplusFontSmall}).Build()
	if err != nil {
		log.Printf("Error creating statusLabel: %v", err)
		statusLabel, _ = widget.NewLabelBuilder().Text("Status").Build()
	}

	versionLabel, err := widget.NewLabelBuilder().Text("Version 1.0.0").Style(style.Style{Font: mplusFontSmall}).Build()
	if err != nil {
		log.Printf("Error creating versionLabel: %v", err)
		versionLabel, _ = widget.NewLabelBuilder().Text("Version").Build()
	}

	spacer, err := container.NewContainerBuilder().Flex(1).Build()
	if err != nil {
		log.Printf("Error creating spacer: %v", err)
		spacer, _ = container.NewContainerBuilder().Build()
	}
	spacer.SetMinSize(10, 0)

	footer, err := container.NewContainerBuilder().
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
	if err != nil {
		log.Printf("Error creating footer: %v", err)
		footer, _ = container.NewContainerBuilder().Build()
	}

	// --- ルートコンテナ (アプリケーション全体の親) ---
	rootContainer, err := container.NewContainerBuilder().
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
	if err != nil {
		log.Printf("Error creating rootContainer: %v", err)
		rootContainer, _ = container.NewContainerBuilder().
			Size(initialWidth, initialHeight).
			Build()
	}

	return &Game{
		root:   rootContainer,
		width:  initialWidth,
		height: initialHeight,
	}
}

func createMainContent() (*container.Container, error) {
	refreshButton, err := widget.NewButtonBuilder().Text("Refresh").OnClick(func() { fmt.Println("Refresh!") }).Style(style.Style{Font: mplusFont}).Build()
	if err != nil {
		return nil, fmt.Errorf("error creating refreshButton: %w", err)
	}

	exportButton, err := widget.NewButtonBuilder().Text("Export").OnClick(func() { fmt.Println("Export!") }).
		Style(style.Style{Font: mplusFont, Margin: style.Insets{Left: 20}}).
		Build()
	if err != nil {
		return nil, fmt.Errorf("error creating exportButton: %w", err)
	}

	controlPanel, err := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, Justify: layout.AlignEnd, Gap: 10}).
		Size(0, 40).
		AddChildren(refreshButton, exportButton).
		Build()
	if err != nil {
		return nil, fmt.Errorf("error creating controlPanel: %w", err)
	}

	card1, err := createCard("Total Revenue", "$45,231.89", "Increased by 20%")
	if err != nil {
		return nil, fmt.Errorf("error creating card1: %w", err)
	}

	card2, err := createCard("Subscriptions", "+2,350", "Gained 180 this week")
	if err != nil {
		return nil, fmt.Errorf("error creating card2: %w", err)
	}

	card3, err := createCard("Active Users", "9,876", "Currently online")
	if err != nil {
		return nil, fmt.Errorf("error creating card3: %w", err)
	}

	cardArea, err := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, AlignItems: layout.AlignStart, Gap: 15}).
		Flex(1).
		AddChildren(card1, card2, card3).
		Build()
	if err != nil {
		return nil, fmt.Errorf("error creating cardArea: %w", err)
	}
	cardArea.SetMinSize(300, 100) // カードエリアが極端に縮まらないように設定

	return container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch, Gap: 10}).
		Flex(1).
		Style(style.Style{Padding: style.Insets{Top: 15, Right: 15, Bottom: 15, Left: 15}}).
		AddChildren(controlPanel, cardArea).
		Build()
}

func createCard(title, mainText, subText string) (component.Widget, error) {
	titleLabel, err := widget.NewLabelBuilder().Text(title).Style(style.Style{Font: mplusFontSmall, TextColor: color.Gray{Y: 100}}).Build()
	if err != nil {
		return nil, fmt.Errorf("error creating titleLabel: %w", err)
	}

	mainLabel, err := widget.NewLabelBuilder().Text(mainText).Style(style.Style{Font: mplusFont, TextColor: color.Black}).Size(0, 40).Build()
	if err != nil {
		return nil, fmt.Errorf("error creating mainLabel: %w", err)
	}

	subLabel, err := widget.NewLabelBuilder().Text(subText).Style(style.Style{Font: mplusFontSmall, TextColor: color.Gray{Y: 120}}).Build()
	if err != nil {
		return nil, fmt.Errorf("error creating subLabel: %w", err)
	}

	spacer, err := container.NewContainerBuilder().Flex(1).Build()
	if err != nil {
		return nil, fmt.Errorf("error creating spacer: %w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("error creating card: %w", err)
	}

	card.SetMinSize(150, 80) // カード自体の最小サイズも設定
	return card, nil
}

func (g *Game) Update() error {
	// --- イベント処理 ---
	// マウスカーソルの位置を取得
	cx, cy := ebiten.CursorPosition()
	// カーソル下のコンポーネントを特定
	target := g.root.HitTest(cx, cy)

	// マウスイベントをUIツリーにディスパッチ
	event.GetDispatcher().Dispatch(target, cx, cy)

	// --- UIの更新 ---
	// UIツリー全体の更新処理を呼び出す。
	// これにより、ダーティフラグが立っているコンテナのレイアウトが再計算され、
	// 必要に応じてウィジェットの状態が更新される。
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

func (g *Game) Cleanup() {
	// ゲーム終了時にリソースを解放
	if cleanup, ok := g.root.(interface{ Cleanup() }); ok {
		cleanup.Cleanup()
	}
	// イベントディスパッチャをリセット
	event.GetDispatcher().Reset()
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("UI Library Example (Improved)")
	game := NewGame()

	// ゲーム終了時にクリーンアップ処理を実行
	defer game.Cleanup()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
