/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 17:46
 * description :
 * history :
 */
package user

import (
	"go2o/core/domain/interface/partner/user"
)

var _ user.IUserManager = new(UserManager)

type UserManager struct {
	partnerId int
	rep       user.IUserRep
}

func NewUserManager(partnerId int, rep user.IUserRep) user.IUserManager {
	return &UserManager{
		partnerId: partnerId,
		rep:       rep,
	}
}

// 获取单个用户
func (this *UserManager) GetUser(id int) IUser {
	v := this.rep.GetUserValue(id)
	if v != nil {
		return newUser(v, this.rep)
	}
	return nil
}

// 获取所有配送员
func (this *UserManager) GetDeliveryStaff() []IDeliveryStaff {
	list := this.rep.GetDeliveryStaffPersons(this.partnerId)
	var staffs = make([]IDeliveryStaff, len(list))
	for i, v := range list {
		staffs[i] = newUser(v, this.rep)
	}
	return staffs
}
