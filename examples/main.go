package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

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

// Game はアプリケーションの状態を保持します
type Game struct {
	root               component.Widget
	width, height      int
	mainLayout         *container.Container // サイドバーとメインコンテンツを保持
	mainContent        component.Widget     // 現在のメインコンテンツ
	statusLabel        *widget.Label
	dynamicWidgetCount int // 動的に追加されるウィジェットのカウンター
}

// NewGame は新しいゲームインスタンスを作成します
func NewGame() *Game {
	const initialWidth, initialHeight = 800, 600
	g := &Game{
		width:  initialWidth,
		height: initialHeight,
	}

	// --- ヘッダー ---
	header, _ := widget.NewLabelBuilder().
		Text("Furoshiki UI - 機能デモ").
		Size(0, 40).
		Style(style.Style{
			Font:       mplusFont,
			Background: color.RGBA{R: 80, G: 80, B: 90, A: 255},
			TextColor:  color.White,
			Padding:    style.Insets{Left: 15},
		}).
		Build()

	// --- フッター (ステータスバー) ---
	g.statusLabel, _ = widget.NewLabelBuilder().Text("Status: OK").Style(style.Style{Font: mplusFontSmall}).Build()
	versionLabel, _ := widget.NewLabelBuilder().Text("Version 1.0.0").Style(style.Style{Font: mplusFontSmall}).Build()
	spacer, _ := container.NewContainerBuilder().Flex(1).Build()
	spacer.SetMinSize(10, 0)

	footer, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, AlignItems: layout.AlignCenter}).
		Size(0, 25).
		Style(style.Style{
			Background:  color.RGBA{R: 220, G: 220, B: 220, A: 255},
			Padding:     style.Insets{Left: 10, Right: 10},
			BorderColor: color.Gray{Y: 180},
			BorderWidth: 1,
		}).
		AddChildren(g.statusLabel, spacer, versionLabel).
		Build()

	// --- サイドバー ---
	sideButton1, _ := widget.NewButtonBuilder().Text("Dashboard").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { g.showDashboardView() }).
		Build()

	sideButton2, _ := widget.NewButtonBuilder().Text("Analytics").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { g.showAnalyticsView() }).
		Build()

	sideButton3, _ := widget.NewButtonBuilder().Text("Settings").Flex(0).Size(0, 35).
		Style(style.Style{Font: mplusFont}).
		HoverStyle(style.Style{Background: color.RGBA{R: 200, G: 220, B: 255, A: 255}}).
		OnClick(func() { g.showSettingsView() }).
		Build()

	sideBar, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, Justify: layout.AlignStart, AlignItems: layout.AlignStretch, Gap: 5}).
		Size(150, 0).
		Style(style.Style{
			Background: color.RGBA{R: 240, G: 240, B: 240, A: 255},
			Padding:    style.Insets{Top: 10, Right: 10, Bottom: 10, Left: 10},
		}).
		AddChildren(sideButton1, sideButton2, sideButton3).
		Build()
	sideBar.SetMinSize(140, 0)

	// --- メインコンテンツエリアの初期化 ---
	g.mainContent = g.createDashboardView() // 初期ビュー

	// --- 中間レイアウト (サイドバー + メインコンテンツ) ---
	g.mainLayout, _ = container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, AlignItems: layout.AlignStretch, Gap: 0}).
		Flex(1).
		AddChildren(sideBar, g.mainContent).
		Build()

	// --- ルートコンテナ ---
	rootContainer, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch}).
		Size(initialWidth, initialHeight).
		Style(style.Style{Background: color.RGBA{R: 250, G: 250, B: 250, A: 255}}).
		AddChildren(header, g.mainLayout, footer).
		Build()

	g.root = rootContainer
	g.updateStatus("Welcome! Dashboard view displayed.")
	return g
}

// updateStatus はフッターのステータスラベルを更新します
func (g *Game) switchMainContent(newContent component.Widget) {
	if g.mainContent != nil {
		g.mainLayout.RemoveChild(g.mainContent)
	}
	g.mainContent = newContent
	g.mainLayout.AddChild(g.mainContent)
}

