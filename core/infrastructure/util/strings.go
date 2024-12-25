package util

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// 获取汉字首字母
func GetHansFirstLetter(str string) string {
	arr := pinyin.LazyPinyin(str, pinyin.Args{})
	return string([]rune{unicode.ToUpper(rune(arr[0][0]))})
}

// 转换为版本号
func IntVersion(s string) int {
	arr := strings.Split(s, ".")
	for i, v := range arr {
		if l := len(v); l < 3 {
			arr[i] = strings.Repeat("0", 3-l) + v
		}
	}
	intVer, err := strconv.Atoi(strings.Join(arr, ""))
	if err != nil {
		panic(err)
	}
	return intVer
}
