package regex

import "regexp"

// 手机号正则
var PhoneRegex = regexp.MustCompile(`^(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\d{8})$`)
var ContainPhoneRegex = regexp.MustCompile(`(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\d{8})`)

// 特殊字符正则
var InvalidCharsRegexp = regexp.MustCompile(`[#|$%&'\*\+,./:;\<=\>\?@\[\]^{\|}~]`)

// 邮箱正则
var EmailRegex = regexp.MustCompile(`^[A-Za-z0-9_\-]+@[a-zA-Z0-9\\-]+(\.[a-zA-Z0-9]+)+$`)
var ContainEmailRegex = regexp.MustCompile(`[A-Za-z0-9_\-]+@[a-zA-Z0-9\\-]+(\.[a-zA-Z0-9]+)+`)

// 网址正则
var UrlRegex = regexp.MustCompile(`^(((ht|f)tps?):\/\/)?[\w-]+(\.[\w-]+)+([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?$`)
var ContainUrlRegex = regexp.MustCompile(`(((ht|f)tps?):\/\/)?[\w-]+(\.[\w-]+)+([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)

// 中国组织机构代码正则
var ChinaORGCodeRegex = regexp.MustCompile(`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`)
var ContainChinaORGCodeRegex = regexp.MustCompile(`[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}`)

// 银行卡正则
var BankCardRegex = regexp.MustCompile(`^(?:[1-9]{1})(?:\d{15}|\d{18})$`)
var ContainBankCardRegex = regexp.MustCompile(`(?:[1-9]{1})(?:\d{15}|\d{18})`)

// 电话号码正则
var TelphoneRegex = regexp.MustCompile(`^0\d{2,3}-\d{7,8}$`)
var ContainTelphoneRegex = regexp.MustCompile(`0\d{2,3}-\d{7,8}`)

// 中国身份证正则
var ChinaIDCardRegex = regexp.MustCompile(`^\d{6}(18|19|20)\d{2}(0\d|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)$`)
var ContainChinaIDCardRegex = regexp.MustCompile(`\d{6}(18|19|20)\d{2}(0\d|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)`)

// 是否为手机
func IsPhone(s string) bool {
	return PhoneRegex.Match([]byte(s))
}

// 是否包含特殊字符,并返回字符数组
func ContainInvalidChars(txt string) (bool, []string) {
	b := InvalidCharsRegexp.Match([]byte(txt))
	if !b {
		return false, nil
	}
	return true, InvalidCharsRegexp.FindAllString(txt, -1)
}

// 验证电子邮件是否合法
func IsEmail(email string) bool {
	return EmailRegex.MatchString(email)
}

var (
	userRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
)

// 验证用户名是否合法
func IsUser(user string) bool {
	return userRegex.MatchString(user)
}
