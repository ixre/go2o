package mss

type MessageScene int

const (
	// 用户登录注册等场景
	UserScene MessageScene = 1
	// 用户密码场景
	PasswordScene MessageScene = 2
	// 其他场景
	OtherScene MessageScene = 3
)
