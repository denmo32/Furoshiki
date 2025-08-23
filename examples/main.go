package main

import (
	"fmt"
	"furoshiki/component"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/style"
	"furoshiki/theme"
	"furoshiki/ui"
	"furoshiki/widget"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Game はEbitenのゲームインターフェースを実装します
type Game struct {
	root *container.Container

	// デモ用の状態管理
	dynamicContainer   *container.Container
	buttonsToDisable   []component.Widget
	widgetCount        int
	areButtonsDisabled bool
}

// NewGame は新しいGameインスタンスを初期化して返します。
func NewGame() *Game {
	g := &Game{}

	// --- フォントの読み込み ---
	mplusFont := loadFont()

	// --- テーマの準備と設定 ---
	appTheme := theme.GetCurrent()
	appTheme.SetDefaultFont(mplusFont)
	// ボタンの色をデフォルトから少しカスタマイズ
	appTheme.Button.Normal.Background = style.PColor(color.RGBA{R: 230, G: 230, B: 240, A: 255})
	appTheme.Button.Hovered.Background = style.PColor(color.RGBA{R: 220, G: 220, B: 235, A: 255})
	theme.SetCurrent(appTheme)

	// --- UIの構築 ---
	// [改善] UI構築のフローを、新しい `AssignTo` メソッドを用いてより宣言的で一貫性の
	// あるものにリファクタリングしました。これにより、UIの階層構造を保ったまま、
	// 特定のウィジェットへの参照を安全に取得できます。
	// (※ このコードを動作させるには、まず component.Builder と ui.Builder に
	//      AssignTo メソッドを追加する必要があります)
	var okBtn, cancelBtn *widget.Button // 無効化対象のボタンへの参照を保持する変数

	root, err := ui.VStack(func(b *ui.Builder) {
		b.Size(400, 600).Padding(10).Gap(10).BackgroundColor(appTheme.BackgroundColor)

		// 1. タイトル
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Furoshiki Demo").Size(380, 40).BackgroundColor(appTheme.PrimaryColor).TextColor(color.White).TextAlign(style.TextAlignCenter).VerticalAlign(style.VerticalAlignMiddle).BorderRadius(8)
		})

		// 2. 動的ウィジェット操作パネル
		b.HStack(func(b *ui.Builder) {
			b.Size(380, 40).Gap(10)
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Add Widget").Flex(1).OnClick(func(e event.Event) { g.addWidget() })
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Remove Widget").Flex(1).OnClick(func(e event.Event) { g.removeWidget() })
			})
		})

		// 3. 動的ウィジェットが追加されるコンテナ
		// RelayoutBoundaryをtrueにすることで、このコンテナ内の変更が親に影響しなくなります。
		// [改善] 以前は一時的なビルダーを作成していましたが、ネストされたVStack内で
		// `AssignTo` を使うことで、コードがシンプルかつ直感的になります。
		b.VStack(func(b *ui.Builder) {
			b.Size(380, 200).
				Padding(10).
				Gap(5).
				BackgroundColor(appTheme.SecondaryColor).
				RelayoutBoundary(true).
				// [変更点] ClipChildren(true) を呼び出して、このコンテナの境界外にはみ出す
				// 子ウィジェット（動的に追加されるラベル）が描画されないようにします。
				ClipChildren(true).
				AssignTo(&g.dynamicContainer) // `AssignTo`でコンテナのインスタンスを直接取得
		})

		// 4. 無効化機能のデモ
		// [改善] `AssignTo` を使い、HStackの宣言的な構造の中でボタンへの参照を取得します。
		// これにより、一時的なビルダーやAddChildrenの呼び出しが不要になります。
		b.HStack(func(b *ui.Builder) {
			b.Size(380, 40).Gap(10)
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("OK").Flex(1).AssignTo(&okBtn)
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Cancel").Flex(1).AssignTo(&cancelBtn)
			})
		})

		// ボタンへの参照をスライスに格納します。
		// okBtnとcancelBtnは*widget.Button型ですが、component.Widgetインターフェースを満たすため、
		// このように代入できます。
		g.buttonsToDisable = []component.Widget{okBtn, cancelBtn}

		b.Button(func(btn *widget.ButtonBuilder) {
			btn.Text("Disable/Enable Buttons").Size(380, 40).OnClick(func(e event.Event) { g.toggleButtonsDisabled() })
		})

	}).Build()

	if err != nil {
		log.Fatalf("Failed to build UI: %v", err)
	}
	g.root = root
	return g
}

func (g *Game) addWidget() {
	g.widgetCount++
	// [変更点] ラベルのサイズを少し小さくして、はみ出しをより分かりやすくします
	newLabel, _ := widget.NewLabelBuilder().Text(fmt.Sprintf("Dynamic Label #%d", g.widgetCount)).Size(360, 25).BackgroundColor(color.White).Padding(5).Build()
	g.dynamicContainer.AddChild(newLabel)
}

func (g *Game) removeWidget() {
	children := g.dynamicContainer.GetChildren()
	if len(children) > 0 {
		lastChild := children[len(children)-1]
		g.dynamicContainer.RemoveChild(lastChild)
	}
}

func (g *Game) toggleButtonsDisabled() {
	g.areButtonsDisabled = !g.areButtonsDisabled
	for _, btn := range g.buttonsToDisable {
		btn.SetDisabled(g.areButtonsDisabled)
	}
}

// Update はゲームの状態を更新します
func (g *Game) Update() error {
	g.root.Update()
	dispatcher := event.GetDispatcher()
	cx, cy := ebiten.CursorPosition()
	target := g.root.HitTest(cx, cy)
	dispatcher.Dispatch(target, cx, cy)
	return nil
}

// Draw はゲームを描画します
func (g *Game) Draw(screen *ebiten.Image) {
	g.root.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, FPS: %0.2f, Widgets: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), g.widgetCount))
}

// Layout はゲームの画面サイズを返します
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 600
}

func main() {
	ebiten.SetWindowSize(400, 600)
	ebiten.SetWindowTitle("Furoshiki UI Demo (Enhanced)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatalf("Ebiten run game failed: %v", err)
	}
}

func loadFont() font.Face {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14, // 少し小さくして多くの情報を表示
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return mplusFont
}