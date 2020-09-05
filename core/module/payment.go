package module

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/storage"
	"go2o/core/module/bank"
	"go2o/core/module/pay"
)

type PaymentModule struct {
	app     gof.App
	storage storage.Interface
	ptArr   []*bank.PaymentPlatform
	gateway *pay.Gateway
}

func (p *PaymentModule) Submit(userId int64, data map[string]string) error {
	return p.gateway.Submit(userId, data)
}

func (p *PaymentModule) CreateToken(userId int64) string {
	return p.gateway.CreatePostToken(userId)
}

func (p *PaymentModule) CheckAndPayment(userId int64, data map[string]string) error {
	return p.gateway.CheckAndPayment(userId, data["trade_no"], data["trade_pwd"])
}

// 模块数据
func (p *PaymentModule) SetApp(app gof.App) {
	p.app = app
	p.storage = app.Storage()
	p.gateway = pay.NewGateway(p.storage)
}

// 初始化模块
func (p *PaymentModule) Init() {

}

// 获取支付平台
func (p *PaymentModule) GetPayPlatform() []*bank.PaymentPlatform {
	if p.ptArr == nil {
		p.ptArr = []*bank.PaymentPlatform{
			bank.Alipay, bank.ChinaPay,
			bank.Tenpay, bank.KuaiBill}
	}
	return p.ptArr
}
