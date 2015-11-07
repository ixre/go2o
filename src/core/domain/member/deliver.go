/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 21:09
 * description :
 * history :
 */

package member

import (
	"go2o/src/core/domain/interface/member"
	"strings"
)

var _ member.IDeliver = new(deliver)

type deliver struct {
	value     *member.DeliverAddress
	memberRep member.IMemberRep
}

func newDeliver(v *member.DeliverAddress, memberRep member.IMemberRep) member.IDeliver {
	return &deliver{
		value:     v,
		memberRep: memberRep,
	}
}

func (this *deliver) GetDomainId() int {
	return this.value.Id
}

func (this *deliver) GetValue() member.DeliverAddress {
	return *this.value
}

func (this *deliver) SetValue(v *member.DeliverAddress) error {
	if this.value.MemberId == v.MemberId {
		v.Address = strings.TrimSpace(v.Address)
		v.RealName = strings.TrimSpace(v.RealName)
		v.Phone = strings.TrimSpace(v.Phone)

		if err := this.checkValue(v); err != nil {
			return err
		}

		this.value = v
	}
	return nil
}

func (this *deliver) checkValue(v *member.DeliverAddress) error {
	if len(v.Address) < 6 {
		return member.ErrDeliverAddressLen
	}

	if len(v.RealName) < 2 {
		return member.ErrDeliverContactPersonName
	}

	if !phoneRegex.MatchString(v.Phone) {
		return member.ErrDeliverContactPhone
	}
	return nil
}

func (this *deliver) Save() (int, error) {
	return this.memberRep.SaveDeliver(this.value)
}
