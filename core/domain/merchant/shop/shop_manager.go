/**
 * Copyright 2015 @ z3q.net.
 * name : shop_manager.go
 * author : jarryliu
 * date : 2016-05-28 12:13
 * description :
 * history :
 */
package shop

import (
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"time"
)

var _ shop.IShopManager = new(shopManagerImpl)

type shopManagerImpl struct {
	_merchant merchant.IMerchant
	_rep      shop.IShopRep
	_shops    []shop.IShop
}

func NewShopManagerImpl(m merchant.IMerchant, rep shop.IShopRep) shop.IShopManager {
	return &shopManagerImpl{
		_merchant: m,
		_rep:      rep,
	}
}

// 获取站点配置
func (this *shopManagerImpl) GetSiteConf() shop.ShopSiteConf {
	return shop.ShopSiteConf{}
	//if this._siteConf == nil {
	//    this._siteConf = this._rep.GetSiteConf(this.GetAggregateRootId())
	//}
	//return *this._siteConf
}

// 保存站点配置
func (this *shopManagerImpl) SaveSiteConf(v *shop.ShopSiteConf) error {
	return nil
	//    this._siteConf = v
	//    this._siteConf.MerchantId = this._value.Id
	//    return this._rep.SaveSiteConf(this.GetAggregateRootId(), this._siteConf)
	//
}

// 新建商店
func (this *shopManagerImpl) CreateShop(v *shop.Shop) shop.IShop {
	v.CreateTime = time.Now().Unix()
	v.MerchantId = this._merchant.GetAggregateRootId()
	return newShop(this._merchant, this, v, this._rep)
}

// 获取所有商店
func (this *shopManagerImpl) GetShops() []shop.IShop {
	if this._shops == nil {
		shops := this._rep.GetShopsOfMerchant(this._merchant.GetAggregateRootId())
		this._shops = make([]shop.IShop, len(shops))
		for i, v := range shops {
			this._shops[i] = this.CreateShop(v)
		}
	}
	return this._shops
}

// 获取营业中的商店
func (this *shopManagerImpl) GetBusinessInShops() []shop.IShop {
	var list []shop.IShop = make([]shop.IShop, 0)
	for _, v := range this._shops {
		if v.GetValue().State == enum.ShopBusinessIn {
			list = append(list, v)
		}
	}
	return list
}

// 获取商店
func (this *shopManagerImpl) GetShop(shopId int) shop.IShop {
	//	v := this.rep.GetValueShop(this.GetAggregateRootId(), shopId)
	//	if v == nil {
	//		return nil
	//	}
	//	return this.CreateShop(v)
	shops := this.GetShops()

	for _, v := range shops {
		if v.GetValue().Id == shopId {
			return v
		}
	}
	return nil
}

// 删除门店
func (this *shopManagerImpl) DeleteShop(shopId int) error {
	//todo : 检测订单数量
	return this._rep.DeleteShop(this._merchant.GetAggregateRootId(), shopId)
}
