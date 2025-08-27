# Furoshiki UI for Ebitengine

Furoshikiは、Go言語によるEbitengine向けUIライブラリです。日本の風呂敷のように、シンプルさと柔軟性を両立し、直感的で宣言的なAPIによってUI構築を容易にすることを目指しています。

## 設計思想 (Design Philosophy)

Furoshikiは、以下の設計思想に基づいています。

-   **宣言的なAPI**: UIの構造を、手続き的なコードではなく、見たままの階層構造として記述できるAPIを提供します。これにより、コードの可読性が向上し、UIの全体像を把握しやすくなります。
-   **柔軟性と拡張性**: Web技術で広く受け入れられているFlexboxレイアウトをコアに採用し、レスポンシブで複雑なレイアウトを簡単に実現します。また、すべてのコンポーネントは共通のインターフェースを実装しており、Goの「Composition over Inheritance（継承より合成）」の哲学に基づき、独自のカスタムコンポーネントを容易に作成・組み込みできます。
-   **関心の分離**: ライブラリは `component`, `container`, `layout`, `style`, `event`, `theme` といった責務が明確なパッケージに分割されています。これにより、ライブラリの各機能が理解しやすくなり、メンテナンス性も向上します。
-   **直感的なビルダー**: `ui`パッケージの高レベルなビルダーとメソッドチェーンにより、少ないコードで流れるようにUIを構築できます。`BackgroundColor()`, `Padding()`, `Margin()`, `Border()` のようなヘルパーメソッドが、冗長なスタイル定義を不要にします。

## コアコンセプト (Core Concepts)

Furoshikiを理解するための主要な概念です。

-   **Widget**: UIを構成するすべての要素（ボタン、ラベルなど）の基本となるインターフェースです。位置、サイズ、スタイル、イベントハンドラなどの共通機能を提供します。
-   **Container**: 他の `Widget` を子要素として内包できる特殊なウィジェットです。`Layout` と組み合わせることで、子要素の配置を管理します。
-   **Layout**: `Container` 内の子要素をどのように配置するかを決定するロジックです。
    -   **FlexLayout**: CSS Flexboxにインスパイアされた、強力で柔軟なレイアウトシステム。アイテムの折り返し (`Wrap`) にも対応し、`HStack`と`VStack`のバックボーンです。
    -   **GridLayout**: 子要素を均等な格子状に配置します。シンプルな表形式のレイアウトに適しています。
    -   **AdvancedGridLayout**: CSS Gridにインスパイアされ、列・行ごとの可変サイズ指定（固定ピクセル、重み）や、セルの結合（colspan/rowspan）が可能です。複雑な表形式のレイアウトに適しています。
    -   **AbsoluteLayout**: 子要素を座標で自由に配置します。`ZStack`で使われ、UI要素の重ね合わせに適しています。
-   **テーマ**: `theme`パッケージは、アプリケーション全体のスタイルを一元管理する仕組みを提供します。これにより、UIの一貫性を保ちつつ、デザインの変更を容易にします。
-   **ビルダーパターン**: `NewButtonBuilder()` のようなビルダーを使い、メソッドチェーンでプロパティを設定することで、安全かつ流れるようにコンポーネントを構築します。
-   **宣言的UIヘルパー**: `ui.VStack()`, `ui.HStack()`, `ui.ZStack()` などの関数を使い、ネストされたUI構造を直感的に構築できます。

## 主な機能 (Key Features)

### 1. テーマによるスタイルの一元管理

`theme`パッケージにより、アプリケーション全体の色、フォント、ウィジェットごとのデフォルトスタイル（通常時、ホバー時、押下時など）を一箇所で定義できます。これにより、UIの一貫性を保ち、デザインの変更が容易になります。

### 2. 強力でレスポンシブなレイアウトシステム

`FlexLayout`, `GridLayout`, そして `AdvancedGridLayout` を使用することで、モダンなUIレイアウトを簡単に構築できます。

-   **Flexbox**: `Direction`, `Justify`, `AlignItems`, `Gap`, `Flex`値などをサポート。複数行にわたるアイテムの折り返し (`Wrap`) と、行間の揃え (`AlignContent`) を完全にサポートしているため、ウィンドウサイズに応じて変化するレスポンシブなレイアウトも実現可能です。
-   **Grid**: `Columns`, `Rows`, `HorizontalGap`, `VerticalGap` を指定して均等な格子状に配置。
-   **Advanced Grid**: `Columns`と`Rows`に固定ピクセル (`ui.Fixed(100)`) や重み (`ui.Weight(1)`) を指定でき、ウィジェットを複数のセルにまたがって (`colspan`, `rowspan`) 配置できます。

### 3. コンテンツに応じた動的なサイズ調整

ウィジェットは自身のコンテンツに応じてサイズを自動調整する機能を持ち、動的なUI構築を強力にサポートします。

-   **テキストの自動折り返し**: `.WrapText(true)` を設定するだけで、ウィジェットの幅を超えたテキストが自動的に折り返されます。
-   **高さの自動調整**: テキストが折り返されると、レイアウトシステムがその内容量に合わせてウィジェットの高さを自動的に計算・調整します。これにより、APIから取得した可変長のテキストなどもレイアウトを崩すことなく安全に表示できます。

