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
	// 自定义短信发送,推送到消息队列由外部系统处理
	CUSTOM SmsProvider = 1
	// HTTP短信
	HTTP SmsProvider = 2
	// 通用HTTP接口
	TECENT_CLOUD SmsProvider = 3
	// 阿里云短信
	ALIYUN SmsProvider = 4
	// 创蓝短信
	CHUANGLAN SmsProvider = 5
)
