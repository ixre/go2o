/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-19 17:49
 * description :
 * history :
 */

package dproxy

import (
	"com/domain/interface/partner"
)

type partnerService struct {
	partnerRep partner.IPartnerRep
}

func (this *partnerService) GetPartner(partnerId int) *partner.ValuePartner {
	pt := this.partnerRep.GetPartner(partnerId)
	if pt != nil {
		v := pt.GetValue()
		return &v
	}
	return nil
}

func (this *partnerService) SaveSaleConf(partnerId int, v *partner.SaleConf) error {
	v.PartnerId = partnerId
	return this.partnerRep.SaveSaleConf(v)
}

func (this *partnerService) SaveSiteConf(partnerId int, v *partner.SiteConf) error {
	v.PartnerId = partnerId
	return this.partnerRep.SaveSiteConf(v)
}
func (this *partnerService) GetSaleConf(partnerId int) *partner.SaleConf {
	pt := this.partnerRep.GetPartner(partnerId)
	conf := pt.GetSaleConf()
	return &conf
}

func (this *partnerService) GetSiteConf(partnerId int) *partner.SiteConf {
	pt := this.partnerRep.GetPartner(partnerId)
	conf := pt.GetSiteConf()
	return &conf
}

func (this *partnerService) GetShopsOfPartner(partnerId int) []*partner.ValueShop {
	pt := this.partnerRep.GetPartner(partnerId)
	shops := pt.GetShops()
	sv := make([]*partner.ValueShop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

func (this *partnerService) GetShopValueById(partnerId, shopId int) *partner.ValueShop {
	pt := this.partnerRep.GetPartner(partnerId)
	v := pt.GetShop(shopId).GetValue()
	return &v
}

func (this *partnerService) SaveShop(partnerId int, v *partner.ValueShop) (int, error) {
	pt := this.partnerRep.GetPartner(partnerId)
	shop := pt.GetShop(v.Id)
	shop.SetValue(v)
	return shop.Save()
}

func (this *partnerService) DeleteShop(partnerId, shopId int) error {
	pt := this.partnerRep.GetPartner(partnerId)
	return pt.DeleteShop(shopId)
}
