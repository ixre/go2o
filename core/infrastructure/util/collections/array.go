package collections

// 是否包含元素
func InArray[T comparable](arr []T, e T) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

// 是否包含任意元素
func AnyArray[T comparable](arr []T, f func(e T) bool) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

// 在数组中筛选
func FilterArray[T comparable](arr []T, f func(e T) bool) []T {
	ret := make([]T, 0)
	for _, v := range arr {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// MapList 映射列表
func MapList[T any, K any](arr []T, f func(T) K) []K {
	var ret = make([]K, len(arr))
	for i := range arr {
		ret[i] = f(arr[i])
	}
	return ret
}