### 4. 直感的で宣言的なUI構築

`ui`パッケージのヘルパー関数と、各ウィジェットのビルダーを組み合わせることで、非常に短いコードで複雑なUIを構築できます。

-   `VStack()`: 垂直方向に子要素を配置します。
-   `HStack()`: 水平方向に子要素を配置します。
-   `ZStack()`: 子要素を重ねて配置します。
-   `Grid()`: 子要素を均等な格子状に配置します。
-   `AdvancedGrid()`: セル結合や可変サイズの列・行を持つ高度なグリッドを構築します。
-   `Spacer()`: `HStack`や`VStack`内で利用可能なスペースを埋めるための伸縮可能な空白を追加します。

### 5. スクロールとクリッピング (`ScrollView`)

`ScrollView`ウィジェットを使用することで、コンテナの表示領域を超えるコンテンツをスクロール表示できます。子要素はコンテナの境界で自動的にクリッピング（切り抜き）されるため、複雑なリストや長文の表示に最適です。

### 6. 状態ベースのイベントシステム

`AddOnClick` のような一般的なUIイベントに加え、`Pressed`（押下時）や`Disabled`（無効時）といったウィジェットの状態が内部で管理されます。これにより、ユーザーの操作に対してよりリッチな視覚的フィードバックを簡単に実装できます。また、一つのイベントに対して複数のハンドラを登録することも可能です。

## 使い方 (Usage Example)

スタイルはテーマから自動的に適用されるため、UIの構造定義に集中できます。動的なコンテンツも `.WrapText(true)` のような設定で簡単かつ安全に扱えます。

```go
// main.go
package main

import (
	"fmt"
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

func NewGame() *Game {
	// 1. フォントを読み込み、テーマを準備
	appTheme := theme.GetCurrent()
	appTheme.SetDefaultFont(basicfont.Face7x13)
	theme.SetCurrent(appTheme)

	// 2. UIを宣言的に構築
	root, _ := ui.VStack(func(b *ui.FlexBuilder) {
		b.Size(400, 500).
			Padding(20).
			Gap(15).
			AlignItems(layout.AlignStretch). // Stretch to use full width
			BackgroundColor(appTheme.BackgroundColor)

		// タイトル
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("Furoshiki UI Demo").
				Size(0, 40). // Width is stretched
				BackgroundColor(appTheme.PrimaryColor).
				TextColor(color.White).
				TextAlign(style.TextAlignCenter)
		})

		// 説明文 (自動折り返し)
		b.Label(func(l *widget.LabelBuilder) {
			l.Text("This label contains a long text that will wrap automatically thanks to the WrapText property.").
				WrapText(true) // Enable text wrapping
		})

		// ボタン
		b.HStack(func(b *ui.FlexBuilder) {
			b.Size(0, 50).Gap(20)

			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("OK").
					Flex(1). // スペースを均等に分ける
					AddOnClick(func(e *event.Event) { fmt.Println("OK clicked") })
			})
			b.Button(func(btn *widget.ButtonBuilder) {
				btn.Text("Cancel").
					Flex(1).
					AddOnClick(func(e *event.Event) { fmt.Println("Cancel clicked") })
			})
		})

	}).Build()

	// (Ebitenゲームループの実行部分は省略)
}

注意点：レイアウトとAbsolutePosition

ウィジェットの .AbsolutePosition(x, y) メソッドは、親コンテナが AbsoluteLayout (主に ui.ZStack で作成) の場合にのみ有効です。

FlexLayout (VStack や HStack), GridLayout, AdvancedGridLayout の中では、子の位置はレイアウトシステムによって自動的に計算・管理されます。そのため、これらのレイアウト内で .AbsolutePosition() を使用しても設定は無視されるため効果はありません。これは意図された挙動です。
今後のロードマップ (Roadmap)

Furoshikiは、より表現力豊かで使いやすいライブラリを目指して、以下の機能開発を計画しています。
基本的なウィジェットの拡充

現在、基本的なUIを構築するためのウィジェットを提供していますが、より多様なアプリケーションに対応するため、以下のコンポーネントの追加を計画しています。

    TextInput: ユーザーからのテキスト入力を受け付けるフィールド。

    Image: 画像を表示するためのウィジェット。

    Checkbox: オン/オフを切り替えるチェックボックス。

    Slider: 特定の範囲から値を選択するためのスライダー。

    RadioButton: 複数の選択肢から一つを選ぶためのラジオボタン。

高度な機能

    相対的なサイズ指定: 親コンポーネントのサイズに対する割合（例: width: "50%"）でのサイズ指定機能の導入を検討します。これにより、さらに柔軟なレスポンシブデザインが可能になります。

    アニメーション対応: スタイルのプロパティ（色、サイズ、位置など）を時間経過で滑らかに変化させるための基本的な仕組みを導入し、より動的なUI表現を可能にすることを目指します。

    データバインディング: ウィジェットのプロパティ（例: Labelのテキスト）をアプリケーション側のデータモデルと自動的に同期させる仕組みを検討します。これにより、UIの状態管理が簡素化されます。