func (g *Game) updateStatus(message string) {
	g.statusLabel.SetText("Status: " + message)
}

// --- ビュー作成関数 ---

func (g *Game) showDashboardView() {
	g.switchMainContent(g.createDashboardView())
	g.updateStatus("Dashboard view displayed.")
}

func (g *Game) showAnalyticsView() {
	g.switchMainContent(g.createAnalyticsView())
	g.updateStatus("Analytics view displayed.")
}

func (g *Game) showSettingsView() {
	g.switchMainContent(g.createSettingsView())
	g.updateStatus("Settings view displayed.")
}

func (g *Game) createDashboardView() component.Widget {
	g.dynamicWidgetCount = 0
	// 動的ウィジェットを保持するコンテナ
	dynamicArea, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, AlignItems: layout.AlignStart, Gap: 10, Wrap: true}).
		Flex(1).
		Style(style.Style{
			Padding:     style.Insets{Top: 10, Right: 10, Bottom: 10, Left: 10},
			BorderColor: color.Gray{Y: 200},
			BorderWidth: 1,
		}).
		Build()
	dynamicArea.SetMinSize(100, 100)

	addWidgetButton, _ := widget.NewButtonBuilder().Text("Add Widget").Style(style.Style{Font: mplusFont}).
		OnClick(func() {
			g.dynamicWidgetCount++
			newLabel, _ := widget.NewLabelBuilder().
				Text("Widget " + strconv.Itoa(g.dynamicWidgetCount)).
				Size(100, 40).
				Style(style.Style{
					Font:       mplusFontSmall,
					Background: color.RGBA{B: 100, A: 50},
					Padding:    style.Insets{Top: 5, Right: 5, Bottom: 5, Left: 5},
				}).
				Build()
			dynamicArea.AddChild(newLabel)
			g.updateStatus(fmt.Sprintf("Added Widget %d.", g.dynamicWidgetCount))
		}).Build()

	removeWidgetButton, _ := widget.NewButtonBuilder().Text("Remove Last").Style(style.Style{Font: mplusFont}).
		OnClick(func() {
			children := dynamicArea.GetChildren()
			if len(children) > 0 {
				lastChild := children[len(children)-1]
				dynamicArea.RemoveChild(lastChild)
				g.updateStatus("Removed last widget.")
			} else {
				g.updateStatus("No widgets to remove.")
			}
		}).Build()

	controlPanel, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, Justify: layout.AlignStart, Gap: 10}).
		Size(0, 60).
		Style(style.Style{Padding: style.Insets{Top: 10, Bottom: 10}}).
		AddChildren(addWidgetButton, removeWidgetButton).
		Build()

	view, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch}).
		Flex(1).
		Style(style.Style{Padding: style.Insets{Top: 15, Right: 15, Bottom: 15, Left: 15}}).
		AddChildren(controlPanel, dynamicArea).
		Build()
	return view
}

func (g *Game) createAnalyticsView() component.Widget {
	title, _ := widget.NewLabelBuilder().Text("Analytics Overview").Style(style.Style{Font: mplusFont}).Build()
	placeholder, _ := widget.NewLabelBuilder().
		Text("This is a placeholder for analytics data and charts.").
		Size(0, 0).
		Flex(1).
		Style(style.Style{
			Font:        mplusFontSmall,
			Background:  color.RGBA{A: 30},
			BorderColor: color.Gray{Y: 200},
			BorderWidth: 1,
			Padding:     style.Insets{Top: 20, Left: 20},
		}).
		Build()

	view, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch, Gap: 15}).
		Flex(1).
		Style(style.Style{Padding: style.Insets{Top: 15, Right: 15, Bottom: 15, Left: 15}}).
		AddChildren(title, placeholder).
		Build()
	return view
}

