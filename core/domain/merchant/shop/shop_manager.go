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
	"go2o/core/domain/interface/domain/enum"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"strings"
	"time"
)

var _ shop.IShopManager = new(shopManagerImpl)

type shopManagerImpl struct {
	merchant     merchant.IMerchant
	repo         shop.IShopRepo
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
}

// 创建线上店铺
func (s *shopManagerImpl) CreateOnlineShop(o *shop.OnlineShop) (shop.IShop, error) {
	o.VendorId = s.merchant.GetAggregateRootId()
	if o.ShopName == "" {
		o.ShopName = s.merchant.GetValue().Name
	}
	is := s.repo.CreateShop(o)
	io := is.(shop.IOnlineShop)
	err := io.SetShopValue(o)
	if err == nil {
		err = is.Save()
	}
	return is, err
}

func NewShopManagerImpl(m merchant.IMerchant, rep shop.IShopRepo,
	valueRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) shop.IShopManager {
	return &shopManagerImpl{
		merchant:     m,
		repo:         rep,
		valueRepo:    valueRepo,
		registryRepo: registryRepo,
	}
}

// 新建商店
func (s *shopManagerImpl) CreateShop(v *shop.Shop) shop.IShop {
	v.CreateTime = time.Now().Unix()
	v.VendorId = int64(s.merchant.GetAggregateRootId())
	return NewShop2(v, s.repo, s.valueRepo, s.registryRepo)
}

// 获取所有商店
func (s *shopManagerImpl) GetShops() []shop.IShop {
	shopList := s.repo.GetShopsOfMerchant(s.merchant.GetAggregateRootId())
	shops := make([]shop.IShop, len(shopList))
	for i, v := range shopList {
		shops[i] = s.CreateShop(&v)
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
func (s *shopManagerImpl) GetShop(shopId int) shop.IShop {
	shops := s.GetShops()
	for _, v := range shops {
		if int(v.GetValue().Id) == shopId {
			return v
		}
	}
	return nil
}

// 获取店铺
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
	sp := s.GetShop(int(shopId))
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

func (s *shopManagerImpl) deleteOnlineShop(mchId int64, sp shop.IShop) error {
	return errors.New("暂不支持删除线上商店")
	shopId := sp.GetDomainId()
	err := s.repo.DeleteOnlineShop(mchId, int64(shopId))
	return err
}

func (s *shopManagerImpl) deleteOfflineShop(mchId int64, sp shop.IShop) error {
	shopId := sp.GetDomainId()
	err := s.repo.DeleteOfflineShop(mchId,int64(shopId))
	return err
}
