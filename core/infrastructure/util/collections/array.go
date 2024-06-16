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

// 在数组中查找单个元素
func FindArray[T interface{}](arr []T, f func(e T) bool) (t T) {
	for _, v := range arr {
		if f(v) {
			return v
		}
	}
	return
}

// MapList 映射列表
func MapList[T any, M any](arr []T, f func(T) M) []M {
	var ret = make([]M, len(arr))
	for i := range arr {
		ret[i] = f(arr[i])
	}
	return ret
}

// Map 映射
func Map[K, I comparable, V, M any](m map[K]V, f func(K, V) (I, M)) map[I]M {
	var ret = make(map[I]M, len(m))
	for k, v := range m {
		i, m := f(k, v)
		ret[i] = m
	}
	return ret
}
