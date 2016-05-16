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
	"go2o/src/core/domain/interface/merchant/user"
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
func (this *userRep) SaveRole(v *user.RoleValue) (int, error) {
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
func (this *userRep) SavePerson(v *user.PersonValue) (int, error) {
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
func (this *userRep) SaveCredential(v *user.CredentialValue) (int, error) {
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
func (this *userRep) GetPersonValue(id int) *user.PersonValue {
	e := new(user.PersonValue)
	err := this.Connector.GetOrm().Get(e, id)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送人员
func (this *userRep) GetDeliveryStaffPersons(partnerId int) []*user.PersonValue {
	e := make([]*user.PersonValue, 0)
	err := this.Connector.GetOrm().Select(e, "select * from usr_person")
	if err != nil {
		return nil
	}
	return e
}
