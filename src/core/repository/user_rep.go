/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2015-02-15 10:22
 * description :
 * history :
 */
package repository

import (
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/partner/user"
)

var _ user.IUserRep = new(UserRep)

type UserRep struct {
	db.Connector
}

func NewUserRep(c db.Connector) user.IUserRep {
	return &UserRep{
		Connector: c,
	}
}

// 保存角色
func (this *UserRep) SaveRole(v *user.RoleValue) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM usr_role", &v.Id)
	}
	return v.Id, err
}

// 保存人员
func (this *UserRep) SavePerson(v *user.PersonValue) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM usr_person", &v.Id)
	}
	return v.Id, err
}

// 保存凭据
func (this *UserRep) SaveCredential(v *user.CredentialValue) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM usr_credential", &v.Id)
	}
	return v.Id, err
}

// 获取人员
func (this *UserRep) GetPersonValue(id int) *user.PersonValue {
	e := new(user.PersonValue)
	err := this.Connector.GetOrm().Get(e, id)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送人员
func (this *UserRep) GetDeliveryStaffPersons(partnerId int) []*user.PersonValue {
	e := make([]*user.PersonValue, 0)
	err := this.Connector.GetOrm().Select(e, "select * from usr_person")
	if err != nil {
		return nil
	}
	return e
}
