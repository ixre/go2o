package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 解析模板中的参数
func ResolveMessage(templateContent string, param []string) string {
	//　替换字符标签{name}标签为{0}
	re := regexp.MustCompile(`\$*{(.+?)}`)
	holders := re.FindAllStringSubmatch(templateContent, -1)
	varPos := make(map[string]int)
	i := 0
	for _, v := range holders {
		s, k := v[0], v[1]
		if e, ok := varPos[k]; !ok {
			// 存储变量位置
			varPos[k] = i
		} else {
			// 替换已出现的模板变量
			templateContent = strings.Replace(templateContent, s, fmt.Sprintf("{%d}", e), 1)
			continue
		}
		if _, err := strconv.Atoi(k); err != nil {
			templateContent = strings.Replace(templateContent, s, fmt.Sprintf("{%d}", i), 1)
		}
		i++
	}
	// 替换值
	for k, v := range param {
		templateContent = strings.ReplaceAll(templateContent, fmt.Sprintf("{%d}", k), v)
	}
	return templateContent
}
