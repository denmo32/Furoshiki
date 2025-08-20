# Furoshiki UI for Ebitengine

Furoshikiは、Go言語によるEbitengine向けUIライブラリです。日本の風呂敷のように、シンプルさと柔軟性を両立し、直感的で宣言的なAPIによってUI構築を容易にすることを目指しています。

## 設計思想 (Design Philosophy)

Furoshikiは、以下の設計思想に基づいています。

-   **宣言的なAPI**: UIの構造を、手続き的なコードではなく、見たままの階層構造として記述できるAPIを提供します。これにより、コードの可読性が向上し、UIの全体像を把握しやすくなります。
-   **柔軟性と拡張性**: Web技術で広く受け入れられているFlexboxレイアウトをコアに採用し、レスポンシブで複雑なレイアウトを簡単に実現します。また、すべてのコンポーネントは共通のインターフェースを実装しており、Goの「Composition over Inheritance（継承より合成）」の哲学に基づき、独自のカスタムコンポーネントを容易に作成・組み込みできます。
-   **関心の分離**: ライブラリは `component`, `container`, `layout`, `style`, `event` といった責務が明確なパッケージに分割されています。これにより、ライブラリの各機能が理解しやすくなり、メンテナンス性も向上します。
-   **直感的なビルダー**: `ui`パッケージの高レベルなビルダーとメソッドチェーンにより、少ないコードで流れるようにUIを構築できます。`BackgroundColor()`や`Padding()`のようなヘルパーメソッドが、冗長なスタイル定義を不要にします。

## コアコンセプト (Core Concepts)

Furoshikiを理解するための主要な概念です。

-   **Widget**: UIを構成するすべての要素（ボタン、ラベルなど）の基本となるインターフェースです。位置、サイズ、スタイル、イベントハンドラなどの共通機能を提供します。
-   **Container**: 他の `Widget` を子要素として内包できる特殊なウィジェットです。`Layout` と組み合わせることで、子要素の配置を管理します。
-   **Layout**: `Container` 内の子要素をどのように配置するかを決定するロジックです。
    -   **FlexLayout**: CSS Flexboxにインスパイアされた、強力で柔軟なレイアウトシステム。`HStack`と`VStack`のバックボーンです。
    -   **AbsoluteLayout**: 子要素を座標で自由に配置します。`ZStack`で使われ、UI要素の重ね合わせに適しています。
    -   **GridLayout**: 子要素を格子状に配置します。設定画面やインベントリに適しています。
-   **ビルダーパターン**: `NewButtonBuilder()` のようなビルダーを使い、メソッドチェーンでプロパティを設定することで、安全かつ流れるようにコンポーネントを構築します。
-   **宣言的UIヘルパー**: `ui.VStack()`, `ui.HStack()`, `ui.ZStack()` などの関数を使い、ネストされたUI構造を直感的に構築できます。

## 主な機能 (Key Features)

### 1. 強力なレイアウトシステム

`FlexLayout`と`GridLayout`を使用することで、モダンなUIレイアウトを簡単に構築できます。

-   **Flexbox**: `Direction`, `Justify`, `AlignItems`, `Gap`, `Flex`値などをサポート。新しく追加された `Spacer` ウィジェットを使えば、動的な空白の管理も容易です。
-   **Grid**: `Columns`, `Rows`, `HorizontalGap`, `VerticalGap` を指定して格子状に配置。

### 2. 直感的で宣言的なUI構築

`ui`パッケージのヘルパー関数と、各ウィジェットのビルダーを組み合わせることで、非常に短いコードで複雑なUIを構築できます。

-   `VStack()`: 垂直方向に子要素を配置します。
-   `HStack()`: 水平方向に子要素を配置します。
-   `ZStack()`: 子要素を重ねて配置します。
-   `Grid()`: 子要素を格子状に配置します。
-   `Spacer()`: `HStack`や`VStack`内で利用可能なスペースを埋めるための伸縮可能な空白を追加します。

### 3. シンプルなイベントシステム

`OnClick` のような一般的なUIイベントを、`event.Event`オブジェクトを受け取るコールバック関数としてコンポーネントに簡単に関連付けることができます。

### 4. 動的なUI操作

`AddChild()` や `RemoveChild()` メソッドを使って、実行中にUIツリーの構造を安全に変更できます。また、`SetVisible(bool)` を使ってコンポーネントの表示・非表示を切り替えることも可能です。

## 使い方 (Usage Example)

`ui` パッケージのヘルパー関数と、ビルダーのスタイルヘルパーメソッドを組み合わせることで、UIの構造と見た目を極めて直感的に記述できます。

```go
// 使用例
root, _ := ui.VStack(func(b *ui.ContainerBuilder) {
    b.Size(400, 500).      // コンテナ全体のサイズ
      Padding(20).         // 全体のパディング
      Gap(15).             // 子要素間のギャップ
      AlignItems(layout.AlignCenter).
      BackgroundColor(color.RGBA{R: 245, G: 245, B: 245, A: 255})

    // タイトル
    b.Label(func(l *widget.LabelBuilder) {
        l.Text("Furoshiki UI Demo").
          Size(360, 40).
          BackgroundColor(color.RGBA{R: 70, G: 130, B: 180, A: 255}).
          TextColor(color.White).
          TextAlign(style.TextAlignCenter).
          VerticalAlign(style.VerticalAlignMiddle).
          BorderRadius(8)
    })

    // ボタンエリア
    b.HStack(func(b *ui.ContainerBuilder) {
        b.Size(360, 50).Gap(20)

        b.Button(func(btn *widget.ButtonBuilder) {
            btn.Text("OK").
              Flex(1). // スペースを均等に分ける
              OnClick(func(e event.Event) { fmt.Println("OK clicked") })
        })
        b.Button(func(btn *widget.ButtonBuilder) {
            btn.Text("Cancel").
              Flex(1). // スペースを均等に分ける
              OnClick(func(e event.Event) { fmt.Println("Cancel clicked") })
        })
    })

    // Spacerを使って残りのスペースを押し広げる
    b.Spacer()

    // フッターラベル
    b.Label(func(l *widget.LabelBuilder) {
        l.Text("Footer Area")
    })

}).Build()

注意点：レイアウトとPosition

    ウィジェットの .Position(x, y) メソッドは、親コンテナが AbsoluteLayout (主に ui.ZStack で作成) の場合にのみ有効です。

    FlexLayout (VStack や HStack) の中では、子の位置はレイアウトシステムによって自動的に計算・管理されるため、.Position() の設定は無視されます。これは意図された挙動です。

今後のロードマップ (Roadmap)

Furoshikiは、より表現力豊かで使いやすいライブラリを目指して、以下の機能開発を計画しています。

    [計画中] テーマ＆スタイルシートシステム: アプリケーション全体のデザインをThemeオブジェクトとして一元管理し、.Class("classname")でスタイルを適用する仕組みを導入します。

    [計画中] 基本ウィジェットの拡充: TextInput, Image, Checkbox, Slider など、基本的なUIコンポーネントを追加します。

    [計画中] クリッピングとスクロール: 描画領域を制限するクリッピング機能を実装し、ScrollViewコンテナを導入します。