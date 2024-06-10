package validate

import "regexp"

var (
	userRegex  = regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
	emailRegex = regexp.MustCompile(`^[A-Za-z0-9_\-]+@[a-zA-Z0-9\\-]+(\.[a-zA-Z0-9]+)+$`)
	phoneRegex = regexp.MustCompile(`^(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\d{8})$`)
)

// 验证用户名是否合法
func IsUser(user string) bool {
	return userRegex.MatchString(user)
}

// 验证电子邮件是否合法
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// 验证手机号码是否合法
func IsPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}
