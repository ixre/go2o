package employee

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/employee"
)

var _ employee.IEmployeeManager = new(employeeManagerImpl)

type employeeManagerImpl struct {
}

func NewEmployeeManager(mch merchant.IMerchant,employeeRepo employee.IEmployeeRepo) employee.IEmployeeManager {
	return &employeeManagerImpl{}
}

// Create implements employee.IEmployeeManager.
func (e *employeeManagerImpl) Create(memberId int) error {
	panic("unimplemented")
}