func (g *Game) createSettingsView() component.Widget {
	// デモ用ウィジェット
	box1, _ := widget.NewLabelBuilder().Text("Box 1").Size(80, 50).Style(style.Style{Background: color.RGBA{R: 255, G: 200, B: 200, A: 255}, Font: mplusFontSmall}).Build()
	box2, _ := widget.NewLabelBuilder().Text("Box 2").Size(80, 50).Style(style.Style{Background: color.RGBA{R: 200, G: 255, B: 200, A: 255}, Font: mplusFontSmall}).Build()
	box3, _ := widget.NewLabelBuilder().Text("Box 3").Size(80, 50).Style(style.Style{Background: color.RGBA{R: 200, G: 200, B: 255, A: 255}, Font: mplusFontSmall}).Build()

	// レイアウト変更デモ用のコンテナ
	layoutDemoContainer := &layout.FlexLayout{Direction: layout.DirectionRow, Justify: layout.AlignCenter, AlignItems: layout.AlignCenter, Gap: 10}
	demoArea, _ := container.NewContainerBuilder().
		Layout(layoutDemoContainer).
		Flex(1).
		Style(style.Style{
			BorderColor: color.Gray{Y: 200},
			BorderWidth: 1,
			Padding:     style.Insets{Top: 10, Right: 10, Bottom: 10, Left: 10},
		}).
		AddChildren(box1, box2, box3).
		Build()

	// レイアウト方向を切り替えるボタン
	toggleDirButton, _ := widget.NewButtonBuilder().Text("Toggle Layout Direction").Style(style.Style{Font: mplusFont}).
		OnClick(func() {
			if layoutDemoContainer.Direction == layout.DirectionRow {
				layoutDemoContainer.Direction = layout.DirectionColumn
				g.updateStatus("Layout changed to: Column")
			} else {
				layoutDemoContainer.Direction = layout.DirectionRow
				g.updateStatus("Layout changed to: Row")
			}
			// レイアウト変更を適用するためにダーティフラグを立てる
			demoArea.MarkDirty(true)
		}).Build()

	// 表示・非表示を切り替えるウィジェット
	toggleTarget, _ := widget.NewLabelBuilder().Text("Toggle Me!").Size(120, 40).Style(style.Style{Background: color.RGBA{R: 255, G: 255, B: 150, A: 255}, Font: mplusFontSmall}).Build()

	toggleVisibilityButton, _ := widget.NewButtonBuilder().Text("Toggle Visibility").Style(style.Style{Font: mplusFont}).
		OnClick(func() {
			isVisible := toggleTarget.IsVisible()
			toggleTarget.SetVisible(!isVisible)
			if isVisible {
				g.updateStatus("Widget hidden.")
			} else {
				g.updateStatus("Widget shown.")
			}
		}).Build()

	controlPanel, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionRow, Justify: layout.AlignStart, AlignItems: layout.AlignCenter, Gap: 10}).
		Size(0, 60).
		Style(style.Style{Padding: style.Insets{Top: 10, Bottom: 10}}).
		AddChildren(toggleDirButton, toggleVisibilityButton, toggleTarget).
		Build()

	view, _ := container.NewContainerBuilder().
		Layout(&layout.FlexLayout{Direction: layout.DirectionColumn, AlignItems: layout.AlignStretch, Gap: 10}).
		Flex(1).
		Style(style.Style{Padding: style.Insets{Top: 15, Right: 15, Bottom: 15, Left: 15}}).
		AddChildren(controlPanel, demoArea).
		Build()
	return view
}

// --- Ebitengine Game Loop ---

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	// g.root.HitTestは component.Widget を返しますが、これは event.eventTarget インターフェースを満たしているので、
	// そのまま event.GetDispatcher().Dispatch に渡すことができます。
	target := g.root.HitTest(cx, cy)
	event.GetDispatcher().Dispatch(target, cx, cy)
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
	if cleanup, ok := g.root.(interface{ Cleanup() }); ok {
		cleanup.Cleanup()
	}
	event.GetDispatcher().Reset()
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Furoshiki UI - Interactive Demo")
	game := NewGame()
	defer game.Cleanup()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}