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

func (this *partnerService) SaveSaleConf(partnerId int, v *partner.SaleConf) error {
	v.PtId = partnerId
	return this.partnerRep.SaveSaleConf(v)
}

func (this *partnerService) SaveSiteConf(partnerId int, v *partner.SiteConf) error {
	v.PtId = partnerId
	return this.partnerRep.SaveSiteConf(v)
}
func (this *partnerService) GetSaleConf(partnerId int) *partner.SaleConf {
	pt := this.partnerRep.CreatePartner(&partner.ValuePartner{Id: partnerId})
	conf := pt.GetSaleConf()
	return &conf
}

func (this *partnerService) GetSiteConf(partnerId int) *partner.SiteConf {
	pt := this.partnerRep.CreatePartner(&partner.ValuePartner{Id: partnerId})
	conf := pt.GetSiteConf()
	return &conf
}
