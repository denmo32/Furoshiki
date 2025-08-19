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
	dispatcher *event.Dispatcher // [追加] イベントディスパッチャへの参照を追加
}

// Update はゲームの状態を更新します
// [改善] イベントディスパッチャを呼び出し、UIイベントを処理するように変更します。
func (g *Game) Update() error {
	// 1. UIツリーの状態更新とレイアウト計算
	g.root.Update()

	// 2. マウスカーソル下のウィジェットを特定し、イベントをディスパッチ
	cx, cy := ebiten.CursorPosition()
	target := g.root.HitTest(cx, cy) // ヒットしたウィジェットを取得
	// component.Widgetはevent.EventTargetインターフェースを構造的に満たすため、直接渡せます。
	g.dispatcher.Dispatch(target, cx, cy) // ディスパッチャに処理を委譲

	return nil
}

// Draw はゲームを描画します
// [修正] メソッドレシーバーのタイプミス (g.Game) を正しい構文 (g *Game) に修正しました。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	g.root.Draw(screen)
	ebitenutil.DebugPrint(screen, "Furoshiki UI Demo")
}

// Layout はゲームの画面サイズを返します
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// [更新] GridLayoutのデモを追加したため、高さを増やします
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

	// [修正] style.Styleのフィールドがポインタになったため、&でアドレスを渡すように変更
	// 全ウィジェットに適用する基本スタイルを定義
	baseStyle := style.Style{
		Font: &mplusFont,
	}

	root, err := ui.VStack(func(b *ui.ContainerBuilder) {
		b.Padding(20).Gap(15).
			Justify(layout.AlignStart).
			AlignItems(layout.AlignCenter).
			Size(400, 450) // [更新] 高さを増やします

		// タイトルラベル
		b.Label(func(l *widget.LabelBuilder) {
			l.Style(baseStyle) // 基本スタイルを適用
			l.Text("Furoshiki UI Demo")
			l.Size(300, 40)
			// 色のスタイルをマージ
			// [修正] ui.Styleヘルパーが返すStyleもポインタフィールドを持つため、そのままマージ可能
			l.Style(ui.Style(color.RGBA{R: 70, G: 130, B: 180, A: 255}, color.White))
		})

		// ボタンを水平に配置
		b.HStack(func(b *ui.ContainerBuilder) {
			b.Gap(20).Size(300, 50)

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle) // 基本スタイルを適用
				btn.Text("OK")
				btn.Size(100, 40)
				btn.OnClick(func() {
					fmt.Println("OK button clicked!")
				})
			})

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle) // 基本スタイルを適用
				btn.Text("Cancel")
				btn.Size(100, 40)
				btn.OnClick(func() {
					fmt.Println("Cancel button clicked!")
				})
			})
		})

		// 重ね合わせコンテナ
		b.ZStack(func(b *ui.ContainerBuilder) {
			b.Size(300, 100) // [更新] 高さを調整
			// コンテナ自体に背景色を設定
			// [修正] color.Color 型の変数に具象型の値を入れることで、そのポインタを *color.Color として渡せるようにします。
			bgColor := color.Color(color.RGBA{R: 220, G: 220, B: 220, A: 255})
			b.Style(style.Style{
				Background: &bgColor,
			})

			// 前景としてボタンを一つだけ追加
			// [改善] 角丸とホバー時の不透明度スタイルを追加して、描画機能のデモを行います。
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Style(baseStyle) // 基本スタイルを適用

				// 角丸スタイルを追加
				radius := float32(8.0)
				btn.Style(style.Style{BorderRadius: &radius})

				btn.Text("Overlay Button")
				btn.Position(90, 35) // [更新] Y座標を調整
				btn.Size(120, 30)
				btn.OnClick(func() {
					fmt.Println("Overlay button clicked!")
				})

				// ホバー時のスタイル（少し半透明にする）
				opacity := 0.8
				btn.HoverStyle(style.Style{Opacity: &opacity})
			})
		})

		// [追加] グリッドレイアウトのデモ
		b.Grid(func(g *ui.GridContainerBuilder) {
			// [修正] メソッドチェーンが正しく動作するように修正
			g.Size(300, 120).
				Columns(3).
				HorizontalGap(10).
				VerticalGap(10)

			// 3x2のグリッドに6つのボタンを配置
			for i := 0; i < 6; i++ {
				// ループ変数をキャプチャして、各ボタンが正しい番号を持つようにする
				buttonIndex := i + 1
				g.Button(func(btn *widget.ButtonBuilder) {
					btn.Style(baseStyle)
					btn.Text(fmt.Sprintf("Grid %d", buttonIndex))
					// サイズはGridLayoutによって自動的に設定されるため、ここでは指定しない
					btn.OnClick(func() {
						fmt.Printf("Grid button %d clicked!\n", buttonIndex)
					})
				})
			}
		})

	}).Build()

	if err != nil {
		panic(err)
	}

	// [修正] Game構造体を初期化する際に、イベントディスパッチャも取得
	game := &Game{
		root:       root,
		dispatcher: event.GetDispatcher(),
	}

	ebiten.SetWindowSize(400, 450) // [更新] 高さを増やします
	ebiten.SetWindowTitle("Furoshiki UI Demo")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}