/**
 * Copyright 2015 @ z3q.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"go2o/core/domain/interface/merchant"
)

var _ merchant.IConfManager = new(confManagerImpl)

type confManagerImpl struct {
	_rep      merchant.IMerchantRep
	_merchant merchant.IMerchant

	_saleConf *merchant.SaleConf
}

func newConfigManagerImpl(m merchant.IMerchant,
	rep merchant.IMerchantRep) merchant.IConfManager {
	return &confManagerImpl{
		_merchant: m,
		_rep:      rep,
	}
}

func (this *confManagerImpl) getMerchantId() int {
	return this._merchant.GetAggregateRootId()
}

// 获取销售配置
func (this *confManagerImpl) GetSaleConf() merchant.SaleConf {
	if this._saleConf == nil {
		this._saleConf = this._rep.GetSaleConf(this.getMerchantId())
		this.verifySaleConf(this._saleConf)
	}
	return *this._saleConf
}

// 保存销售配置
func (this *confManagerImpl) SaveSaleConf(v *merchant.SaleConf) error {
	this.GetSaleConf()
	if v.FlowConvertCsn < 0 || v.PresentConvertCsn < 0 ||
		v.ApplyCsn < 0 || v.TransCsn < 0 ||
		v.FlowConvertCsn > 1 || v.PresentConvertCsn > 1 ||
		v.ApplyCsn > 1 || v.TransCsn > 1 {
		return merchant.ErrSalesPercent
	}

	this.verifySaleConf(v)

	this._saleConf = v
	this._saleConf.MerchantId = this.getMerchantId()

	return this._rep.SaveSaleConf(this.getMerchantId(), this._saleConf)
}

// 验证销售设置
func (this *confManagerImpl) verifySaleConf(v *merchant.SaleConf) {
	if v.OrderTimeOutMinute <= 0 {
		v.OrderTimeOutMinute = 1440 // 一天
	}

	if v.OrderConfirmAfterMinute <= 0 {
		v.OrderConfirmAfterMinute = 60 // 一小时后自动确认
	}

	if v.OrderTimeOutReceiveHour <= 0 {
		v.OrderTimeOutReceiveHour = 7 * 24 // 7天后自动确认
	}
}
