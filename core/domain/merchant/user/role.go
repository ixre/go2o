/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:42
 * description :
 * history :
 */
package user

import (
	"errors"
	"go2o/core/domain/interface/merchant/user"
)

var _ user.IRole = new(Role)

type Role struct {
	value *user.RoleValue
	rep   user.IUserRepo
}

func newRole(v *user.RoleValue, rep user.IUserRepo) user.IRole {
	return &Role{
		value: v,
		rep:   rep,
	}
}

func (this *Role) GetDomainId() int32 {
	return this.value.Id
}

func (this *Role) GetValue() user.RoleValue {
	return *this.value
}

func (this *Role) SetValue(v *user.RoleValue) error {
	if v.Id == this.value.Id && v.Id > 0 {
		this.value = v
		return nil
	}
	return errors.New("no such value")
}

func (this *Role) Save() (int32, error) {
	return this.rep.SaveRole(this.value)
}
