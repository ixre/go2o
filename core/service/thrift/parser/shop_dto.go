package parser

import (
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/service/thrift/idl/gen-go/define"
	"strconv"
	"strings"
)

func getShopDto(s shop.IShop) *define.Shop {
	b := s.GetValue()
	dto := &define.Shop{
		ID:       s.GetDomainId(),
		VendorId: b.VendorId,
		ShopType: b.ShopType,
		State:    b.State,
		Name:     b.Name,
		Data:     make(map[string]string),
	}
	return dto
}

func ParseOnlineShop(s shop.IShop) *define.Shop {
	dto := getShopDto(s)
	o := s.(shop.IOnlineShop).GetShopValue()
	dto.Data["ShopId"] = strconv.Itoa(int(o.ShopId))
	dto.Data["ServiceTel"] = o.ServiceTel
	dto.Data["Alias"] = o.Alias
	dto.Data["Host"] = o.Host
	dto.Data["Logo"] = o.Logo
	dto.Data["ShopTitle"] = o.ShopTitle
	dto.Data["ShopNotice"] = o.ShopNotice
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

func Parse2OnlineShop(s *define.Shop) (*shop.Shop, *shop.OnlineShop) {
	sv := parse2Shop(s)
	ov := &shop.OnlineShop{}
	id, _ := strconv.Atoi(s.Data["ShopId"])
	ov.ShopId = int32(id)
	ov.Address = s.Data["Address"]
	ov.ServiceTel = s.Data["ServiceTel"]
	ov.Logo = s.Data["Logo"]
	ov.ShopNotice = s.Data["ShopNotice"]
	ov.ShopTitle = s.Data["ShopTitle"]
	return sv, ov
}
func parse2Shop(s *define.Shop) *shop.Shop {
	return &shop.Shop{
		Id:       s.ID,
		Name:     s.Name,
		VendorId: s.VendorId,
		ShopType: s.ShopType,
		State:    s.ShopType,
	}
}
