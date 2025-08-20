package main

import (
	"fmt"
	"furoshiki/container"
	"furoshiki/event"
	"furoshiki/layout"
	"furoshiki/style"
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
	root       *container.Container
	dispatcher *event.Dispatcher
}

// Update はゲームの状態を更新します
func (g *Game) Update() error {
	// 1. UIツリーの状態更新とレイアウト計算
	// この中で、ダーティマークされたコンテナのレイアウトが再計算されます。
	g.root.Update()

	// 2. マウスカーソル下のウィジェットを特定し、イベントをディスパッチします。
	cx, cy := ebiten.CursorPosition()
	target := g.root.HitTest(cx, cy) // ヒットしたウィジェットを取得
	// component.Widgetはevent.EventTargetインターフェースを構造的に満たすため、直接渡せます。
	g.dispatcher.Dispatch(target, cx, cy) // ディスパッチャにイベント処理を委譲

	return nil
}

// Draw はゲームを描画します
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	g.root.Draw(screen)
	ebitenutil.DebugPrint(screen, "Furoshiki UI Demo")
}

// Layout はゲームの画面サイズを返します
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 450
}

var mplusFont font.Face

func main() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	// [改良] スタイル設定にstyle.PFontヘルパーを使用
	baseStyle := style.Style{
		Font: style.PFont(mplusFont),
	}

	// uiパッケージのヘルパーを使って宣言的にUIを構築します
	root, err := ui.VStack(func(b *ui.ContainerBuilder) {
		b.Padding(20).Gap(15).
			Justify(layout.AlignStart).
			AlignItems(layout.AlignCenter).
			Size(400, 450)

		// タイトルラベル
		b.Label(func(l *widget.LabelBuilder) {
			l.Style(baseStyle) // 基本スタイルを適用
			l.Text("Furoshiki UI Demo")
			l.Size(300, 40)
			// 色のスタイルをマージします
			l.Style(ui.Style(color.RGBA{R: 70, G: 130, B: 180, A: 255}, color.White))
			// [改良] 新しいテキスト揃えAPIの使用例
			l.Style(style.Style{
				TextAlign:    style.PTextAlignType(style.TextAlignCenter),
				VerticalAlign: style.PVerticalAlignType(style.VerticalAlignMiddle),
			})
		})

		// ボタンを水平に配置
		b.HStack(func(b *ui.ContainerBuilder) {
			b.Gap(20).Size(300, 50)

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle)
				btn.Text("OK")
				btn.Size(100, 40)
				// [改良] OnClickの新しいシグネチャを使用。イベント引数eは無視できます。
				btn.OnClick(func(e event.Event) {
					fmt.Printf("OK button clicked at (%d, %d)!\n", e.X, e.Y)
				})
			})

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle)
				btn.Text("Cancel")
				btn.Size(100, 40)
				btn.OnClick(func(_ event.Event) { // 引数を `_` で無視する例
					fmt.Println("Cancel button clicked!")
				})
			})
		})

		// 重ね合わせコンテナ (ZStack)
		b.ZStack(func(b *ui.ContainerBuilder) {
			b.Size(300, 100)
			// [改良] style.PColorヘルパーを使い、一時変数が不要に
			b.Style(style.Style{
				Background: style.PColor(color.RGBA{R: 220, G: 220, B: 220, A: 255}),
			})

			// 前景としてボタンを一つ追加します。
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle)

				// [改良] PFloat32, PFloat64ヘルパーを使い、より直感的に
				btn.Style(style.Style{
					BorderRadius: style.PFloat32(8.0),
				})
				btn.Text("Overlay Button")
				btn.Position(90, 35) // ZStack内での相対位置
				btn.Size(130, 30)
				btn.OnClick(func(_ event.Event) {
					fmt.Println("Overlay button clicked!")
				})
				btn.HoverStyle(style.Style{
					Opacity: style.PFloat64(0.8),
				})
			})
		})

		// グリッドレイアウトのデモ
		b.Grid(func(g *ui.GridContainerBuilder) {
			g.Size(300, 120). // グリッドコンテナ自体のサイズ
						Columns(3).
						HorizontalGap(10).
						VerticalGap(10)

			// 3x2のグリッドに6つのボタンを配置
			for i := 0; i < 6; i++ {
				// ループ変数をキャプチャして、各ボタンが正しい番号を持つようにします
				buttonIndex := i + 1
				g.Button(func(btn *widget.ButtonBuilder) {
					btn.Style(baseStyle)
					btn.Text(fmt.Sprintf("Grid %d", buttonIndex))
					// サイズはGridLayoutによって自動的に設定されるため、ここでは指定しません
					btn.OnClick(func(_ event.Event) {
						fmt.Printf("Grid button %d clicked!\n", buttonIndex)
					})
				})
			}
		})

	}).Build()

	if err != nil {
		panic(err)
	}

	game := &Game{
		root:       root,
		dispatcher: event.GetDispatcher(),
	}

	ebiten.SetWindowSize(400, 450)
	ebiten.SetWindowTitle("Furoshiki UI Demo")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}