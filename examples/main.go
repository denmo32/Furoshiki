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
	// 背景色はルートコンテナのスタイルで設定するため、ここでのFillは不要になります。
	// screen.Fill(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	g.root.Draw(screen)

	// FPSなどのデバッグ情報を表示
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

// Layout はゲームの画面サイズを返します
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 500
}

var mplusFont font.Face

// アプリケーション全体で共有する基本スタイルを定義
var baseTextStyle style.Style
var baseButtonStyle style.Style
var steelBlue = color.RGBA{R: 70, G: 130, B: 180, A: 255}
var lightGray = color.RGBA{R: 220, G: 220, B: 220, A: 255}
var darkGray = color.RGBA{R: 105, G: 105, B: 105, A: 255}

func main() {
	// --- フォントの読み込み ---
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

	// --- スタイルの事前定義 ---
	baseTextStyle = style.Style{
		Font: style.PFont(mplusFont),
	}
	baseButtonStyle = style.Merge(baseTextStyle, style.Style{
		Background:  style.PColor(lightGray),
		TextColor:   style.PColor(color.Black),
		BorderColor: style.PColor(darkGray),
		BorderWidth: style.PFloat32(1),
		Padding:     style.PInsets(style.Insets{Top: 5, Right: 10, Bottom: 5, Left: 10}),
	})

	// --- UIの構築 ---
	// [大幅改良] uiパッケージの宣言的ヘルパーと新しいスタイルAPIで、より直感的にUIを構築します。
	root, err := ui.VStack(func(b *ui.ContainerBuilder) {
		b.Size(400, 500). // ウィンドウサイズに合わせる
			Padding(20).    // 全体にパディング
			Gap(15).        // 子要素間のギャップ
			Justify(layout.AlignStart).
			AlignItems(layout.AlignCenter).
			BackgroundColor(color.RGBA{R: 245, G: 245, B: 245, A: 255})

		// 1. タイトルラベル
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Furoshiki UI Demo").
				Style(baseTextStyle). // フォントスタイルを適用
				Size(360, 40).
				BackgroundColor(steelBlue).
				TextColor(color.White).
				TextAlign(style.TextAlignCenter).
				VerticalAlign(style.VerticalAlignMiddle).
				BorderRadius(8)
		})

		// 2. ボタンを配置する水平コンテナ
		b.HStack(func(b *ui.ContainerBuilder) {
			b.Size(360, 50).
				Gap(20) // ボタン間のギャップ

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("OK").
					Style(baseButtonStyle).
					Flex(1). // Flex=1でスペースを均等に分け合う
					HoverStyle(style.Style{Opacity: style.PFloat64(0.8)}).
					OnClick(func(e event.Event) {
						fmt.Printf("OK button clicked at (%d, %d)!\n", e.X, e.Y)
					})
			})

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Cancel").
					Style(baseButtonStyle).
					Flex(1). // Flex=1でスペースを均等に分け合う
					HoverStyle(style.Style{Opacity: style.PFloat64(0.8)}).
					OnClick(func(_ event.Event) {
						fmt.Println("Cancel button clicked!")
					})
			})
		})

		// 3. 重ね合わせコンテナ (ZStack)
		b.ZStack(func(b *ui.ContainerBuilder) {
			b.Size(360, 100).
				BackgroundColor(lightGray).
				BorderRadius(4)

			// ZStack内の要素はPositionで相対位置を指定
			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Overlapping Content").
					Style(baseTextStyle).
					Position(20, 15) // コンテナ左上からの相対位置
			})

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Overlay Button").
					Style(baseButtonStyle).
					BorderRadius(20). // 丸いボタン
					Position(180, 50).
					Size(160, 40).
					OnClick(func(_ event.Event) {
						fmt.Println("Overlay button clicked!")
					})
			})
		})

		// 4. [新機能] Spacerを使って残りの垂直スペースを埋める
		b.Spacer()

		// 5. グリッドレイアウトのデモ
		b.Grid(func(g *ui.GridContainerBuilder) {
			g.Size(360, 120).
				Columns(3).
				HorizontalGap(10).
				VerticalGap(10)

			// 3x2のグリッドに6つのボタンを配置
			for i := 0; i < 6; i++ {
				buttonIndex := i + 1
				g.Button(func(btn *widget.ButtonBuilder) {
					btn.Text(fmt.Sprintf("Grid %d", buttonIndex)).
						Style(baseButtonStyle).
						OnClick(func(_ event.Event) {
							fmt.Printf("Grid button %d clicked!\n", buttonIndex)
						})
				})
			}
		})

	}).Build()

	if err != nil {
		log.Fatalf("Failed to build UI: %v", err)
	}

	game := &Game{
		root:       root,
		dispatcher: event.GetDispatcher(),
	}

	ebiten.SetWindowSize(400, 500)
	ebiten.SetWindowTitle("Furoshiki UI Demo (Improved)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Ebiten run game failed: %v", err)
	}
}