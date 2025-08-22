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
    -   **FlexLayout**: CSS Flexboxにインスパイアされた、強力で柔軟なレイアウトシステム。`HStack`と`VStack`のバックボーンです。
    -   **AbsoluteLayout**: 子要素を座標で自由に配置します。`ZStack`で使われ、UI要素の重ね合わせに適しています。
    -   **GridLayout**: 子要素を格子状に配置します。設定画面やインベントリに適しています。
-   **テーマ**: `theme`パッケージは、アプリケーション全体のスタイルを一元管理する仕組みを提供します。これにより、UIの一貫性を保ちつつ、デザインの変更を容易にします。
-   **ビルダーパターン**: `NewButtonBuilder()` のようなビルダーを使い、メソッドチェーンでプロパティを設定することで、安全かつ流れるようにコンポーネントを構築します。
-   **宣言的UIヘルパー**: `ui.VStack()`, `ui.HStack()`, `ui.ZStack()` などの関数を使い、ネストされたUI構造を直感的に構築できます。

## 主な機能 (Key Features)

### 1. テーマによるスタイルの一元管理

新しく導入された`theme`パッケージにより、アプリケーション全体の色、フォント、ウィジェットごとのデフォルトスタイル（通常時、ホバー時、押下時など）を一箇所で定義できます。これにより、UIの一貫性を保ち、デザインの変更が容易になります。

### 2. 強力なレイアウトシステム

`FlexLayout`と`GridLayout`を使用することで、モダンなUIレイアウトを簡単に構築できます。

-   **Flexbox**: `Direction`, `Justify`, `AlignItems`, `Gap`, `Flex`値などをサポート。`Spacer` ウィジェットを使えば、動的な空白の管理も容易です。
-   **Grid**: `Columns`, `Rows`, `HorizontalGap`, `VerticalGap` を指定して格子状に配置。

### 3. 直感的で宣言的なUI構築

`ui`パッケージのヘルパー関数と、各ウィジェットのビルダーを組み合わせることで、非常に短いコードで複雑なUIを構築できます。

-   `VStack()`: 垂直方向に子要素を配置します。
-   `HStack()`: 水平方向に子要素を配置します。
-   `ZStack()`: 子要素を重ねて配置します。
-   `Grid()`: 子要素を格子状に配置します。
-   `Spacer()`: `HStack`や`VStack`内で利用可能なスペースを埋めるための伸縮可能な空白を追加します。

### 4. 状態ベースのイベントシステム

`OnClick` のような一般的なUIイベントに加え、`Pressed`（押下時）や`Disabled`（無効時）といったウィジェットの状態が内部で管理されます。これにより、ユーザーの操作に対してよりリッチな視覚的フィードバックを簡単に実装できます。

### 5. 動的なUI操作

`AddChild()` や `RemoveChild()` メソッドを使って、実行中にUIツリーの構造を安全に変更できます。また、`SetVisible(bool)` や `SetDisabled(bool)` を使ってコンポーネントの表示・非表示や有効・無効を切り替えることも可能です。

## 使い方 (Usage Example)

新しく導入された`theme`パッケージにより、UIの構築がさらにシンプルになりました。スタイルはテーマから自動的に適用されるため、UIの構造定義に集中できます。

```go
// main.go
func main() {
    // 1. フォントを読み込み、テーマを準備
    mplusFont := loadFont() // (フォント読み込み処理は省略)
    appTheme := theme.GetCurrent()
    appTheme.SetDefaultFont(mplusFont)
    theme.SetCurrent(appTheme)

    // 2. UIを宣言的に構築
    root, _ := ui.VStack(func(b *ui.ContainerBuilder) {
        b.Size(400, 500).
          Padding(20).
          Gap(15).
          AlignItems(layout.AlignCenter).
          BackgroundColor(appTheme.BackgroundColor)

        // タイトル (テーマの色を部分的に上書きし、枠線を追加)
        b.Label(func(l *widget.LabelBuilder) {
            l.Text("Furoshiki UI Demo").
              Size(360, 40).
              BackgroundColor(appTheme.PrimaryColor).
              TextColor(color.White).
              TextAlign(style.TextAlignCenter).
              VerticalAlign(style.VerticalAlignMiddle).
              BorderRadius(8).
              Border(1, color.Gray{Y: 180})
        })

        // ボタン (スタイルはテーマから自動適用)
        b.HStack(func(b *ui.ContainerBuilder) {
            b.Size(360, 50).Gap(20)

            b.Button(func(btn *widget.ButtonBuilder) {
                btn.Text("OK").
                  Flex(1). // スペースを均等に分ける
                  OnClick(func(e event.Event) { fmt.Println("OK clicked") })
            })
            b.Button(func(btn *widget.ButtonBuilder) {
                btn.Text("Cancel").
                  Flex(1).
                  OnClick(func(e event.Event) { fmt.Println("Cancel clicked") })
            })
        })

        b.Spacer()

        b.Label(func(l *widget.LabelBuilder) {
            l.Text("Footer Area")
        })

    }).Build()

    // (Ebitenゲームループの実行部分は省略)
}
```

## 注意点：レイアウトとAbsolutePosition

ウィジェットの `.AbsolutePosition(x, y)` メソッドは、親コンテナが `AbsoluteLayout` (主に `ui.ZStack` で作成) の場合にのみ有効です。

`FlexLayout` (`VStack` や `HStack`) や `GridLayout` の中では、子の位置はレイアウトシステムによって自動的に計算・管理されます。そのため、これらのレイアウト内で `.AbsolutePosition()` を使用しても設定は無視されるため効果はありません。これは意図された挙動です。

## 今後のロードマップ (Roadmap)

Furoshikiは、より表現力豊かで使いやすいライブラリを目指して、以下の機能開発を計画しています。

-   **[完了] テーマ＆スタイルシートシステム**: アプリケーション全体のデザインを`Theme`オブジェクトとして一元管理し、ウィジェットのスタイルを自動で適用するシステムを導入しました。これにより、UIの一貫性を保ちつつ、コードの記述量を削減できます。

-   **[計画中] 基本ウィジェットの拡充**: `TextInput`, `Image`, `Checkbox`, `Slider` など、基本的なUIコンポーネントを追加します。

-   **[計画中] クリッピングとスクロール**: 描画領域を制限するクリッピング機能を実装し、`ScrollView`コンテナを導入します。