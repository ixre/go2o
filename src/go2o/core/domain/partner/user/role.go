/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 15:42
 * description :
 * history :
 */
package user

import (
	"errors"
	"go2o/core/domain/interface/partner/user"
)

var _ user.IRole = new(Role)

type Role struct {
	value *user.RoleValue
	rep   user.IUserRep
}

func newRole(v *user.RoleValue, rep user.IUserRep) user.IRole {
	return &Role{
		value: v,
		rep:   rep,
	}
}

func (this *Role) GetDomainId() int {
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

func (this *Role) Save() (int, error) {
	return this.rep.SaveRole(this.value)
}
