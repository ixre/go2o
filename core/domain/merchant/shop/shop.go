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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure/lbs"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var _ shop.IOfflineShop = new(offlineShopImpl)
var _ shop.IOnlineShop = new(onlineShopImpl)
var (
	shopAliasRegexp = regexp.MustCompile("^[A-Za-z0-9-]{3,11}$")
)

type shopImpl struct {
	shopRepo  shop.IShopRepo
	value     *shop.Shop
	valueRepo valueobject.IValueRepo
}

func NewShop(v *shop.Shop, shopRepo shop.IShopRepo,
	valRepo valueobject.IValueRepo) shop.IShop {
	s := &shopImpl{
		shopRepo:  shopRepo,
		value:     v,
		valueRepo: valRepo,
	}
	switch s.Type() {
	case shop.TypeOnlineShop:
		return newOnlineShopImpl(s, valRepo)
	case shop.TypeOfflineShop:
		return newOfflineShopImpl(s, valRepo)
	}
	panic("未知的商店类型")
}

func (s *shopImpl) GetDomainId() int32 {
	return s.value.Id
}

// 商店类型
func (s *shopImpl) Type() int32 {
	return s.value.ShopType
}

func (s *shopImpl) GetValue() shop.Shop {
	return *s.value
}

func (s *shopImpl) SetValue(v *shop.Shop) error {
	if err := s.check(v); err != nil {
		return err
	}
	s.value.Name = v.Name
	s.value.SortNum = v.SortNum
	if s.GetDomainId() <= 0 {
		s.value.State = shop.StateNormal
		s.value.OpeningState = shop.OStateNormal
	}
	return nil
}

func (s *shopImpl) check(v *shop.Shop) error {
	v.Name = strings.TrimSpace(v.Name)
	if s.checkNameExists(v) {
		return shop.ErrSameNameShopExists
	}
	return nil
}

func (s *shopImpl) checkNameExists(v *shop.Shop) bool {
	i := 0
	tmp.Db().ExecScalar("SELECT COUNT(0) FROM mch_shop WHERE name= $1 AND id <> $2", &i,
		v.Name, v.Id)
	return i > 0
}

func (s *shopImpl) Save() (int32, error) {
	return s.shopRepo.SaveShop(s.value)
}

// 开启店铺
func (s *shopImpl) TurnOn() error {
	s.value.State = shop.StateNormal
	_, err := s.Save()
	return err
}

// 停用店铺
func (s *shopImpl) TurnOff(reason string) error {
	s.value.State = shop.StateStopped
	_, err := s.Save()
	return err
}

// 商店营业
func (s *shopImpl) Opening() error {
	s.value.OpeningState = shop.OStateNormal
	_, err := s.Save()
	return err
}

// 商店暂停营业
func (s *shopImpl) Pause() error {
	s.value.OpeningState = shop.OStatePause
	_, err := s.Save()
	return err
}

type offlineShopImpl struct {
	*shopImpl
	//todo: lng,lat要去掉
	_lng     float64
	_lat     float64
	_shopVal *shop.OfflineShop
}

func newOfflineShopImpl(s *shopImpl, valueRepo valueobject.IValueRepo) shop.IShop {
	var v *shop.OfflineShop
	if s.GetDomainId() > 0 {
		v = s.shopRepo.GetOfflineShop(s.GetDomainId())
	}
	if v == nil {
		dv := shop.DefaultOfflineShop
		v = &dv
		v.ShopId = s.GetDomainId()
	}

	return &offlineShopImpl{
		shopImpl: s,
		_shopVal: v,
	}
}

