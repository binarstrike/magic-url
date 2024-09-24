package utils

// If returns t if condition is true, otherwise it returns f.
//
//	If(5 < 10, 100, 500) // return 100
func If[T any](condition bool, t, f T) T {
	if condition {
		return t
	}

	return f
}
