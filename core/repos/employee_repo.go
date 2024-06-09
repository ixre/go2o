package repos

import (
	"github.com/ixre/go2o/core/domain/interface/merchant/employee"
	"github.com/ixre/gof/db/orm"
)

var _ employee.IEmployeeRepo = new(employeeRepoImpl)

type employeeRepoImpl struct {
	_orm orm.Orm
}

func NewEmployeeRepo(o orm.Orm) employee.IEmployeeRepo {
	return &employeeRepoImpl{
		_orm: o,
	}
}
