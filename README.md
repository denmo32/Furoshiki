Furoshiki UI for Ebitengine
Furoshikiは、Go言語によるEbitengine向けUIライブラリです。日本の風呂敷のように、シンプルさと柔軟性を両立し、直感的で宣言的なAPIによってUI構築を容易にすることを目指しています。

設計思想 (Design Philosophy)
Furoshikiは、以下の設計思想に基づいています。
    宣言的なAPI: UIの構造を、手続き的なコードではなく、見たままの階層構造として記述できるAPIを提供します。これにより、コードの可読性が向上し、UIの全体像を把握しやすくなります。
    柔軟性と拡張性: Web技術で広く受け入れられているFlexboxレイアウトをコアに採用し、レスポンシブで複雑なレイアウトを簡単に実現します。また、すべてのコンポーネントは共通のインターフェースを実装しており、Goの「Composition over Inheritance（継承より合成）」の哲学に基づき、独自のカスタムコンポーネントを容易に作成・組み込みできます。
    関心の分離: ライブラリは component, container, layout, style, event といった責務が明確なパッケージに分割されています。これにより、ライブラリの各機能が理解しやすくなり、メンテナンス性も向上します。
    段階的な抽象化: ライブラリは、柔軟性の高い低レベルなAPI（プリミティブ）と、それらを組み合わせて作られた利便性の高い高レベルなAPIの2層構造を提供することを目指します。ユーザーは、簡単なUIは高レベルAPIで迅速に構築し、複雑な要件には低レベルAPIで対応するなど、必要に応じて抽象度のレベルを選択できます。

コアコンセプト (Core Concepts)
Furoshikiを理解するための主要な概念です。
    Component: UIを構成するすべての要素（ボタン、ラベルなど）の基本となるインターフェースです。位置、サイズ、スタイル、イベントハンドラなどの共通機能を提供します。
    Container: 他の Component を子要素として内包できる特殊なコンポーネントです。Layout と組み合わせることで、子要素の配置を管理します。
    Layout: Container 内の子要素をどのように配置するかを決定するロジックです。
        FlexLayout: CSS Flexboxにインスパイアされた、強力で柔軟なレイアウトシステム。
        AbsoluteLayout: 子要素を絶対座標で自由に配置します。ダイアログやオーバーレイに適しています。
    Builder Pattern: NewButtonBuilder() のようなビルダーを使い、メソッドチェーンでプロパティを設定することで、安全かつ流れるようにコンポーネントを構築します。
    Style: style.Style 構造体を通じて、背景色、境界線、パディング、マージンといったコンポーネントの見た目を定義します。

主な機能 (Key Features)
1. 強力なFlexboxレイアウト
FlexLayout を使用することで、モダンなUIレイアウトを簡単に構築できます。
    Direction: 子要素を縦 (DirectionColumn) または横 (DirectionRow) に並べます。
    Justify: 主軸方向の揃え位置（先頭、中央、末尾）を指定します。
    AlignItems: 交差軸方向の揃え位置（先頭、中央、末尾、引き伸ばし）を指定します。
    Gap: 子要素間の間隔を設定します。
    Flex値: 子要素に Flex(1) のように値を設定すると、親コンテナのサイズ変更に応じて自動的に伸縮します。

2. イベントシステム
OnClick, OnMouseEnter, OnMouseLeave といった一般的なUIイベントを、コンポーネントに簡単に関連付けることができます。

3. 動的なUI操作
AddChild() や RemoveChild() メソッドを使って、実行中にUIツリーの構造を安全に変更できます。また、SetVisible(bool) を使ってコンポーネントの表示・非表示を切り替えることも可能です。非表示のコンポーネントは、更新・描画・レイアウト計算の対象から除外されます。

今後のロードマップ (Roadmap)
Furoshikiは、より表現力豊かで使いやすいライブラリを目指して、以下の機能開発を計画しています。
    [計画中] テーマ＆スタイルシートシステム: アプリケーション全体のデザインをThemeオブジェクトとして一元管理し、.Class("classname")でスタイルを適用する仕組みを導入します。これにより、スタイル定義とUI構造を分離し、コードの簡潔化とデザインの一貫性を向上させます。
    [計画中] 宣言的UIヘルパー (Functional Builders): UIの階層構造をコードのネストで表現できる高レベルAPIを導入し、より直感的で宣言的なUI構築を可能にします。

// 将来の理想的なコード例
root, _ := ui.VStack(func(b *ui.Builder) { // 縦積みコンテナ
    b.Padding(10).Gap(5)
    b.Label(func(l *ui.LabelBuilder) { l.Text("Title") })
    b.Button(func(btn *ui.ButtonBuilder) {
        btn.Text("Submit").Class("primary")
    })
}).Build()

	[計画中] コアウィジェットの拡充: TextInput, Image, Checkbox, Slider など、基本的なUIコンポーネントを追加します。
	[計画中] クリッピング: 描画領域を制限するクリッピング機能を実装します。