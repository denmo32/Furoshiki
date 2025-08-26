package main

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/ui"
	"furoshiki/widget"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game はEbitenのゲーム構造体を保持します。
type Game struct {
	root        component.Container
	contentArea *container.Container
	currentDemo component.Widget
}

// NewGame は新しいGameインスタンスを作成し、UIを構築します。
func NewGame() *Game {
	g := &Game{}

	// --- テーマとフォントの初期設定 ---
	appTheme := theme.GetCurrent()
	appTheme.SetDefaultFont(basicfont.Face7x13)
	theme.SetCurrent(appTheme)

	// --- UIの全体構造を構築 ---
	// root変数は不要なため、_で破棄します
	_, err := ui.VStack(func(b *ui.FlexBuilder) {
		b.Size(screenWidth, screenHeight).
			BackgroundColor(appTheme.BackgroundColor).
			Padding(10).
			Gap(10).
			AssignTo(&g.root) // g.rootにコンテナインスタンスを代入

		// --- 1. ナビゲーションバー ---
		b.HStack(func(b *ui.FlexBuilder) {
			b.Size(0, 30).Gap(10) // 幅は自動、高さ30

			// 各デモへの切り替えボタン
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Flex Layout").Flex(1).OnClick(func(e *event.Event) {
					g.switchToDemo(g.createFlexLayoutDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Grid Layout").Flex(1).OnClick(func(e *event.Event) {
					g.switchToDemo(g.createGridLayoutDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("ZStack Layout").Flex(1).OnClick(func(e *event.Event) {
					g.switchToDemo(g.createZStackDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("ScrollView").Flex(1).OnClick(func(e *event.Event) {
					g.switchToDemo(g.createScrollViewDemo)
				})
			})
		})

		// --- 2. デモ表示エリア ---
		b.VStack(func(b *ui.FlexBuilder) {
			b.Flex(1). // 残りのスペースをすべて使用
				AssignTo(&g.contentArea)
		})

	}).Build()

	if err != nil {
		log.Fatalf("UI build failed: %v", err)
	}

	// 初期表示のデモを設定
	g.switchToDemo(g.createFlexLayoutDemo)

	return g
}

// switchToDemo は表示するデモを切り替えます。
func (g *Game) switchToDemo(demoCreator func() (component.Widget, error)) {
	if g.contentArea == nil {
		return
	}

	// 古いデモが存在すれば削除
	if g.currentDemo != nil {
		g.contentArea.RemoveChild(g.currentDemo) // RemoveChildはCleanupも呼び出す
	}

	// 新しいデモを生成して追加
	newDemo, err := demoCreator()
	if err != nil {
		log.Printf("Failed to create demo: %v", err)
		// エラー発生時はエラーメッセージを表示するラベルに切り替える
		// ui.Labelではなく、widget.NewLabelBuilderを使用します
		errorLabelBuilder := widget.NewLabelBuilder()
		errorLabelBuilder.Text(fmt.Sprintf("Error: %v", err)).
			TextColor(color.RGBA{R: 255, A: 255})
		errorLabel, _ := errorLabelBuilder.Build()
		newDemo = errorLabel
	}

	g.currentDemo = newDemo
	g.contentArea.AddChild(g.currentDemo)
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
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Furoshiki UI Demo")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

// --- Demo Creation Functions ---

// createFlexLayoutDemo はFlexLayoutのデモ用ウィジェットを生成します。
func (g *Game) createFlexLayoutDemo() (component.Widget, error) {
	return ui.VStack(func(b *ui.FlexBuilder) {
		b.Flex(1).Padding(10).Gap(10).Border(1, color.Gray{Y: 100})

		// --- HStackのデモ ---
		b.Label(func(l *widget.LabelBuilder) { l.Text("HStack (Justify: AlignCenter, AlignItems: AlignCenter)") })
		b.HStack(func(b *ui.FlexBuilder) {
			b.Size(0, 80).Padding(5).Gap(5).Border(1, color.Gray{Y: 150}).
				Justify(layout.AlignCenter).AlignItems(layout.AlignCenter)

			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Button A").Size(80, 0) })
			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Button B").Size(100, 40) })
			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Button C").Size(60, 60) })
		})

		// --- Spacerのデモ ---
		b.Label(func(l *widget.LabelBuilder) { l.Text("HStack with Spacer") })
		b.HStack(func(b *ui.FlexBuilder) {
			b.Size(0, 50).Padding(5).Gap(5).Border(1, color.Gray{Y: 150}).AlignItems(layout.AlignCenter)

			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Left") })
			b.Spacer() // Spacerが残りの空間を埋める
			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Right") })
		})

		// --- VStackのデモ ---
		b.Label(func(l *widget.LabelBuilder) { l.Text("VStack (AlignItems: AlignStretch)") })
		b.VStack(func(b *ui.FlexBuilder) {
			b.Flex(1).Padding(5).Gap(5).Border(1, color.Gray{Y: 150}).
				AlignItems(layout.AlignStretch) // 子要素の幅がコンテナに合わせられる

			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Stretched Button 1").Size(0, 30) })
			b.Label(func(l *widget.LabelBuilder) { l.Text("This label is also stretched.") })
			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Stretched Button 2").Size(0, 30) })
		})

	}).Build()
}

// createGridLayoutDemo はGridLayoutのデモ用ウィジェットを生成します。
func (g *Game) createGridLayoutDemo() (component.Widget, error) {
	return ui.Grid(func(b *ui.GridBuilder) {
		b.Flex(1).Padding(10).Border(1, color.Gray{Y: 100}).
			Columns(5).
			HorizontalGap(8).
			VerticalGap(8)

		for i := 1; i <= 20; i++ {
			itemNum := i
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text(strconv.Itoa(itemNum)).OnClick(func(e *event.Event) {
					log.Printf("Grid item %d clicked", itemNum)
				})
			})
		}
	}).Build()
}

// createZStackDemo はZStack (AbsoluteLayout) のデモ用ウィジェットを生成します。
func (g *Game) createZStackDemo() (component.Widget, error) {
	return ui.ZStack(func(b *ui.ZStackBuilder) {
		b.Flex(1).Padding(10).Border(1, color.Gray{Y: 100})

		// 背景
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Background").
				Size(400, 300).
				BackgroundColor(color.Gray{Y: 220}).
				AbsolutePosition(50, 50)
		})

		// 中景
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Middle Layer").
				Size(300, 200).
				BackgroundColor(color.RGBA{R: 100, G: 150, B: 200, A: 255}).
				TextColor(color.White).
				AbsolutePosition(100, 100)
		})

		// 前景
		b.Button(func(btn *widget.ButtonBuilder) {
			btn.Text("Foreground Button").
				Size(150, 50).
				AbsolutePosition(175, 175).
				OnClick(func(e *event.Event) { log.Println("Foreground button clicked!") })
		})

		// ZStackの範囲外（クリッピングされる場合）
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Partially Visible").
				BackgroundColor(color.RGBA{R: 200, G: 100, B: 100, A: 255}).
				AbsolutePosition(480, 200)
		})
	}).Build()
}

