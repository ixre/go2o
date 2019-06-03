/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:40
 * description :
 * history :
 */
package user

import (
	"errors"
	"go2o/core/domain/interface/merchant/user"
)

var _ user.IPerson = new(Person)

type Person struct {
	value *user.PersonValue
	rep   user.IUserRepo
}

func newPerson(v *user.PersonValue, rep user.IUserRepo) user.IPerson {
	return &Person{
		value: v,
		rep:   rep,
	}
}

func (p *Person) GetDomainId() int32 {
	return p.value.Id
}

func (p *Person) GetValue() user.PersonValue {
	return *p.value
}

func (p *Person) SetValue(v *user.PersonValue) error {
	if v.Id == p.value.Id && v.Id > 0 {
		p.value = v
		return nil
	}
	return errors.New("no such value")
}

func (p *Person) Save() (int32, error) {
	return p.rep.SavePerson(p.value)
}
