package utils

// IfThen は三項演算子のように動作します。
// 条件(cond)がtrueであればvtrueを、falseであればvfalseを返します。
func IfThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
