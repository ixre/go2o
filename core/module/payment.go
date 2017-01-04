package module

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
	"go2o/core/module/bank"
)

type PaymentModule struct {
	app     gof.App
	storage storage.Interface
	ptArr   []*bank.PaymentPlatform
}

// 模块数据
func (m *PaymentModule) SetApp(app gof.App) {
	m.app = app
	m.storage = app.Storage()
}

// 初始化模块
func (m *PaymentModule) Init() {

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
