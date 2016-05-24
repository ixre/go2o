/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 07:55
 * description :
 * history :
 */

package merchant

import (
	"go2o/core/domain/interface/merchant"
)

var _ merchant.IShop = new(Shop)

type Shop struct {
	partnerRep merchant.IMerchantRep
	value      *merchant.ValueShop
	partner    *Merchant
	lng        float64
	lat        float64
}

func newShop(pt *Merchant, v *merchant.ValueShop, partnerRep merchant.IMerchantRep) merchant.IShop {
	return &Shop{
		partnerRep: partnerRep,
		value:      v,
		partner:    pt,
	}
}

func (this *Shop) GetDomainId() int {
	return this.value.Id
}

func (this *Shop) GetValue() merchant.ValueShop {
	return *this.value
}

func (this *Shop) SetValue(v *merchant.ValueShop) error {
	//	if this.value.Address != v.Address ||
	//		len(this.value.Location) == 0 {
	//		lng, lat, err := lbs.GetLocation(v.Address)
	//		if err != nil {
	//			return err
	//		}
	//		this.value.Location = fmt.Sprintf("%f,%f", lng, lat)
	//}
	//this.value.DeliverRadius = v.DeliverRadius
	this.value.Address = v.Address
	this.value.Name = v.Name
	this.value.SortNumber = v.SortNumber
	this.value.MerchantId = this.partner.GetAggregateRootId()
	this.value.Phone = v.Phone
	this.value.State = v.State
	return nil
}

func (this *Shop) Save() (int, error) {
	this.partner.clearShopCache()
	return this.partnerRep.SaveShop(this.value)
}

//// 获取经维度
//func (this *Shop) GetLngLat() (float64, float64) {
//	if this.lng == 0 || this.lat == 0 {
//		var err error
//		this.lng, this.lat, err = lbs.GetLocation(this.value.Address)
//		if err != nil {
//			log.Error(err)
//		}
//	}
//	return this.lng, this.lat
//}

//// 是否可以配送
//// 返回是否可以配送，以及距离(米)
//func (this *Shop) CanDeliver(lng, lat float64) (bool, int) {
//	shopLng, shopLat := this.GetLngLat()
//	distance := lbs.GetLocDistance(shopLng, shopLat, lng, lat)
//	i := int(distance)
//	return i <= this.value.DeliverRadius*1000, i
//}

//// 是否可以配送
//// 返回是否可以配送，以及距离(米)
//func (this *Shop) CanDeliverTo(address string) (bool, int) {
//	lng, lat, err := lbs.GetLocation(address)
//	if err != nil {
//		log.Error(err)
//		return false, -1
//	}
//	return this.CanDeliver(lng, lat)
//}
