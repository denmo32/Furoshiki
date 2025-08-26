package utils

// IfThen は三項演算子のように動作します。
func IfThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// 最新版のGoでは組み込みのmax関数が利用可能です。
// 独自定義ののmax関数は削除しました。