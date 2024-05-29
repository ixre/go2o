package collections

func InArray[T comparable](arr []T, e T) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

// MapList 映射列表
func MapList[T any, K any](arr []T, f func(T) K) []K {
	var ret = make([]K, len(arr))
	for i := range arr {
		ret[i] = f(arr[i])
	}
	return ret
}
