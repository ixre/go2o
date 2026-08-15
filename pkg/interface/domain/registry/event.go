package registry

// 注册中心推送事件
type RegistryPushEvent struct {
	// 是否为用户定义的注册键
	IsUser bool
	// 配置键
	Key string
	// 配置值
	Value string
}