// 设置值
func (s *offlineShopImpl) SetShopValue(v *shop.OfflineShop) error {
	//	if s.value.Address != v.Address ||
	//		len(s.value.Location) == 0 {
	//		lng, lat, err := lbs.GetLocation(v.Address)
	//		if err != nil {
	//			return err
	//		}
	//		s.value.Location = fmt.Sprintf("%f,%f", lng, lat)
	//}
	//s.value.CoverRadius = v.CoverRadius

	s._shopVal.Address = v.Address
	s._shopVal.Tel = v.Tel
	s._shopVal.CoverRadius = v.CoverRadius
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
	return i <= s._shopVal.CoverRadius*1000, i
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
func (s *offlineShopImpl) Save() (int32, error) {
	create := s.GetDomainId() <= 0
	id, err := s.shopImpl.Save()
	if err == nil {
		s._shopVal.ShopId = id
		err = s.shopRepo.SaveOfflineShop(s._shopVal, create)
	}
	return id, err
}

// 获取商店信息
func (s *offlineShopImpl) Data() *shop.ComplexShop {
	sv := s.value
	ov := s._shopVal
	v := &shop.ComplexShop{
		ID:       s.GetDomainId(),
		VendorId: sv.VendorId,
		ShopType: sv.ShopType,
		Name:     sv.Name,
		State:    sv.ShopType,
		Data:     make(map[string]string),
	}
	address := s.valueRepo.AreaString(ov.Province, ov.City, ov.District, ov.Address)
	v.Data["Address"] = address
	v.Data["Tel"] = ov.Tel
	v.Data["CoverRadius"] = strconv.Itoa(int(ov.CoverRadius))
	v.Data["Lat"] = strconv.FormatFloat(float64(ov.Lat), 'g', 2, 32)
	v.Data["Lng"] = strconv.FormatFloat(float64(ov.Lng), 'g', 2, 32)
	return v
}

type onlineShopImpl struct {
	*shopImpl
	_shopVal *shop.OnlineShop
	valRepo  valueobject.IValueRepo
}

func newOnlineShopImpl(s *shopImpl, valRepo valueobject.IValueRepo) shop.IShop {
	var v *shop.OnlineShop
	if s.GetDomainId() > 0 {
		v = s.shopRepo.GetOnlineShop(s.GetDomainId())
	}
	if v == nil {
		dv := shop.DefaultOnlineShop
		v = &dv
		v.ShopId = s.GetDomainId()
	}
	return &onlineShopImpl{
		shopImpl: s,
		_shopVal: v,
		valRepo:  valRepo,
	}
}

func (s *onlineShopImpl) checkShopAlias(alias string) error {
	alias = strings.ToLower(alias)
	if !shopAliasRegexp.Match([]byte(alias)) {
		return shop.ErrShopAliasFormat
	}

	//todo: 非法关键字
	//arr := strings.Split(conf.ShopIncorrectAliasWords, "|")
	arr := strings.Split("", "|")
	for _, v := range arr {
		if strings.Index(alias, v) != -1 {
			return shop.ErrShopAliasIncorrect
		}
	}
	if s.shopRepo.ShopAliasExists(alias, s.GetDomainId()) {
		return shop.ErrShopAliasUsed
	}
	return nil
}

// 设置值
func (s *onlineShopImpl) SetShopValue(v *shop.OnlineShop) (err error) {
	v.Logo = strings.TrimSpace(v.Logo)
	s._shopVal.ServiceTel = v.ServiceTel
	s._shopVal.Address = v.Address
	if len(v.Host) > 0 {
		s._shopVal.Host = v.Host
	}
	if s.GetDomainId() == 0 {
		if v.Logo == "" {
			return shop.ErrShopNoLogo
		}
	}
	if len(v.Logo) > 0 {
		s._shopVal.Logo = v.Logo
	}
	if len(v.Alias) > 0 {
		err = s.checkShopAlias(v.Alias)
		if err == nil {
			s._shopVal.Alias = v.Alias
		}
	}
	s._shopVal.ShopTitle = v.ShopTitle
	s._shopVal.ShopNotice = v.ShopNotice
	return err
}

// 获取值
func (s *onlineShopImpl) GetShopValue() shop.OnlineShop {
	return *s._shopVal
}

// 保存
func (s *onlineShopImpl) Save() (int32, error) {
	create := s.GetDomainId() <= 0
	if create {
		if s.shopRepo.ShopCount(s.value.VendorId, shop.TypeOnlineShop) > 0 {
			return 0, shop.ErrSupportSingleOnlineShop
		}
		s._shopVal.Alias = s.generateShopAlias()
	}
	id, err := s.shopImpl.Save()
	if err == nil {
		s._shopVal.ShopId = id
		err = s.shopRepo.SaveOnlineShop(s._shopVal, create)
	}
	return id, err
}

func (s *onlineShopImpl) generateShopAlias() string {
	return "shop" + strconv.Itoa(util.RandInt(8))
	//todo: ???
	for {
		id := "shop" + strconv.Itoa(util.RandInt(8))
		if err := s.checkShopAlias(id); err == nil {
			return id
		}
	}
	return ""
}

// 获取商店信息
func (s *onlineShopImpl) Data() *shop.ComplexShop {
	sv := s.value
	ov := s._shopVal
	v := &shop.ComplexShop{
		ID:       s.GetDomainId(),
		VendorId: sv.VendorId,
		ShopType: sv.ShopType,
		Name:     sv.Name,
		State:    sv.ShopType,
		Data:     make(map[string]string),
	}
	v.Data["Host"] = ov.Host
	v.Data["Logo"] = ov.Logo
	v.Data["ServiceTel"] = ov.ServiceTel
	return v
}
