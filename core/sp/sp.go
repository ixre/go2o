package sp

import "github.com/ixre/go2o/core/domain/interface/registry"

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

func (c *ServiceProviderConfiguration) Configure() {

}
