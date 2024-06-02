package util

import (
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// 获取汉字首字母
func GetHansFirstLetter(str string) string {
	arr := pinyin.LazyPinyin(str, pinyin.Args{})
	return string([]rune{unicode.ToUpper(rune(arr[0][0]))})
}
