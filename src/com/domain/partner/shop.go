/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-23 07:55
 * description :
 * history :
 */

package partner

import (
	"com/domain/interface/partner"
)

var _ partner.IShop = new(Shop)

type Shop struct {
	partnerRep partner.IPartnerRep
	value      *partner.ValueShop
	partner    *Partner
}

func newShop(pt *Partner, v *partner.ValueShop, partnerRep partner.IPartnerRep) partner.IShop {
	return &Shop{
		partnerRep: partnerRep,
		value:      v,
		partner:    pt,
	}
}

func (this *Shop) GetDomainId() int {
	return this.value.Id
}

func (this *Shop) GetValue() partner.ValueShop {
	return *this.value
}

func (this *Shop) SetValue(v *partner.ValueShop) error {
	this.value = v
	return nil
}

func (this *Shop) Save() (int, error) {
	this.partner.clearShopCache()
	return this.partnerRep.SaveShop(this.value)
}
