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
	root component.Widget // UIツリーのルート

	// AssignToで紐付けられ、動的に更新されるウィジェットへの参照
	detailTitleLabel *widget.Label
	detailInfoLabel  *widget.Label
}

// NewGame は新しいGameインスタンスを作成し、UIを構築します。
func NewGame() *Game {
	g := &Game{}

	// --- 1. テーマとフォントの準備 ---
	// アプリケーション全体のデフォルトテーマを取得します。
	appTheme := theme.GetCurrent()
	// basicfontはデモ用です。実際のアプリではより高品質なフォントを読み込むことを推奨します。
	appTheme.SetDefaultFont(basicfont.Face7x13)
	// 更新したテーマをアプリケーションの現在のテーマとして設定します。
	theme.SetCurrent(appTheme)

	// --- 2. UIの宣言的な構築 ---
	// ui.HStackを使用して、UI全体を左右に分割するレイアウトを作成します。
	root, err := ui.HStack(func(b *ui.Builder) {
		b.Size(800, 600).
			BackgroundColor(appTheme.BackgroundColor).
			Padding(10). // ウィンドウ全体のパディング
			Gap(10)      // 左右ペイン間のギャップ

		// --- 左ペイン：スクロールビュー ---
		b.ScrollView(func(sv *widget.ScrollViewBuilder) {
			sv.Size(250, 580). // 親の高さ(600) - 親のPadding(10*2) = 580
						Border(1, color.Gray{Y: 200})

			// ScrollViewのコンテンツとして、垂直スタック(VStack)を持つコンテナを設定します。
			// このコンテナがスクロール対象となります。
			content, _ := ui.VStack(func(b *ui.Builder) {
				b.Padding(8).Gap(5) // スクロール領域内のパディングとアイテム間のギャップ

				// 多数のボタンを動的に生成して、スクロールが必要な状況を作ります。
				for i := 1; i <= 50; i++ {
					itemNumber := i
					b.Button(func(btn *widget.ButtonBuilder) {
						btn.Text(fmt.Sprintf("Item %d", itemNumber)).
							Size(220, 30).
							// OnClickイベントハンドラを設定します。
							// ボタンがクリックされたら、右ペインのラベルのテキストを更新します。
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

			// 作成したコンテナをScrollViewのコンテンツとして設定します。
			sv.Content(content)
		})

		// --- 右ペイン：詳細表示エリア ---
		b.VStack(func(b *ui.Builder) {
			b.Flex(1). // 残りの水平スペースをすべて使用します
					Padding(10).
					Gap(10).
					Border(1, color.Gray{Y: 200}).
					AlignItems(layout.AlignStretch)

			// タイトルラベル
			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Details").
					Size(0, 30). // 高さは30, 幅は親に合わせる (Flex未設定のため)
					TextColor(color.White).
					BackgroundColor(appTheme.PrimaryColor).
					// AssignToメソッドを使い、ビルド中のラベルインスタンスへの参照を
					// Game構造体のフィールドに安全に格納します。
					AssignTo(&g.detailTitleLabel)
			})

			// 情報表示用ラベル
			b.Label(func(l *widget.LabelBuilder) {
				l.Text("Please select an item from the list on the left.").
					Flex(1). // 垂直方向の余ったスペースを埋めます
					TextAlign(style.TextAlignLeft).
					VerticalAlign(style.VerticalAlignTop).
					AssignTo(&g.detailInfoLabel)
			})

			// SpacerはFlexLayout内で余白を埋めるためのウィジェットです。
			// ここではFlex(1)のラベルがあるため実質的な効果はありませんが、使用例として示します。
			b.Spacer()

			// フッター
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
	// Ebitenから現在のマウスカーソル座標を取得します。
	cx, cy := ebiten.CursorPosition()
	// イベントディスパッチャのシングルトンインスタンスを取得します。
	dispatcher := event.GetDispatcher()

	// UIツリーのルートでヒットテストを実行し、カーソル直下のウィジェットを特定します。
	target := g.root.HitTest(cx, cy)
	// マウスイベント（ホバー、クリックなど）を適切なウィジェットにディスパッチします。
	dispatcher.Dispatch(target, cx, cy)

	// UIツリー全体の更新処理を呼び出します。
	// これにより、ダーティマークされたウィジェットの再レイアウトや状態更新が行われます。
	g.root.Update()
	return nil
}

// Draw はゲームを描画します。
func (g *Game) Draw(screen *ebiten.Image) {
	// 背景色で画面をクリアします。
	screen.Fill(color.RGBA{50, 50, 50, 255})
	// UIツリーのルートを描画します。これにより、すべての子孫ウィジェットが再帰的に描画されます。
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
