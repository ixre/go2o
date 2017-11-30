/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package rsi

import (
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/gen-code/thrift/define"
	"go2o/core/service/thrift/parser"
	"go2o/core/variable"
)

type shopService struct {
	_rep     shop.IShopRepo
	_mchRepo merchant.IMerchantRepo
	_query   *query.ShopQuery
}

func NewShopService(rep shop.IShopRepo, mchRepo merchant.IMerchantRepo,
	query *query.ShopQuery) *shopService {
	return &shopService{
		_rep:     rep,
		_mchRepo: mchRepo,
		_query:   query,
	}
}

// 获取商铺
func (s *shopService) GetStore(vendorId int32) (*define.Store, error) {
	mch := s._mchRepo.GetMerchant(vendorId)
	if mch != nil {
		shop := mch.ShopManager().GetOnlineShop()
		if shop != nil {
			return parser.ParseOnlineShop(shop), nil
		}
	}
	return nil, nil
}

func (s *shopService) GetStoreById(shopId int32) (*define.Store, error) {
	vendorId := s._query.GetMerchantId(shopId)
	return s.GetStore(vendorId)
}

// 打开或关闭商店
func (s *shopService) TurnShop(shopId int32, on bool, reason string) (*define.Result_, error) {
	var err error
	sp := s._rep.GetShop(shopId)
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if on {
			err = sp.TurnOn()
		} else {
			err = sp.TurnOff(reason)
		}
	}
	return parser.Result(shopId, err), nil
}

// 设置商店是否营业
func (s *shopService) OpenShop(shopId int32, on bool, reason string) (*define.Result_, error) {
	var err error
	sp := s._rep.GetShop(shopId)
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if on {
			err = sp.Opening()
		} else {
			err = sp.Pause()
		}
	}
	return parser.Result(shopId, err), nil
}

func (ss *shopService) getMerchantId(shopId int32) int32 {
	return ss._query.GetMerchantId(shopId)
}

func (ss *shopService) GetMerchantId(shopId int32) int32 {
	return ss._query.GetMerchantId(shopId)
}

// 根据主机查询商户编号
func (ss *shopService) GetShopIdByHost(host string) (mchId int32, shopId int32) {
	return ss._query.QueryShopIdByHost(host)
}

// 获取商店的数据
func (ss *shopService) GetShopData(mchId, shopId int32) *shop.ComplexShop {
	mch := ss._mchRepo.GetMerchant(mchId)
	sp := mch.ShopManager().GetShop(shopId)
	if sp != nil {
		return sp.Data()
	}
	return nil
}

func (ss *shopService) GetShopValueById(mchId, shopId int32) *shop.Shop {
	mch := ss._mchRepo.GetMerchant(mchId)
	if mch != nil {
		v := mch.ShopManager().GetShop(shopId).GetValue()
		return &v
	}
	return nil
}

// 保存线上商店
func (ss *shopService) SaveStore(s *define.Store) error {
	mch := ss._mchRepo.GetMerchant(s.VendorId)
	if mch != nil {
		v, v1 := parser.Parse2OnlineShop(s)
		mgr := mch.ShopManager()
		sp := mgr.GetOnlineShop()
		// 创建商店
		if sp == nil {
			sp = mgr.CreateShop(v)
		}
		err := sp.SetValue(v)
		if err == nil {
			ofs := sp.(shop.IOnlineShop)
			err = ofs.SetShopValue(v1)
			if err == nil {
				_, err = sp.Save()
			}
		}
		return err
	}
	return merchant.ErrNoSuchMerchant
}

// 保存门店
func (ss *shopService) SaveOfflineShop(s *shop.Shop, v *shop.OfflineShop) error {
	mch := ss._mchRepo.GetMerchant(s.VendorId)
	if mch != nil {
		mgr := mch.ShopManager()
		var sp shop.IShop
		if s.Id > 0 {
			// 保存商店
			sp = mgr.GetShop(s.Id)
		} else {
			//创建商店
			sp = mgr.CreateShop(s)
		}
		err := sp.SetValue(s)
		if err == nil {
			ofs := sp.(shop.IOfflineShop)
			err = ofs.SetShopValue(v)
			if err == nil {
				_, err = sp.Save()
			}
		}
		return err
	}
	return merchant.ErrNoSuchMerchant
}

func (ss *shopService) SaveShop(mchId int32, v *shop.Shop) (int32, error) {
	mch := ss._mchRepo.GetMerchant(mchId)
	if mch != nil {
		var shop shop.IShop
		if v.Id > 0 {
			shop = mch.ShopManager().GetShop(v.Id)
			if shop == nil {
				return 0, errors.New("门店不存在")
			}
		} else {
			shop = mch.ShopManager().CreateShop(v)
		}
		err := shop.SetValue(v)
		if err != nil {
			return v.Id, err
		}
		return shop.Save()
	}
	return 0, merchant.ErrNoSuchMerchant
}

func (ss *shopService) DeleteShop(mchId, shopId int32) error {
	mch := ss._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ShopManager().DeleteShop(shopId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取线上商城配置
func (ss *shopService) GetOnlineShopConf(shopId int32) *shop.OnlineShop {
	mchId := ss.getMerchantId(shopId)
	mch := ss._mchRepo.GetMerchant(mchId)
	if mch != nil {
		s := mch.ShopManager().GetShop(shopId)
		if s == nil {
			v := s.(shop.IOnlineShop).GetShopValue()
			return &v
		}
	}
	return nil
}

// 获取商城
func (ss *shopService) GetOnlineShops(vendorId int32) []*shop.Shop {
	mch := ss._mchRepo.GetMerchant(vendorId)
	shops := mch.ShopManager().GetShops()
	sv := []*shop.Shop{}
	for _, v := range shops {
		if v.Type() == shop.TypeOnlineShop {
			vv := v.GetValue()
			sv = append(sv, &vv)
		}
	}
	return sv
}

// 获取指定的营业中的店铺
func (ss *shopService) PagedOnBusinessOnlineShops(begin, end int, where, order string) (int, []*dto.ListOnlineShop) {
	n, rows := ss._query.PagedOnBusinessOnlineShops(begin, end, where, order)
	if len(rows) > 0 {
		for _, v := range rows {
			v.Logo = format.GetResUrl(v.Logo)
			if v.Host == "" {
				v.Host = v.Alias + "." + variable.Domain
			}
		}
	}
	return n, rows
}
