package collections

func InArray[T comparable](arr []T, e T) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}
