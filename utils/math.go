package utils

// IfThen は三項演算子のように動作します。
func IfThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Clamp は値を指定された最小値と最大値の範囲内に収めます。
// Go 1.22で cmp.Clamp が導入されるまでの代替として使用できます。
func Clamp[T ~int | ~float64 | ~float32](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// 最新版のGoでは組み込みのmax/min関数が利用可能なため、独自定義の関数は削除しました。
// 旧バージョンとの互換性維持は想定していません。