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
	"image"
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
			b.Size(0, 30).Gap(5) // 幅は自動、高さ30

			// 各デモへの切り替えボタン
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Flex Layout").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createFlexLayoutDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Flex Wrap").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createFlexWrapDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Text Wrap").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createWrapTextDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Grid Layout").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createGridLayoutDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Advanced Grid").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createAdvancedGridLayoutDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("ZStack Layout").Flex(1).AddOnClick(func(e *event.Event) {
					g.switchToDemo(g.createZStackDemo)
				})
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("ScrollView").Flex(1).AddOnClick(func(e *event.Event) {
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

// Layout はEbitenにゲームの画面サイズを伝え、UIのレイアウトを計算します。
// この関数はEbitenによって毎フレーム呼び出されます。
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	w, h := screenWidth, screenHeight

	// UIの再レイアウトが必要な場合のみ、Measure/Arrangeパスを実行します。
	if g.root.NeedsRelayout() {
		// 1. Measure Pass: UI全体の要求サイズを計算します。
		// ルートコンテナの場合、利用可能なサイズは画面全体です。
		g.root.Measure(image.Point{X: w, Y: h})

		// 2. Arrange Pass: 計算された情報に基づき、全ウィジェットを配置します。
		// ルートコンテナは画面全体を占有します。
		bounds := image.Rect(0, 0, w, h)
		g.root.SetPosition(bounds.Min.X, bounds.Min.Y)
		g.root.SetSize(bounds.Dx(), bounds.Dy())
		g.root.Arrange(bounds)

		// 3. Dirty Flag Cleanup: レイアウトが完了したので、ダーティフラグをクリアします。
		clearDirty(g.root)
	}

	return w, h
}

// clearDirty は、指定されたウィジェットとそのすべての子孫のダーティフラグを再帰的にクリアします。
func clearDirty(w component.Widget) {
	w.ClearDirty()
	if c, ok := w.(component.Container); ok {
		for _, child := range c.GetChildren() {
			clearDirty(child)
		}
	}
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
			// [修正] .Flex(1) を追加して、このVStackが親コンテナ内で利用可能な垂直スペースを
			// すべて使用するようにします。これにより、AlignStretchの効果が正しく表示されます。
			b.Flex(1).Padding(5).Gap(5).Border(1, color.Gray{Y: 150}).
				AlignItems(layout.AlignStretch) // 子要素の幅がコンテナに合わせられる

			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Stretched Button 1").Size(0, 30) })
			b.Label(func(l *widget.LabelBuilder) { l.Text("This label is also stretched.") })
			b.Button(func(btn *widget.ButtonBuilder) { btn.Text("Stretched Button 2").Size(0, 30) })
		})

	}).Build()
}

// createFlexWrapDemo はFlexLayoutの折り返し機能のデモ用ウィジェットを生成します。
func (g *Game) createFlexWrapDemo() (component.Widget, error) {
	return ui.VStack(func(b *ui.FlexBuilder) {
		b.Flex(1).Padding(10).Gap(10).Border(1, color.Gray{Y: 100})

		// --- 基本的な折り返しのデモ ---
		b.Label(func(l *widget.LabelBuilder) { l.Text("HStack with Wrap(true)") })
		b.HStack(func(b *ui.FlexBuilder) {
			b.Size(0, 100).Padding(5).Gap(8).Border(1, color.Gray{Y: 150}).
				Wrap(true) // 折り返しを有効にする

			// 多数のアイテムを追加して折り返しを発生させる
			for i := 1; i <= 18; i++ {
				b.Button(func(btn *widget.ButtonBuilder) {
					btn.Text("Item "+strconv.Itoa(i)).Size(100, 30)
				})
			}
		})

		// --- AlignContent のデモ ---
		b.Label(func(l *widget.LabelBuilder) { l.Text("HStack with Wrap(true) and AlignContent(AlignCenter)") })
		b.HStack(func(b *ui.FlexBuilder) {
			b.Flex(1).Padding(5).Gap(8).Border(1, color.Gray{Y: 150}).
				Wrap(true).                      // 折り返しを有効にする
				AlignContent(layout.AlignCenter) // 折り返したライン全体をコンテナの中央に配置

			for i := 1; i <= 18; i++ {
				b.Button(func(btn *widget.ButtonBuilder) {
					btn.Text("Item "+strconv.Itoa(i)).Size(100, 30)
				})
			}
		})
	}).Build()
}

