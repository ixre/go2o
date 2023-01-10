package mss

type MessageScene int
type SmsProvider int

const (
	// 用户登录注册等场景
	UserScene MessageScene = 1
	// 用户密码场景
	PasswordScene MessageScene = 2
	// 其他场景
	OtherScene MessageScene = 3
)

const (
	HTTP SmsProvider = 1
	// 通用HTTP接口
	TECENT_CLOUD SmsProvider = 2
	// 阿里云短信
	ALIYUN SmsProvider = 3
	// 创蓝短信
	CHUANGLAN SmsProvider = 4
)
