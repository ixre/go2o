/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 07:55
 * description :
 * history :
 */

package shop

import (
	"errors"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/tmp"
	"github.com/ixre/gof/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	shopAliasRegexp = regexp.MustCompile("^[A-Za-z0-9-]{3,11}$")
)

func NewShop(v *shop.OnlineShop, shopRepo shop.IShopRepo,
	valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) shop.IShop {
	return &onlineShopImpl{
		_shopVal: v,
		valRepo:  valRepo,
		shopRepo: shopRepo,
		registryRepo: registryRepo,
	}
}

var _ shop.IShop = new(onlineShopImpl)
var _ shop.IOnlineShop = new(onlineShopImpl)

type onlineShopImpl struct {
	_mch merchant.IMerchant
	_shopVal     *shop.OnlineShop
	valRepo      valueobject.IValueRepo
	shopRepo     shop.IShopRepo
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
}

func (s *onlineShopImpl) GetDomainId() int {
	return int(s._shopVal.Id)
}

func (s *onlineShopImpl) check(v *shop.Shop) error {
	v.Name = strings.TrimSpace(v.Name)
	if s.checkNameExists(v) {
		return shop.ErrSameNameShopExists
	}
	return nil
}

func (s *onlineShopImpl) checkNameExists(v *shop.Shop) bool {
	i := 0
	tmp.Db().ExecScalar("SELECT COUNT(0) FROM mch_shop WHERE name= $1 AND id <> $2", &i,
		v.Name, v.Id)
	return i > 0
}

func (s *onlineShopImpl) Type() int32 {
	return shop.TypeOnlineShop
}

func (s *onlineShopImpl) GetValue() shop.Shop {
	panic("implement me")
}

func (s *onlineShopImpl) SetValue(*shop.Shop) error {
	panic("implement me")
}

func (s *onlineShopImpl) TurnOn() error {
	panic("implement me")
}

func (s *onlineShopImpl) TurnOff(reason string) error {
	panic("implement me")
}

func (s *onlineShopImpl) Opening() error {
	panic("implement me")
}

func (s *onlineShopImpl) Pause() error {
	panic("implement me")
}

func (s *onlineShopImpl) checkShopAlias(alias string) error {
	alias = strings.ToLower(alias)
	if strings.HasPrefix(alias, "shop") {
		return errors.New("非法的店铺别名")
	}
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

// SetShopValue 设置值
func (s *onlineShopImpl) SetShopValue(v *shop.OnlineShop) (err error) {
	mv := s._mch.GetValue()
	v.Logo = strings.TrimSpace(v.Logo)
	dst := s._shopVal
	if s.GetDomainId() <= 0 {
		unix := time.Now().Unix()
		if v.VendorId <= 0 {
			return merchant.ErrNoSuchMerchant
		}
		dst.Logo = shop.DefaultOnlineShop.Logo
		dst.CreateTime = unix
		dst.State = shop.StateAwaitInitial
		dst.Alias = s.generateShopAlias()
	}
	if len(v.Logo) > 0 {
		dst.Logo = v.Logo
	}
	if len(v.Alias) > 0 && v.Alias != dst.Alias {
		err = s.checkShopAlias(v.Alias)
		if err == nil {
			dst.Alias = v.Alias
		}
	}
	// 判断自营
	if mv.SelfSales == 1{
		dst.Flag |= shop.FlagSelfSale
	}else{
		dst.Flag ^= shop.FlagSelfSale
	}
	dst.Host = v.Host
	dst.Telephone = v.Telephone
	dst.Addr = v.Addr
	dst.State = v.State
	dst.Host = v.Host
	dst.ShopTitle = v.ShopTitle
	dst.ShopName = v.ShopName
	dst.ShopNotice = v.ShopNotice
	return err
}

// GetShopValue 获取值
func (s *onlineShopImpl) GetShopValue() shop.OnlineShop {
	return *s._shopVal
}

// Save 保存
func (s *onlineShopImpl) Save() error {
	if s.GetDomainId() > 0 {
		if len(s._shopVal.Alias) == 0{
			s._shopVal.Alias = s.generateShopAlias()
		}
		_, err := s.shopRepo.SaveOnlineShop(s._shopVal)
		return err
	}
	return s.createShop()
}

func (s *onlineShopImpl) createShop() error {
	if s.shopRepo.CheckShopExists(s._shopVal.VendorId) {
		return shop.ErrSupportSingleOnlineShop
	}
	unix := time.Now().Unix()
	s._shopVal.Alias = s.generateShopAlias()
	s._shopVal.CreateTime = unix
	if s._shopVal.Logo == "" {
		s._shopVal.Logo = shop.DefaultOnlineShop.Logo
	}
	id, err := s.shopRepo.SaveOnlineShop(s._shopVal)
	if err == nil {
		s._shopVal.Id = id
	}
	return err
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

// Data 获取商店信息
func (s *onlineShopImpl) Data() *shop.ComplexShop {
	ov := s._shopVal
	v := &shop.ComplexShop{
		ID:       int64(s.GetDomainId()),
		VendorId: ov.VendorId,
		ShopType: shop.TypeOnlineShop,
		Name:     ov.ShopTitle,
		State:    shop.StateNormal,
		Data:     make(map[string]string),
	}
	v.Data["Host"] = ov.Host
	v.Data["Logo"] = ov.Logo
	v.Data["ServiceTel"] = ov.Telephone
	return v
}

func (s *onlineShopImpl) GetLocateDomain() string {
	panic("implement me")
}

func (s *onlineShopImpl) BindDomain(domain string) error {
	panic("implement me")
}
