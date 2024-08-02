package regex

import "regexp"

var phoneRegex = regexp.MustCompile(`^(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\d{8})$`)
var invalidCharsRegexp = regexp.MustCompile(`[#|$%&'\*\+,./:;\<=\>\?@\[\]^{\|}~]`)
var emailRegex = regexp.MustCompile(`^[A-Za-z0-9_\-]+@[a-zA-Z0-9\\-]+(\.[a-zA-Z0-9]+)+$`)

// 是否为手机
func IsPhone(s string) bool {
	return phoneRegex.Match([]byte(s))
}

// 是否包含特殊字符,并返回字符数组
func ContainInvalidChars(txt string) (bool, []string) {
	b := invalidCharsRegexp.Match([]byte(txt))
	if !b {
		return false, nil
	}
	return true, invalidCharsRegexp.FindAllString(txt, -1)
}

// 验证电子邮件是否合法
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

var (
	userRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
)

// 验证用户名是否合法
func IsUser(user string) bool {
	return userRegex.MatchString(user)
}
