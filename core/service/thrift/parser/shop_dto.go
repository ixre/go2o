package parser

import (
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/valueobject"
	"go2o/gen-code/thrift/define"
	"strconv"
	"strings"
)

func getShopDto(s shop.IShop) *define.Shop {
	b := s.GetValue()
	dto := &define.Shop{
		ID:           s.GetDomainId(),
		VendorId:     b.VendorId,
		ShopType:     b.ShopType,
		State:        b.State,
		OpeningState: b.OpeningState,
		Name:         b.Name,
		Data:         make(map[string]string),
	}
	return dto
}

func parse2Shop(s *define.Shop) *shop.Shop {
	return &shop.Shop{
		Id:           s.ID,
		Name:         s.Name,
		VendorId:     s.VendorId,
		ShopType:     s.ShopType,
		State:        s.State,
		OpeningState: s.OpeningState,
	}
}

func ParseOnlineShop(s shop.IShop) *define.Store {
	b := s.GetValue()
	o := s.(shop.IOnlineShop).GetShopValue()
	dto := &define.Store{
		ID:           s.GetDomainId(),
		VendorId:     b.VendorId,
		State:        b.State,
		OpeningState: b.OpeningState,
		Name:         b.Name,
		Alias:        o.Alias,
		StorePhone:   o.ServiceTel,
		Host:         o.Host,
		Logo:         o.Logo,
		StoreTitle:   o.ShopTitle,
		StoreNotice:  o.ShopNotice,
	}
	return dto
}

func ParseOfflineShop(s shop.IShop, valRepo valueobject.IValueRepo) *define.Shop {
	dto := getShopDto(s)
	o := s.(shop.IOfflineShop).GetShopValue()
	areaNames := valRepo.GetAreaNames([]int32{o.Province, o.City, o.District})
	dto.Data["ShopId"] = strconv.Itoa(int(o.ShopId))
	dto.Data["ServiceTel"] = o.Tel
	dto.Data["Province"] = strconv.Itoa(int(o.Province))
	dto.Data["City"] = strconv.Itoa(int(o.City))
	dto.Data["District"] = strconv.Itoa(int(o.District))
	dto.Data["Location"] = strings.Join(areaNames, " ")
	dto.Data["Address"] = o.Address
	dto.Data["Lng"] = strconv.FormatFloat(float64(o.Lng), 'g', 2, 32)
	dto.Data["Lat"] = strconv.FormatFloat(float64(o.Lat), 'g', 2, 32)
	dto.Data["CoverRadius"] = strconv.Itoa(o.CoverRadius)
	return dto
}

func Parse2OnlineShop(s *define.Store) (*shop.Shop, *shop.OnlineShop) {
	sv := &shop.Shop{
		Id:           s.ID,
		Name:         s.Name,
		VendorId:     s.VendorId,
		ShopType:     shop.TypeOnlineShop,
		State:        s.State,
		OpeningState: s.OpeningState,
	}
	ov := &shop.OnlineShop{}
	ov.ShopId = s.ID
	ov.Address = "" //todo:???
	ov.ServiceTel = s.StorePhone
	ov.Logo = s.Logo
	ov.ShopNotice = s.StoreNotice
	ov.ShopTitle = s.StoreTitle
	return sv, ov
}