// createWrapTextDemo はテキストの折り返し機能のデモ用ウィジェットを生成します。
func (g *Game) createWrapTextDemo() (component.Widget, error) {
	longText := "This is a very long text that should wrap automatically when it reaches the edge of the widget. The layout system will then adjust the widget's height to accommodate all the wrapped lines of text."
	return ui.VStack(func(b *ui.FlexBuilder) {
		b.Flex(1).Padding(10).Gap(10).Border(1, color.Gray{Y: 100}).ClipChildren(true)

		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Label with WrapText(true)")
		})
		b.Label(func(l *widget.LabelBuilder) {
			l.Text(longText).
				WrapText(true). // テキストの折り返しを有効化
				Size(300, 0).   // 幅を固定し、高さはレイアウトに任せる
				Border(1, color.Gray{Y: 150}).
				Padding(5)
		})

		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Stretched Label with WrapText(true) and VerticalAlignTop")
		})
		b.Label(func(l *widget.LabelBuilder) {
			l.Text(longText+" "+longText).
				WrapText(true).                        // テキストの折り返しを有効化
				Flex(1).                               // 残りの垂直スペースをすべて使用
				VerticalAlign(style.VerticalAlignTop). // 上揃え
				Border(1, color.Gray{Y: 150}).
				Padding(5)
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
				btn.Text(strconv.Itoa(itemNum)).AddOnClick(func(e *event.Event) {
					log.Printf("Grid item %d clicked", itemNum)
				})
			})
		}
	}).Build()
}

// createAdvancedGridLayoutDemo はAdvancedGridLayoutのデモ用ウィジェットを生成します。
func (g *Game) createAdvancedGridLayoutDemo() (component.Widget, error) {
	return ui.AdvancedGrid(func(b *ui.AdvancedGridBuilder) {
		b.Flex(1).Padding(10).Border(1, color.Gray{Y: 100}).Gap(8)

		// 列と行のサイズを定義
		b.Columns(ui.Fixed(150), ui.Weight(1), ui.Weight(2))
		b.Rows(ui.Fixed(50), ui.Weight(1), ui.Fixed(30))

		// ヘッダー (0行目、0列目から3列にまたがる)
		b.LabelAt(0, 0, 1, 3, func(l *widget.LabelBuilder) {
			l.Text("Header (Spans 3 Columns)").
				BackgroundColor(color.Gray{Y: 180}).
				TextAlign(style.TextAlignCenter)
		})

		// サイドバー (1行目、0列目)
		b.LabelAt(1, 0, 1, 1, func(l *widget.LabelBuilder) {
			l.Text("Sidebar (150px Fixed)").
				BackgroundColor(color.Gray{Y: 200})
		})

		// メインコンテンツ (1行目、1列目)
		b.LabelAt(1, 1, 1, 1, func(l *widget.LabelBuilder) {
			l.Text("Content (Weight 1)").
				BackgroundColor(color.Gray{Y: 220})
		})

		// 追加情報 (1行目、2列目)
		b.LabelAt(1, 2, 1, 1, func(l *widget.LabelBuilder) {
			l.Text("More Info (Weight 2)").
				BackgroundColor(color.Gray{Y: 210})
		})

		// フッター (2行目、0列目から3列にまたがる)
		b.ButtonAt(2, 0, 1, 3, func(btn *widget.ButtonBuilder) {
			btn.Text("Footer Button (Spans 3 Columns, 30px Fixed Height)").
				AddOnClick(func(e *event.Event) {
					log.Println("Advanced Grid Footer clicked!")
				})
		})
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
				AbsolutePosition(50, 50).
				WrapText(true) // 折り返しを有効化
		})

		// 中景
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Middle Layer").
				Size(300, 200).
				BackgroundColor(color.RGBA{R: 100, G: 150, B: 200, A: 255}).
				TextColor(color.White).
				AbsolutePosition(100, 100).
				WrapText(true) // 折り返しを有効化
		})

		// 前景
		b.Button(func(btn *widget.ButtonBuilder) {
			btn.Text("Foreground Button").
				Size(150, 50).
				AbsolutePosition(175, 175).
				AddOnClick(func(e *event.Event) { log.Println("Foreground button clicked!") }).
				WrapText(true) // Buttonでも有効
		})

		// ZStackの範囲外（クリッピングされる場合）
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Partially Visible").
				BackgroundColor(color.RGBA{R: 200, G: 100, B: 100, A: 255}).
				AbsolutePosition(480, 200).
				WrapText(true) // 折り返しを有効化
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
							AddOnClick(func(e *event.Event) {
								log.Printf("Clicked: Item %d", itemNumber)
								if detailTitleLabel != nil {
									detailTitleLabel.SetText(fmt.Sprintf("Details for Item %d", itemNumber))
								}
								if detailInfoLabel != nil {
									detailInfoLabel.SetText(fmt.Sprintf("Here you would see more detailed information about item number %d. This text is updated dynamically when you select an item from the list. It can be quite long, so text wrapping is essential here.", itemNumber))
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
					WrapText(true). // 詳細テキストの折り返しを有効化
					AssignTo(&detailInfoLabel)
			})
		})
	}).Build()
}