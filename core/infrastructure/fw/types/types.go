package types

import (
	"encoding/json"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Ternary returns the first argument if cond is true, otherwise it returns the second.
func Ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

func OrValue[T comparable](v T, or T) (t T) {
	if v == t {
		return or
	}
	return v
}

// 不包含前缀的较短的Title
func Title(str string) string {
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = cases.Title(language.Und).String(v)
	}
	return strings.Join(arr, "")
}

// Title 下划线转驼峰
func CamelTitle(str string, shortUpper bool) string {
	arr := strings.Split(str, "_")
	n := make([]string, len(arr)-1)
	for i, v := range arr[1:] {
		n[i] = cases.Title(language.Und).String(v)
	}
	arr = append(arr[:1], n...)
	return strings.Join(arr, "")
}

// DeepClone 深拷贝
func DeepClone[T any](v *T) (t *T) {
	dst := new(T)
	*dst = *v
	return dst
}

// 转换为Json对象
func ParseJSONObject(v interface{}) (map[string]interface{}, error) {
	var data []byte
	switch v := v.(type) {
	case string:
		data = []byte(v)
	default:
		md, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = md
	}
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err == nil {
		return result, nil
	}
	return nil, err
}
