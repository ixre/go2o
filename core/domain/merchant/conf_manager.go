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
	rep      merchant.IMerchantRep
	merchant merchant.IMerchant
	saleConf *merchant.SaleConf
	valRep   valueobject.IValueRep
}

func newConfigManagerImpl(m merchant.IMerchant,
	rep merchant.IMerchantRep, valRep valueobject.IValueRep) merchant.IConfManager {
	return &confManagerImpl{
		merchant: m,
		rep:      rep,
		valRep:   valRep,
	}
}

func (c *confManagerImpl) getMerchantId() int {
	return c.merchant.GetAggregateRootId()
}

// 获取销售配置
func (c *confManagerImpl) GetSaleConf() merchant.SaleConf {
	if c.saleConf == nil {
		c.saleConf = c.rep.GetMerchantSaleConf(c.getMerchantId())
		if c.saleConf != nil {
			c.verifySaleConf(c.saleConf)
		} else {
			c.saleConf = &merchant.SaleConf{
				MerchantId: c.getMerchantId(),
			}
			c.loadGlobSaleConf(c.saleConf)
		}
	}
	return *c.saleConf
}

func (c *confManagerImpl) loadGlobSaleConf(dst *merchant.SaleConf) error {
	cfg := c.valRep.GetGlobMchSaleConf()
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
func (c *confManagerImpl) UseGlobSaleConf() error {
	c.GetSaleConf()
	c.loadGlobSaleConf(c.saleConf)
	return c.rep.SaveMerchantSaleConf(c.saleConf)
}

// 保存销售配置
func (c *confManagerImpl) SaveSaleConf(v *merchant.SaleConf) error {
	if v.CashBackPercent >= 1 || (v.CashBackTg1Percent+
		v.CashBackTg2Percent+v.CashBackMemberPercent) > 1 {
		return merchant.ErrSalesPercent
	}
	c.GetSaleConf()
	if err := c.verifySaleConf(v); err != nil {
		return err
	}
	c.saleConf = v
	c.saleConf.MerchantId = c.getMerchantId()
	return c.rep.SaveMerchantSaleConf(c.saleConf)
}

// 验证销售设置
func (c *confManagerImpl) verifySaleConf(v *merchant.SaleConf) error {
	cfg := c.valRep.GetGlobMchSaleConf()
	if !cfg.FxSalesEnabled && v.FxSalesEnabled == 1 {
		return merchant.ErrEnabledFxSales
	}
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
	return nil
}
