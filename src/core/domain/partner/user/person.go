/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 17:40
 * description :
 * history :
 */
package user

import (
	"errors"
	"go2o/src/core/domain/interface/partner/user"
)

var _ user.IPerson = new(Person)

type Person struct {
	value *user.PersonValue
	rep   user.IUserRep
}

func newPerson(v *user.PersonValue, rep user.IUserRep) user.IPerson {
	return &Person{
		value: v,
		rep:   rep,
	}
}

func (this *Person) GetDomainId() int {
	return this.value.Id
}

func (this *Person) GetValue() user.PersonValue {
	return *this.value
}

func (this *Person) SetValue(v *user.PersonValue) error {
	if v.Id == this.value.Id && v.Id > 0 {
		this.value = v
		return nil
	}
	return errors.New("no such value")
}

func (this *Person) Save() (int, error) {
	return this.rep.SavePerson(this.value)
}
