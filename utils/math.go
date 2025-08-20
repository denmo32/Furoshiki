package utils

// Max は2つの整数のうち大きい方を返します。
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IfThen は三項演算子のように動作します。
// 条件(cond)がtrueであればvtrueを、falseであればvfalseを返します。
func IfThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}