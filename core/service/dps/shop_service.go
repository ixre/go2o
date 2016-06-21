/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps

import (
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/query"
	"log"
)

type shopService struct {
	_rep    shop.IShopRep
	_mchRep merchant.IMerchantRep
	_query  *query.ShopQuery
}

func NewShopService(rep shop.IShopRep, mchRep merchant.IMerchantRep,
	query *query.ShopQuery) *shopService {
	return &shopService{
		_rep:    rep,
		_mchRep: mchRep,
		_query:  query,
	}
}

func (this *shopService) getMerchantId(shopId int) int {
	return this._query.GetMerchantId(shopId)
}

func (this *shopService) GetMerchantId(shopId int) int {
	return this._query.GetMerchantId(shopId)
}

// 根据主机查询商户编号
func (this *shopService) GetShopIdByHost(host string) (mchId int, shopId int) {
	return this._query.QueryShopIdByHost(host)
}

// 获取商店的数据
func (this *shopService) GetShopData(mchId, shopId int) *shop.ShopDto {
	mch, _ := this._mchRep.GetMerchant(mchId)
	sp := mch.ShopManager().GetShop(shopId)
	if sp != nil {
		return sp.Data()
	}
	return nil
}

func (this *shopService) GetShopValueById(mchId, shopId int) *shop.Shop {
	mch, err := this._mchRep.GetMerchant(mchId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	v := mch.ShopManager().GetShop(shopId).GetValue()
	return &v
}

// 保存线上商店
func (this *shopService) SaveOnlineShop(s *shop.Shop, v *shop.OnlineShop) error {
	mch, err := this._mchRep.GetMerchant(s.MerchantId)
	if err == nil {
		mgr := mch.ShopManager()
		var sp shop.IShop
		if s.Id > 0 { // 保存商店
			sp = mgr.GetShop(s.Id)
			err = sp.SetValue(s)
			if err != nil {
				return err
			}
		} else {
			//检测店名是否重复
			if err = this.checkShopName(mgr, s.Id, s.Name); err != nil {
				return err
			}
			// 创建商店
			sp = mgr.CreateShop(s)
		}

		ofs := sp.(shop.IOnlineShop)
		err = ofs.SetShopValue(v)
		if err == nil {
			_, err = sp.Save()
		}
	}
	return err
}

func (this *shopService) checkShopName(mgr shop.IShopManager, id int, name string) error {
	v := mgr.GetShopByName(name)
	if v != nil {
		sv := v.GetValue()
		if sv.Name == sv.Name && sv.Id != id {
			return shop.ErrSameNameShopExists
		}
	}
	return nil
}

// 保存门店
func (this *shopService) SaveOfflineShop(s *shop.Shop, v *shop.OfflineShop) error {
	mch, err := this._mchRep.GetMerchant(s.MerchantId)
	if err == nil {
		mgr := mch.ShopManager()
		var sp shop.IShop
		if s.Id > 0 { // 保存商店
			sp = mgr.GetShop(s.Id)
			err = sp.SetValue(s)
			if err != nil {
				return err
			}
		} else {
			//检测店名是否重复
			if err = this.checkShopName(mgr, s.Id, s.Name); err != nil {
				return err
			}
			//创建商店
			sp = mgr.CreateShop(s)
		}

		ofs := sp.(shop.IOfflineShop)
		err = ofs.SetShopValue(v)
		if err == nil {
			_, err = sp.Save()
		}
	}
	return err
}

func (this *shopService) SaveShop(mchId int, v *shop.Shop) (int, error) {
	mch, err := this._mchRep.GetMerchant(mchId)
	if err != nil {

		log.Println("[ Merchant][ Service]-", err.Error())
		return 0, err
	}
	var shop shop.IShop
	if v.Id > 0 {
		shop = mch.ShopManager().GetShop(v.Id)
		if shop == nil {
			return 0, errors.New("门店不存在")
		}
	} else {
		shop = mch.ShopManager().CreateShop(v)
	}
	err = shop.SetValue(v)
	if err != nil {
		return v.Id, err
	}
	return shop.Save()
}

func (this *shopService) DeleteShop(merchantId, shopId int) error {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Merchant][ Service]-", err.Error())
	}
	return mch.ShopManager().DeleteShop(shopId)
}

// 获取线上商城配置
func (this *shopService) GetOnlineShopConf(shopId int) *shop.OnlineShop {
	merchantId := this.getMerchantId(shopId)
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err == nil {
		s := mch.ShopManager().GetShop(shopId)
		if s == nil {
			v := s.(shop.IOnlineShop).GetShopValue()
			return &v
		}
	}
	return nil
}
