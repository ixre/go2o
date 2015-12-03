/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:46
 * description :
 * history :
 */
package user

import (
	"go2o/src/core/domain/interface/partner/user"
)

var _ user.IUser = new(User)
var _ user.IDeliveryStaff = new(User)

type User struct {
	person user.IPerson
	rep    user.IUserRep
}

func newUser(v *user.PersonValue, rep user.IUserRep) user.IUser {
	var person = newPerson(v, rep)
	return &User{
		person: person,
		rep:    rep,
	}
}

// 获取人员信息
func (this *User) GetPerson() user.IPerson {
	return this.person
}

// 获取凭据
func (this *User) GetCredential(sign string) *user.CredentialValue {
	//todo: not will used
	return nil
}

// 保存凭据
func (this *User) SaveCredential(v *user.CredentialValue) error {
	_, err := this.rep.SaveCredential(v)
	return err
}
