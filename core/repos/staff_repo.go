package repos

import (
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/gof/db/orm"
)

var _ staff.IStaffRepo = new(employeeRepoImpl)

type employeeRepoImpl struct {
	_orm orm.Orm
}

func NewStaffRepo(o orm.Orm) staff.IStaffRepo {
	return &employeeRepoImpl{
		_orm: o,
	}
}
