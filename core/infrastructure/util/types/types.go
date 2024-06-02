package types

// Ternary returns the first argument if cond is true, otherwise it returns the second.
func Ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}
