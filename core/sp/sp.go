package sp

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/lbs"
	"github.com/ixre/go2o/core/infrastructure/util/sms"
	"github.com/ixre/go2o/core/sp/tencent"
)

// 第三方服务商初始化
type ServiceProviderConfiguration struct {
	registryRepo registry.IRegistryRepo
}

// 获取第三方服务商初始化配置
func NewSPConfig(registryRepo registry.IRegistryRepo) *ServiceProviderConfiguration {
	return &ServiceProviderConfiguration{
		registryRepo: registryRepo,
	}
}

// Configure 函数用于配置服务提供商配置
func (c *ServiceProviderConfiguration) Configure() {
	// 注册腾讯云短信服务商
	sms.RegisterProvider(tencent.NewTencentSms(c.registryRepo))
	// 注册腾讯位置服务
	lbs.Configure(tencent.NewLbsService(c.registryRepo))
}
