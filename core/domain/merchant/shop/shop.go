/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 07:55
 * description :
 * history :
 */

package shop

import (
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/lbs"
	"log"
)

var _ shop.IShop = new(ShopImpl)

type ShopImpl struct {
	_manager  *shopManagerImpl
	_shopRep  shop.IShopRep
	_value    *shop.Shop
	_merchant merchant.IMerchant
	_lng      float64
	_lat      float64
}

func newShop(pt merchant.IMerchant, manager *shopManagerImpl,
	v *shop.Shop, shopRep shop.IShopRep) shop.IShop {
	return &ShopImpl{
		_manager:  manager,
		_shopRep:  shopRep,
		_value:    v,
		_merchant: pt,
	}
}

func (this *ShopImpl) GetDomainId() int {
	return this._value.Id
}

func (this *ShopImpl) GetValue() shop.Shop {
	return *this._value
}

func (this *ShopImpl) SetValue(v *shop.Shop) error {
	//	if this.value.Address != v.Address ||
	//		len(this.value.Location) == 0 {
	//		lng, lat, err := lbs.GetLocation(v.Address)
	//		if err != nil {
	//			return err
	//		}
	//		this.value.Location = fmt.Sprintf("%f,%f", lng, lat)
	//}
	//this.value.DeliverRadius = v.DeliverRadius
	this._value.Address = v.Address
	this._value.Name = v.Name
	this._value.SortNumber = v.SortNumber
	this._value.MerchantId = this._merchant.GetAggregateRootId()
	this._value.Phone = v.Phone
	this._value.State = v.State
	return nil
}

func (this *ShopImpl) Save() (int, error) {
	//todo: clear cache
	//this.partner.clearShopCache()
	return this._shopRep.SaveShop(this._value)
}

// 获取经维度
func (this *ShopImpl) GetLngLat() (float64, float64) {
	if this._lng == 0 || this._lat == 0 {
		var err error
		this._lng, this._lat, err = lbs.GetLocation(this._value.Address)
		if err != nil {
			log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		}
	}
	return this._lng, this._lat
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *ShopImpl) CanDeliver(lng, lat float64) (bool, int) {
	//todo:
	return true, -1
	//shopLng, shopLat := this.GetLngLat()
	//distance := lbs.GetLocDistance(shopLng, shopLat, lng, lat)
	//i := int(distance)
	//return i <= this.value.DeliverRadius*1000, i
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *ShopImpl) CanDeliverTo(address string) (bool, int) {
	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		return false, -1
	}
	return this.CanDeliver(lng, lat)
}
