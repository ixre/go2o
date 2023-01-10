package regex

import "regexp"

var phoneRegex = regexp.MustCompile(`^(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\d{8})$`)

// 是否为手机
func IsPhone(s string) bool {
	return phoneRegex.Match([]byte(s))
}
