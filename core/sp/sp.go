package sp

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/util/lbs"
	"github.com/ixre/go2o/core/infrastructure/util/sms"
	"github.com/ixre/go2o/core/infrastructure/util/smtp"
	"github.com/ixre/go2o/core/sp/tencent"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/typeconv"
)

// 第三方服务商初始化
type ServiceProviderConfiguration struct {
	registryRepo registry.IRegistryRepo
	storage      storage.Interface
}

// 获取第三方服务商初始化配置
func NewSPConfig(st storage.Interface, registryRepo registry.IRegistryRepo) *ServiceProviderConfiguration {
	return &ServiceProviderConfiguration{
		registryRepo: registryRepo,
		storage:      st,
	}
}

// Configure 函数用于配置服务提供商配置
func (c *ServiceProviderConfiguration) Configure() {
	// 注册腾讯云短信服务商
	sms.RegisterProvider(tencent.NewTencentSms(c.registryRepo))
	// 注册腾讯位置服务
	lbs.Configure(tencent.NewLbsService(c.registryRepo))
	// 配置Smtp邮箱服务器
	smtp.Configure(getDefaultSmtpServer(c.registryRepo))
	tencent.Configure(c.storage, c.registryRepo)
}

// 通过配置获取Smtp服务器配置信息
func getDefaultSmtpServer(repo registry.IRegistryRepo) *smtp.SmtpConfig {
	repo.CreateUserKey("smtp_host", "smtp.exmail.qq.com", "SMTP服务器地址")
	repo.CreateUserKey("smtp_port", "465", "SMTP服务器端口")
	repo.CreateUserKey("smtp_user", "", "SMTP服务器用户名")
	repo.CreateUserKey("smtp_password", "", "SMTP服务器密码")
	repo.CreateUserKey("smtp_default_from", "Go2o", "SMTP默认发件人,your-name")

	host, _ := repo.GetValue("smtp_host")
	port, _ := repo.GetValue("smtp_port")
	user, _ := repo.GetValue("smtp_user")
	password, _ := repo.GetValue("smtp_password")
	from, _ := repo.GetValue("smtp_default_from")

	return &smtp.SmtpConfig{
		Host:     host,
		Port:     typeconv.Int(port),
		User:     user,
		Password: password,
		From:     from,
	}
}
