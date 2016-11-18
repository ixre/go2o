/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2015-02-15 10:22
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/merchant/user"
)

var _ user.IUserRep = new(userRep)

type userRep struct {
	db.Connector
}

func NewUserRep(c db.Connector) user.IUserRep {
	return &userRep{
		Connector: c,
	}
}

// 保存角色
func (this *userRep) SaveRole(v *user.RoleValue) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 保存人员
func (this *userRep) SavePerson(v *user.PersonValue) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 保存凭据
func (this *userRep) SaveCredential(v *user.CredentialValue) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 获取人员
func (this *userRep) GetPersonValue(id int32) *user.PersonValue {
	e := new(user.PersonValue)
	err := this.Connector.GetOrm().Get(e, id)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送人员
func (this *userRep) GetDeliveryStaffPersons(mchId int32) []*user.PersonValue {
	e := make([]*user.PersonValue, 0)
	err := this.Connector.GetOrm().Select(e, "select * from usr_person")
	if err != nil {
		return nil
	}
	return e
}