// createScrollViewDemo はScrollViewのデモ用ウィジェットを生成します。
func (g *Game) createScrollViewDemo() (component.Widget, error) {
	var detailTitleLabel *widget.Label
	var detailInfoLabel *widget.Label

	return ui.HStack(func(b *ui.FlexBuilder) {
		b.Flex(1).Gap(10)

		// 左ペイン：スクロールビュー
		b.ScrollView(func(sv *widget.ScrollViewBuilder) {
			sv.Size(250, 0).Flex(1).Border(1, color.Gray{Y: 200})

			content, _ := ui.VStack(func(b *ui.FlexBuilder) {
				b.Padding(8).Gap(5)

				for i := 1; i <= 50; i++ {
					itemNumber := i
					b.Button(func(btn *widget.ButtonBuilder) {
						btn.Text(fmt.Sprintf("Item %d", itemNumber)).
							Size(0, 30). // 幅は親に合わせる
							OnClick(func(e *event.Event) {
								log.Printf("Clicked: Item %d", itemNumber)
								if detailTitleLabel != nil {
									detailTitleLabel.SetText(fmt.Sprintf("Details for Item %d", itemNumber))
								}
								if detailInfoLabel != nil {
									detailInfoLabel.SetText(fmt.Sprintf("Here you would see more detailed information about item number %d. This text is updated dynamically.", itemNumber))
								}
							})
					})
				}
			}).Build()

			sv.Content(content)
		})

		// 右ペイン：詳細表示エリア
		b.VStack(func(b *ui.FlexBuilder) {
			b.Flex(1).Padding(10).Gap(10).Border(1, color.Gray{Y: 200})

			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Details").
					Size(0, 30).
					TextColor(color.White).
					BackgroundColor(theme.GetCurrent().PrimaryColor).
					AssignTo(&detailTitleLabel)
			})

			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Please select an item from the list on the left.").
					Flex(1).
					TextAlign(style.TextAlignLeft).
					VerticalAlign(style.VerticalAlignTop).
					AssignTo(&detailInfoLabel)
			})
		})
	}).Build()
}
