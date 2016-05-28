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
	"go2o/core/domain/interface/valueobject"
)

var _ merchant.IConfManager = new(confManagerImpl)

type confManagerImpl struct {
	_rep      merchant.IMerchantRep
	_merchant merchant.IMerchant
	_saleConf *merchant.SaleConf
	_valRep   valueobject.IValueRep
}

func newConfigManagerImpl(m merchant.IMerchant,
	rep merchant.IMerchantRep, valRep valueobject.IValueRep) merchant.IConfManager {
	return &confManagerImpl{
		_merchant: m,
		_rep:      rep,
		_valRep:   valRep,
	}
}

func (this *confManagerImpl) getMerchantId() int {
	return this._merchant.GetAggregateRootId()
}

// 获取销售配置
func (this *confManagerImpl) GetSaleConf() merchant.SaleConf {
	if this._saleConf == nil {
		this._saleConf = this._rep.GetMerchantSaleConf(this.getMerchantId())
		if this._saleConf != nil {
			this.verifySaleConf(this._saleConf)
		} else {
			this._saleConf = &merchant.SaleConf{
				MerchantId: this.getMerchantId(),
			}
			this.loadGlobSaleConf(this._saleConf)
		}
	}
	return *this._saleConf
}

func (this *confManagerImpl) loadGlobSaleConf(dst *merchant.SaleConf) error {
	cfg := this._valRep.GetGlobMerchantSaleConf()
	// 是否启用分销
	if cfg.FxSalesEnabled {
		dst.FxSalesEnabled = 1
	} else {
		dst.FxSalesEnabled = 0
	}
	// 返现比例,0则不返现
	dst.CashBackPercent = cfg.CashBackPercent
	// 一级比例
	dst.CashBackTg1Percent = cfg.CashBackTg1Percent
	// 二级比例
	dst.CashBackTg2Percent = cfg.CashBackTg2Percent
	// 会员比例
	dst.CashBackMemberPercent = cfg.CashBackMemberPercent
	// 自动设置订单
	dst.AutoSetupOrder = cfg.AutoSetupOrder
	// 订单超时分钟数
	dst.OrderTimeOutMinute = cfg.OrderTimeOutMinute
	// 订单自动确认时间
	dst.OrderConfirmAfterMinute = cfg.OrderConfirmAfterMinute
	// 订单超时自动收货
	dst.OrderTimeOutReceiveHour = cfg.OrderTimeOutReceiveHour
	return nil
}

// 使用系统的配置并保存
func (this *confManagerImpl) UseGlobSaleConf() error {
	this.GetSaleConf()
	this.loadGlobSaleConf(this._saleConf)
	return this._rep.SaveMerchantSaleConf(this._saleConf)
}

// 保存销售配置
func (this *confManagerImpl) SaveSaleConf(v *merchant.SaleConf) error {
	if v.CashBackPercent >= 1 || (v.CashBackTg1Percent+
		v.CashBackTg2Percent+v.CashBackMemberPercent) > 1 {
		return merchant.ErrSalesPercent
	}
	this.GetSaleConf()
	this.verifySaleConf(v)
	this._saleConf = v
	this._saleConf.MerchantId = this.getMerchantId()
	return this._rep.SaveMerchantSaleConf(this._saleConf)
}

// 验证销售设置
func (this *confManagerImpl) verifySaleConf(v *merchant.SaleConf) {
	cfg := this._valRep.GetGlobMerchantSaleConf()
	if v.OrderTimeOutMinute <= 0 {
		v.OrderTimeOutMinute = cfg.OrderTimeOutMinute
	}
	if v.OrderConfirmAfterMinute <= 0 {
		v.OrderConfirmAfterMinute = cfg.OrderConfirmAfterMinute
	}
	if v.OrderTimeOutReceiveHour <= 0 {
		v.OrderTimeOutReceiveHour = cfg.OrderTimeOutReceiveHour
	}
	if v.CashBackPercent >= 1 || (v.CashBackTg1Percent+
		v.CashBackTg2Percent+v.CashBackMemberPercent) > 1 {
		v.FxSalesEnabled = 0 //自动关闭分销
	}

}
