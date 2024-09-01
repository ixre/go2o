/**
 * Copyright 2015 @ 56x.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/gof/util"
)

var _ merchant.IConfManager = new(confManagerImpl)

type confManagerImpl struct {
	mchId         int
	repo          merchant.IMerchantRepo
	saleConf      *merchant.SaleConf
	valRepo       valueobject.IValueRepo
	memberRepo    member.IMemberRepo
	_settleConf   *merchant.SettleConf
	tradeConfList []*merchant.TradeConf
}

// GetSettleConf implements merchant.IConfManager.
func (c *confManagerImpl) GetSettleConf() *merchant.SettleConf {
	if c._settleConf == nil {
		c._settleConf = c.repo.SettleRepo().FindBy("mch_id=?", c.mchId)
		if c._settleConf == nil {
			c._settleConf = &merchant.SettleConf{
				Id:          0,
				MchId:       c.mchId,
				OrderTxRate: 0,
				OtherTxRate: 0,
				SubMchNo:    "",
				UpdateTime:  int(time.Now().Unix()),
			}
			c.repo.SettleRepo().Save(c._settleConf)
		}
	}
	return c._settleConf
}

// SaveSettleConf implements merchant.IConfManager.
func (c *confManagerImpl) SaveSettleConf(s *merchant.SettleConf) error {
	if s.MchId != c.mchId {
		return errors.New("商户编号不匹配")
	}
	if s.OtherTxRate > 1 {
		return merchant.ErrTxRate
	}
	if s.OrderTxRate > 1 {
		return merchant.ErrTxRate
	}
	o := c.GetSettleConf()
	o.OrderTxRate = s.OrderTxRate
	o.OtherTxRate = s.OtherTxRate
	o.SubMchNo = s.SubMchNo
	o.UpdateTime = int(time.Now().Unix())
	_, err := c.repo.SettleRepo().Save(o)
	return err
}

func newConfigManagerImpl(mchId int,
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
		c.saleConf = c.repo.GetMerchantSaleConf(int64(c.mchId))
		if c.saleConf != nil {
			c.verifySaleConf(c.saleConf)
		} else {
			c.saleConf = &merchant.SaleConf{
				MchId: c.mchId,
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
		dst.FxSales = 1
	} else {
		dst.FxSales = 0
	}
	// 返现比例,0则不返现
	dst.CbPercent = cfg.CashBackPercent
	// 一级比例
	dst.CbTg1Percent = cfg.CashBackTg1Percent
	// 二级比例
	dst.CbTg2Percent = cfg.CashBackTg2Percent
	// 会员比例
	dst.CbMemberPercent = cfg.CashBackMemberPercent
	// 自动设置订单
	dst.OaOpen = cfg.AutoSetupOrder
	// 订单超时分钟数
	dst.OaTimeoutMinute = cfg.OrderTimeOutMinute
	// 订单自动确认时间
	dst.OaConfirmMinute = cfg.OrderConfirmAfterMinute
	// 订单超时自动收货
	dst.OaReceiveHour = cfg.OrderTimeOutReceiveHour
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
	if v.CbPercent >= 1 || (v.CbTg1Percent+
		v.CbTg2Percent+v.CbMemberPercent) > 1 {
		return merchant.ErrSalesPercent
	}
	c.GetSaleConf()
	if err := c.verifySaleConf(v); err != nil {
		return err
	}
	c.saleConf = v
	c.saleConf.MchId = c.mchId
	return c.repo.SaveMerchantSaleConf(c.saleConf)
}

// 验证销售设置
func (c *confManagerImpl) verifySaleConf(v *merchant.SaleConf) error {
	cfg := c.valRepo.GetGlobMchSaleConf()
	if !cfg.FxSalesEnabled && v.FxSales == 1 {
		return merchant.ErrEnabledFxSales
	}
	if v.OaTimeoutMinute <= 0 {
		v.OaTimeoutMinute = cfg.OrderTimeOutMinute
	}
	if v.OaConfirmMinute <= 0 {
		v.OaConfirmMinute = cfg.OrderConfirmAfterMinute
	}
	if v.OaReceiveHour <= 0 {
		v.OaReceiveHour = cfg.OrderTimeOutReceiveHour
	}
	if v.CbPercent >= 1 || (v.CbTg1Percent+
		v.CbTg2Percent+v.CbMemberPercent) > 1 {
		v.FxSales = 0 //自动关闭分销
	}
	return nil
}

func (c *confManagerImpl) getAllMchBuyerGroups() []*merchant.MchBuyerGroupSetting {
	return c.repo.SelectMchBuyerGroup(int64(c.mchId))
}

// 获取商户的全部客户分组
func (c *confManagerImpl) SelectBuyerGroup() []*merchant.BuyerGroup {
	groups := c.memberRepo.GetManager().GetAllBuyerGroups()
	myGroups := c.getAllMchBuyerGroups()
	mp := make(map[int64]*merchant.MchBuyerGroupSetting)
	for _, v := range myGroups {
		mp[v.GroupId] = v
	}
	list := make([]*merchant.BuyerGroup, len(groups))
	for i, v := range groups {
		e := &merchant.BuyerGroup{
			GroupId: int64(v.ID),
			Name:    v.Name,
		}
		mg, ok := mp[int64(v.ID)]
		if ok {
			if mg.Alias != "" {
				e.Name = mg.Alias
			}
			e.EnableWholesale = mg.EnableWholesale == 1
			e.EnableRetail = mg.EnableRetail == 1
			e.RebatePeriod = int(mg.RebatePeriod)
		}
		list[i] = e
	}
	return list
}

// 保存客户分组
func (c *confManagerImpl) SaveMchBuyerGroup(v *merchant.MchBuyerGroupSetting) (int32, error) {
	g := c.GetGroupByGroupId(int32(v.GroupId))
	g.Alias = v.Alias
	g.EnableRetail = v.EnableRetail
	g.EnableWholesale = v.EnableWholesale
	g.RebatePeriod = v.RebatePeriod
	return util.I32Err(c.repo.SaveMchBuyerGroup(g))
}

// 根据分组编号获取分组设置
func (c *confManagerImpl) GetGroupByGroupId(groupId int32) *merchant.MchBuyerGroupSetting {
	v := c.repo.GetMchBuyerGroupByGroupId(int32(c.mchId), groupId)
	if v != nil {
		return v
	}
	g := c.memberRepo.GetManager().GetBuyerGroup(groupId)
	if g != nil {
		return &merchant.MchBuyerGroupSetting{
			MerchantId:      int64(c.mchId),
			GroupId:         int64(groupId),
			Alias:           g.Name,
			EnableRetail:    1,
			EnableWholesale: 1,
			RebatePeriod:    1,
		}
	}
	return nil
}

// 获取所有的交易设置
func (c *confManagerImpl) GetAllTradeConf_() []*merchant.TradeConf {
	if c.tradeConfList == nil {
		c.tradeConfList = c.repo.SelectMchTradeConf("mch_id= $1", c.mchId)
		if len(c.tradeConfList) == 0 {
			// 零售订单费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:      merchant.TKNormalOrder,
				Flag:           merchant.TFlagNormal,
				AmountBasis:    enum.AmountBasisByPercent,
				TransactionFee: 0,
				TradeRate:      int(0.2 * enum.RATE_PERCENT),
			})
			// 线下支付费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:      merchant.TKTradeOrder,
				Flag:           merchant.TFlagNormal,
				AmountBasis:    enum.AmountBasisByPercent,
				TransactionFee: 0,
				TradeRate:      int(0.2 * enum.RATE_PERCENT),
			})
			// 批发订单费率
			c.tradeConfList = append(c.tradeConfList, &merchant.TradeConf{
				TradeType:      merchant.TKWholesaleOrder,
				Flag:           merchant.TFlagNormal,
				AmountBasis:    enum.AmountBasisByPercent,
				TransactionFee: 0,
				TradeRate:      int(0.1 * enum.RATE_PERCENT),
			})

		}
	}
	return c.tradeConfList
}

// 根据交易类型获取交易设置
func (c *confManagerImpl) GetTradeConf(tradeType int) *merchant.TradeConf {
	for _, v := range c.GetAllTradeConf_() {
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
			origin.TransactionFee = v.TransactionFee
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
