package component

// NOTE: このファイル内のインターフェースは、`ui`のような高レベルパッケージにおける
// ジェネリックヘルパー関数をサポートするために定義されています。これらは様々なビルダー実装に
// 期待される振る舞いを形式化し、型安全で再利用可能なUI構築ロジックを可能にします。

// BuilderInitializer はビルダーを初期化するための契約を定義します。
// これは基底の `component.Builder` によって実装されます。
type BuilderInitializer[T any, W Widget] interface {
	Init(self T, widget W)
}

// ErrorAdder はエラーを蓄積できるビルダーのための契約を定義します。
// これは基底の `component.Builder` によって実装されます。
type ErrorAdder interface {
	AddError(err error)
}

// BuilderFinalizer は最終的なウィジェットを生成できるビルダーのための契約を定義します。
// これは基底の `component.Builder` によって実装されます。
type BuilderFinalizer[W Widget] interface {
	Build() (W, error)
}

// WidgetContainer は子ウィジェットを受け入れることができるビルダーのための契約を定義します。
// このインターフェースは、`Container` ウィジェットをラップする高レベルのビルダー
// (例: `ui.FlexBuilder`) によって実装されることを意図しています。
type WidgetContainer interface {
	// AddChild は、ラップしているコンテナにウィジェットを追加するためにビルダーに期待されます。
	// メソッドチェーンのための戻り値の型は、具象ビルダーの実装によって処理されます。
	AddChild(child Widget)
}