/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:46
 * description :
 * history :
 */
package user

import (
	"go2o/core/domain/interface/merchant/user"
)

var _ user.IUserManager = new(UserManager)

type UserManager struct {
	mchId int32
	rep   user.IUserRepo
}

func NewUserManager(mchId int32, rep user.IUserRepo) user.IUserManager {
	return &UserManager{
		mchId: mchId,
		rep:   rep,
	}
}

// 获取单个用户
func (u *UserManager) GetUser(id int32) user.IUser {
	v := u.rep.GetPersonValue(id)
	if v != nil {
		return newUser(v, u.rep)
	}
	return nil
}

// 获取所有配送员
func (u *UserManager) GetDeliveryStaff() []user.IDeliveryStaff {
	list := u.rep.GetDeliveryStaffPersons(u.mchId)
	var staffs = make([]user.IDeliveryStaff, len(list))
	for i, v := range list {
		staffs[i] = newUser(v, u.rep)
	}
	return staffs
}
