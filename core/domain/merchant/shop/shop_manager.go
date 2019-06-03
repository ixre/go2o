/**
 * Copyright 2015 @ to2.net.
 * name : shop_manager.go
 * author : jarryliu
 * date : 2016-05-28 12:13
 * description :
 * history :
 */
package shop

import (
	"errors"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/valueobject"
	"strings"
	"time"
)

var _ shop.IShopManager = new(shopManagerImpl)

type shopManagerImpl struct {
	merchant  merchant.IMerchant
	rep       shop.IShopRepo
	valueRepo valueobject.IValueRepo
}

func NewShopManagerImpl(m merchant.IMerchant, rep shop.IShopRepo,
	valueRepo valueobject.IValueRepo) shop.IShopManager {
	return &shopManagerImpl{
		merchant:  m,
		rep:       rep,
		valueRepo: valueRepo,
	}
}

// 新建商店
func (s *shopManagerImpl) CreateShop(v *shop.Shop) shop.IShop {
	v.CreateTime = time.Now().Unix()
	v.VendorId = s.merchant.GetAggregateRootId()
	return NewShop(v, s.rep, s.valueRepo)
}

// 获取所有商店
func (s *shopManagerImpl) GetShops() []shop.IShop {
	shopList := s.rep.GetShopsOfMerchant(s.merchant.GetAggregateRootId())
	shops := make([]shop.IShop, len(shopList))
	for i, v := range shopList {
		v2 := v
		shops[i] = s.CreateShop(&v2)
	}
	return shops
}

// 根据名称获取商店
func (s *shopManagerImpl) GetShopByName(name string) shop.IShop {
	name = strings.TrimSpace(name)
	for _, v := range s.GetShops() {
		if strings.TrimSpace(v.GetValue().Name) == name {
			return v
		}
	}
	return nil
}

// 获取营业中的商店
func (s *shopManagerImpl) GetBusinessInShops() []shop.IShop {
	list := make([]shop.IShop, 0)
	for _, v := range s.GetShops() {
		if v.GetValue().State == enum.ShopBusinessIn {
			list = append(list, v)
		}
	}
	return list
}

// 获取商店
func (s *shopManagerImpl) GetShop(shopId int32) shop.IShop {
	shops := s.GetShops()
	for _, v := range shops {
		if v.GetValue().Id == shopId {
			return v
		}
	}
	return nil
}

// 获取商铺
func (s *shopManagerImpl) GetOnlineShop() shop.IShop {
	for _, v := range s.GetShops() {
		if v.Type() == shop.TypeOnlineShop {
			return v
		}
	}
	return nil
}

// 删除门店
func (s *shopManagerImpl) DeleteShop(shopId int32) error {
	//todo : 检测订单数量
	mchId := s.merchant.GetAggregateRootId()
	sp := s.GetShop(shopId)
	if sp != nil {
		switch sp.Type() {
		case shop.TypeOfflineShop:
			return s.deleteOfflineShop(mchId, sp)
		case shop.TypeOnlineShop:
			return s.deleteOnlineShop(mchId, sp)
		}
	}
	return nil
}

func (s *shopManagerImpl) deleteOnlineShop(mchId int32, sp shop.IShop) error {
	return errors.New("暂不支持删除线上商店")
	shopId := sp.GetDomainId()
	err := s.rep.DeleteOnlineShop(mchId, shopId)
	return err
}

func (s *shopManagerImpl) deleteOfflineShop(mchId int32, sp shop.IShop) error {
	shopId := sp.GetDomainId()
	err := s.rep.DeleteOfflineShop(mchId, shopId)
	return err
}
