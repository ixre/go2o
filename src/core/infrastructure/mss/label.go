/**
 * Copyright 2015 @ S1N1 Team.
 * name : label
 * author : jarryliu
 * date : 2015-07-27 09:28
 * description :
 * history :
 */
package mss
import "regexp"

var reg = regexp.MustCompile("\\{([^\\}]+\\}")

// 翻译标签
func Transplate(c string,m map[string]string)string {
	return reg.ReplaceAllStringFunc(c, func(k string) {
		if v, ok := m[k]; ok {
			return v
		}
		return "{"+k+"}"
	})
}