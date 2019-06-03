/**
 * Copyright 2015 @ to2.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/valueobject"
	"time"
)

var _ merchant.IConfManager = new(confManagerImpl)

type confManagerImpl struct {
	mchId         int32
	repo          merchant.IMerchantRepo
	saleConf      *merchant.SaleConf
	valRepo       valueobject.IValueRepo
	memberRepo    member.IMemberRepo
	tradeConfList []*merchant.TradeConf
}

func newConfigManagerImpl(mchId int32,
	repo merchant.IMerchantRepo, memberRepo member.IMemberRepo,
	valRepo valueobject.IValueRepo) merchant.IConfManager {
	return &confManagerImpl{
		mchId:      mchId,
		repo:       repo,
		memberRepo: memberRepo,
		valRepo:    valRepo,
	}
}

// 获取销售配置
func (c *confManagerImpl) GetSaleConf() merchant.SaleConf {
	if c.saleConf == nil {
		c.saleConf = c.repo.GetMerchantSaleConf(c.mchId)
		if c.saleConf != nil {
			c.verifySaleConf(c.saleConf)
		} else {
			c.saleConf = &merchant.SaleConf{
				MerchantId: c.mchId,
			}
			c.loadGlobSaleConf(c.saleConf)
		}
	}
	return *c.saleConf
}

func (c *confManagerImpl) loadGlobSaleConf(dst *merchant.SaleConf) error {
	cfg := c.valRepo.GetGlobMchSaleConf()
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
	return c.repo.SaveMerchantSaleConf(c.saleConf)
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
	c.saleConf.MerchantId = c.mchId
	return c.repo.SaveMerchantSaleConf(c.saleConf)
}

// 验证销售设置
func (c *confManagerImpl) verifySaleConf(v *merchant.SaleConf) error {
	cfg := c.valRepo.GetGlobMchSaleConf()
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

func (c *confManagerImpl) getAllMchBuyerGroups() []*merchant.MchBuyerGroup {
	return c.repo.SelectMchBuyerGroup(c.mchId)
}

// 获取商户的全部客户分组
func (c *confManagerImpl) SelectBuyerGroup() []*merchant.BuyerGroup {
	groups := c.memberRepo.GetManager().GetAllBuyerGroups()
	myGroups := c.getAllMchBuyerGroups()
	mp := make(map[int32]*merchant.MchBuyerGroup)
	for _, v := range myGroups {
		mp[v.GroupId] = v
	}
	list := make([]*merchant.BuyerGroup, len(groups))
	for i, v := range groups {
		list[i] = &merchant.BuyerGroup{
			GroupId: v.ID,
			Name:    v.Name,
		}
		mg, ok := mp[v.ID]
		if ok && mg.Alias != "" {
			list[i].Name = mg.Alias
		}
	}
	return list
}

// 保存客户分组
func (c *confManagerImpl) SaveMchBuyerGroup(v *merchant.MchBuyerGroup) (int32, error) {
	g := c.GetGroupByGroupId(v.GroupId)
	g.Alias = v.Alias
	g.EnableRetail = v.EnableRetail
	g.EnableWholesale = v.EnableWholesale
	g.RebatePeriod = v.RebatePeriod
	return util.I32Err(c.repo.SaveMchBuyerGroup(g))
}

// 根据分组编号获取分组设置
func (c *confManagerImpl) GetGroupByGroupId(groupId int32) *merchant.MchBuyerGroup {
	v := c.repo.GetMchBuyerGroupByGroupId(c.mchId, groupId)
	if v != nil {
		return v
	}
	g := c.memberRepo.GetManager().GetBuyerGroup(groupId)
	if g != nil {
		return &merchant.MchBuyerGroup{
			MchId:           c.mchId,
			GroupId:         groupId,
			Alias:           g.Name,
			EnableRetail:    1,
			EnableWholesale: 1,
			RebatePeriod:    1,
		}
	}
	return nil
}

// 获取所有的交易设置
func (c *confManagerImpl) GetAllTradeConf() []*merchant.TradeConf {
	if c.tradeConfList == nil {
		c.tradeConfList = c.repo.SelectMchTradeConf("mch_id= $1", c.mchId)
		if len(c.tradeConfList) == 0 {
			// 零售订单费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:   merchant.TKNormalOrder,
				Flag:        merchant.TFlagNormal,
				AmountBasis: enum.AmountBasisByPercent,
				TradeFee:    0,
				TradeRate:   int(0.2 * enum.RATE_PERCENT),
			})
			// 线下支付费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:   merchant.TKTradeOrder,
				Flag:        merchant.TFlagNormal,
				AmountBasis: enum.AmountBasisByPercent,
				TradeFee:    0,
				TradeRate:   int(0.2 * enum.RATE_PERCENT),
			})
			// 批发订单费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:   merchant.TKWholesaleOrder,
				Flag:        merchant.TFlagNormal,
				AmountBasis: enum.AmountBasisByPercent,
				TradeFee:    0,
				TradeRate:   int(0.1 * enum.RATE_PERCENT),
			})

		}
	}
	return c.tradeConfList
}

// 根据交易类型获取交易设置
func (c *confManagerImpl) GetTradeConf(tradeType int) *merchant.TradeConf {
	for _, v := range c.GetAllTradeConf() {
		if v.TradeType == tradeType {
			return v
		}
	}
	return nil
}

// 保存交易设置
func (c *confManagerImpl) SaveTradeConf(arr []*merchant.TradeConf) error {
	if arr == nil || len(arr) == 0 {
		return errors.New("trade config array is nil")
	}
	unix := time.Now().Unix()
	for _, v := range arr {
		if v.Flag <= 0 {
			v.Flag = merchant.TFlagNormal
		}
		v.UpdateTime = unix
		origin := c.GetTradeConf(v.TradeType)
		if origin != nil {
			origin.TradeFee = v.TradeFee
			origin.MchId = int64(c.mchId)
			origin.AmountBasis = v.AmountBasis
			origin.Flag = v.Flag
			origin.PlanId = v.PlanId
			origin.TradeRate = v.TradeRate
			origin.UpdateTime = v.UpdateTime
			c.repo.SaveMchTradeConf(origin)
		} else {
			c.repo.SaveMchTradeConf(v)
		}
	}
	c.tradeConfList = nil
	return nil
}
