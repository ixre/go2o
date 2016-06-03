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
var _ shop.IOfflineShop = new(offlineShopImpl)
var _ shop.IOnlineShop = new(onlineShopImpl)

type ShopImpl struct {
	_manager  *shopManagerImpl
	_shopRep  shop.IShopRep
	_value    *shop.Shop
	_merchant merchant.IMerchant
}

func newShop(pt merchant.IMerchant, manager *shopManagerImpl,
	v *shop.Shop, shopRep shop.IShopRep) shop.IShop {
	s := &ShopImpl{
		_manager:  manager,
		_shopRep:  shopRep,
		_value:    v,
		_merchant: pt,
	}
	switch s.Type() {
	case shop.TypeOnlineShop:
		return newOnlineShopImpl(s)
	case shop.TypeOfflineShop:
		return newOfflineShopImpl(s)
	}
	panic("未知的商店类型")
}

func (this *ShopImpl) GetDomainId() int {
	return this._value.Id
}

// 商店类型
func (this *ShopImpl) Type() int {
	return this._value.ShopType
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
	this._value.Name = v.Name
	this._value.SortNumber = v.SortNumber
	this._value.State = v.State
	return nil
}

func (this *ShopImpl) Save() (int, error) {
	this._merchant.ShopManager().Reload() //清除缓存
	return this._shopRep.SaveShop(this._value)
}

// 数据
func (this *ShopImpl) Data() *shop.ShopDto {
	return &shop.ShopDto{
		Id:         this.GetDomainId(),
		MerchantId: this._value.MerchantId,
		ShopType:   this.Type(),
		Name:       this._value.Name,
		State:      this._value.State,
		Data:       nil,
	}
}

type offlineShopImpl struct {
	*ShopImpl
	//todo: lng,lat要去掉
	_lng     float64
	_lat     float64
	_shopVal *shop.OfflineShop
}

func newOfflineShopImpl(s *ShopImpl) shop.IShop {
	var v *shop.OfflineShop
	if s.GetDomainId() > 0 {
		v = s._shopRep.GetOfflineShop(s.GetDomainId())
	}
	if v == nil {
		dv := shop.DefaultOfflineShop
		v = &dv
		v.ShopId = s.GetDomainId()
	}

	return &offlineShopImpl{
		ShopImpl: s,
		_shopVal: v,
	}
}

// 设置值
func (this *offlineShopImpl) SetShopValue(v *shop.OfflineShop) error {
	this._shopVal.Address = v.Address
	this._shopVal.Tel = v.Tel
	this._shopVal.DeliverRadius = v.DeliverRadius
	if v.Lat > 0 && v.Lng > 0 {
		this._shopVal.Lat = v.Lat
		this._shopVal.Lng = v.Lng
	}
	return nil
}

// 获取值
func (this *offlineShopImpl) GetShopValue() shop.OfflineShop {
	return *this._shopVal
}

// 获取经维度
func (this *offlineShopImpl) GetLngLat() (float64, float64) {
	if this._lng == 0 || this._lat == 0 {
		//todo: 基于位置获取坐标,已经将坐标存储到数据库中了
		var err error
		this._lng, this._lat, err = lbs.GetLocation(this._shopVal.Location())
		if err != nil {
			log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		}
	}
	return this._lng, this._lat
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *offlineShopImpl) CanDeliver(lng, lat float64) (bool, int) {
	shopLng, shopLat := this.GetLngLat()
	distance := lbs.GetLocDistance(shopLng, shopLat, lng, lat)
	i := int(distance)
	return i <= this._shopVal.DeliverRadius*1000, i
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *offlineShopImpl) CanDeliverTo(address string) (bool, int) {
	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		return false, -1
	}
	return this.CanDeliver(lng, lat)
}

// 保存
func (this *offlineShopImpl) Save() (int, error) {
	create := this.GetDomainId() <= 0
	id, err := this.ShopImpl.Save()
	if err == nil {
		this._shopVal.ShopId = id
		err = this._shopRep.SaveOfflineShop(this._shopVal, create)
	}
	return id, err
}

// 数据
func (this *offlineShopImpl) Data() *shop.ShopDto {
	v := this.ShopImpl.Data()
	v.Data = this.GetShopValue()
	return v
}

type onlineShopImpl struct {
	*ShopImpl
	_shopVal *shop.OnlineShop
}

func newOnlineShopImpl(s *ShopImpl) shop.IShop {
	var v *shop.OnlineShop
	if s.GetDomainId() > 0 {
		v = s._shopRep.GetOnlineShop(s.GetDomainId())
	}
	if v == nil {
		dv := shop.DefaultOnlineShop
		v = &dv
		v.ShopId = s.GetDomainId()
	}
	return &onlineShopImpl{
		ShopImpl: s,
		_shopVal: v,
	}
}

// 设置值
func (this *onlineShopImpl) SetShopValue(v *shop.OnlineShop) error {
	this._shopVal.Tel = v.Tel
	this._shopVal.Address = v.Address
	if len(v.Alias) > 0 {
		this._shopVal.Alias = v.Alias
	}
	if len(v.Host) > 0 {
		this._shopVal.Host = v.Host
	}
	if len(v.Logo) > 0 {
		this._shopVal.Logo = v.Logo
	}
	this._shopVal.IndexTitle = v.IndexTitle
	this._shopVal.SubTitle = v.SubTitle
	this._shopVal.Notice = v.Notice
	return nil
}

// 获取值
func (this *onlineShopImpl) GetShopValue() shop.OnlineShop {
	return *this._shopVal
}

// 保存
func (this *onlineShopImpl) Save() (int, error) {
	create := this.GetDomainId() <= 0
	id, err := this.ShopImpl.Save()
	if err == nil {
		this._shopVal.ShopId = id
		err = this._shopRep.SaveOnlineShop(this._shopVal, create)
	}
	return id, err
}

// 数据
func (this *onlineShopImpl) Data() *shop.ShopDto {
	v := this.ShopImpl.Data()
	v.Data = this.GetShopValue()
	return v
}
