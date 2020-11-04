/**
 * Copyright 2015 @ to2.net.
 * name : member_test.go
 * author : jarryliu
 * date : 2016-06-25 07:51
 * description :
 * history :
 */
package member

import (
	"go2o/core/domain/interface/member"
	"testing"
)

func TestDeliverAddressSave(t *testing.T) {
	m := NewMember(NewMemberManager(nil, nil), &member.Member{
		Id: 1,
	}, nil, nil, nil)
	d := m.Profile().GetDeliverAddress()[0]
	v := d.GetValue()
	v.Province = 440000
	v.City = 440600
	v.District = 440605
	d.SetValue(&v)
	d.Save()
}
