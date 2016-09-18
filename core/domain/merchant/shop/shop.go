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
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure/lbs"
	"log"
	"regexp"
	"strings"
)

var _ shop.IShop = new(ShopImpl)
var _ shop.IOfflineShop = new(offlineShopImpl)
var _ shop.IOnlineShop = new(onlineShopImpl)
var (
	shopAliasRegexp = regexp.MustCompile("^[A-Za-z0-9-]{3,11}$")
)

type ShopImpl struct {
	_manager *shopManagerImpl
	_shopRep shop.IShopRep
	_value   *shop.Shop
}

func newShop(manager *shopManagerImpl,
	v *shop.Shop, shopRep shop.IShopRep, valRep valueobject.IValueRep) shop.IShop {
	s := &ShopImpl{
		_manager: manager,
		_shopRep: shopRep,
		_value:   v,
	}
	switch s.Type() {
	case shop.TypeOnlineShop:
		return newOnlineShopImpl(s, valRep)
	case shop.TypeOfflineShop:
		return newOfflineShopImpl(s)
	}
	panic("未知的商店类型")
}

func (s *ShopImpl) GetDomainId() int {
	return s._value.Id
}

// 商店类型
func (s *ShopImpl) Type() int {
	return s._value.ShopType
}

func (s *ShopImpl) GetValue() shop.Shop {
	return *s._value
}

func (s *ShopImpl) SetValue(v *shop.Shop) error {
	//	if s.value.Address != v.Address ||
	//		len(s.value.Location) == 0 {
	//		lng, lat, err := lbs.GetLocation(v.Address)
	//		if err != nil {
	//			return err
	//		}
	//		s.value.Location = fmt.Sprintf("%f,%f", lng, lat)
	//}
	//s.value.DeliverRadius = v.DeliverRadius
	if err := s.check(v); err != nil {
		return err
	}
	s._value.Name = v.Name
	s._value.SortNumber = v.SortNumber
	s._value.State = v.State
	return nil
}

func (s *ShopImpl) check(v *shop.Shop) error {
	v.Name = strings.TrimSpace(v.Name)
	if s.checkNameExists(v) {
		return shop.ErrSameNameShopExists
	}
	return nil
}

func (s *ShopImpl) checkNameExists(v *shop.Shop) bool {
	i := 0
	tmp.Db().ExecScalar("SELECT COUNT(0) FROM mch_shop WHERE name=? AND id <> ?", &i,
		v.Name, v.Id)
	return i > 0
}

func (s *ShopImpl) Save() (int, error) {
	return s._shopRep.SaveShop(s._value)
}

// 数据
func (s *ShopImpl) Data() *shop.ShopDto {
	return &shop.ShopDto{
		Id:         s.GetDomainId(),
		MerchantId: s._value.MerchantId,
		ShopType:   s.Type(),
		Name:       s._value.Name,
		State:      s._value.State,
		CreateTime: s._value.CreateTime,
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
func (s *offlineShopImpl) SetShopValue(v *shop.OfflineShop) error {
	s._shopVal.Address = v.Address
	s._shopVal.Tel = v.Tel
	s._shopVal.DeliverRadius = v.DeliverRadius
	s._shopVal.Province = v.Province
	s._shopVal.City = v.City
	s._shopVal.District = v.District
	if v.Lat > 0 && v.Lng > 0 {
		s._shopVal.Lat = v.Lat
		s._shopVal.Lng = v.Lng
	}
	return nil
}

// 获取值
func (s *offlineShopImpl) GetShopValue() shop.OfflineShop {
	return *s._shopVal
}

// 获取经维度
func (s *offlineShopImpl) GetLngLat() (float64, float64) {
	if s._lng == 0 || s._lat == 0 {
		//todo: 基于位置获取坐标,已经将坐标存储到数据库中了
		var err error
		s._lng, s._lat, err = lbs.GetLocation(s._shopVal.Location())
		if err != nil {
			log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		}
	}
	return s._lng, s._lat
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (s *offlineShopImpl) CanDeliver(lng, lat float64) (bool, int) {
	shopLng, shopLat := s.GetLngLat()
	distance := lbs.GetLocDistance(shopLng, shopLat, lng, lat)
	i := int(distance)
	return i <= s._shopVal.DeliverRadius*1000, i
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (s *offlineShopImpl) CanDeliverTo(address string) (bool, int) {
	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		log.Println("[ Go2o][ LBS][ Error] -", err.Error())
		return false, -1
	}
	return s.CanDeliver(lng, lat)
}

// 保存
func (s *offlineShopImpl) Save() (int, error) {
	create := s.GetDomainId() <= 0
	id, err := s.ShopImpl.Save()
	if err == nil {
		s._shopVal.ShopId = id
		err = s._shopRep.SaveOfflineShop(s._shopVal, create)
	}
	return id, err
}

// 数据
func (s *offlineShopImpl) Data() *shop.ShopDto {
	v := s.ShopImpl.Data()
	v.Data = s.GetShopValue()
	return v
}

type onlineShopImpl struct {
	*ShopImpl
	_shopVal *shop.OnlineShop
	valRep   valueobject.IValueRep
}

func newOnlineShopImpl(s *ShopImpl, valRep valueobject.IValueRep) shop.IShop {
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
		valRep:   valRep,
	}
}

func (s *onlineShopImpl) checkShopAlias(alias string) error {
	alias = strings.ToLower(alias)
	if !shopAliasRegexp.Match([]byte(alias)) {
		return shop.ErrShopAliasFormat
	}
	conf := s.valRep.GetRegistry()
	arr := strings.Split(conf.ShopIncorrectAliasWords, "|")
	for _, v := range arr {
		if strings.Index(alias, v) != -1 {
			return shop.ErrShopAliasIncorrect
		}
	}
	if s._shopRep.ShopAliasExists(alias, s.GetDomainId()) {
		return shop.ErrShopAliasUsed
	}
	return nil
}

// 设置值
func (s *onlineShopImpl) SetShopValue(v *shop.OnlineShop) error {
	s._shopVal.Tel = v.Tel
	s._shopVal.Address = v.Address

	if len(s._shopVal.Alias) == 0 { //未设置域名情况下可更新
		if len(v.Alias) == 0 {
			return shop.ErrNotSetAlias
		}
		if err := s.checkShopAlias(v.Alias); err != nil {
			return err
		}
		s._shopVal.Alias = strings.ToLower(v.Alias)
	}
	if len(v.Host) > 0 {
		s._shopVal.Host = v.Host
	}
	if len(v.Logo) > 0 {
		s._shopVal.Logo = v.Logo
	}

	s._shopVal.IndexTitle = v.IndexTitle
	s._shopVal.SubTitle = v.SubTitle
	s._shopVal.Notice = v.Notice
	return nil
}

// 获取值
func (s *onlineShopImpl) GetShopValue() shop.OnlineShop {
	return *s._shopVal
}

// 保存
func (s *onlineShopImpl) Save() (int, error) {
	create := s.GetDomainId() <= 0
	id, err := s.ShopImpl.Save()
	if err == nil {
		s._shopVal.ShopId = id
		err = s._shopRep.SaveOnlineShop(s._shopVal, create)
	}
	return id, err
}

// 数据
func (s *onlineShopImpl) Data() *shop.ShopDto {
	v := s.ShopImpl.Data()
	v.Data = s.GetShopValue()
	return v
}
