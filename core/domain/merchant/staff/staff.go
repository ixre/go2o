package employee

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
)

var _ staff.IStaffManager = new(employeeManagerImpl)

type employeeManagerImpl struct {
}

func NewEmployeeManager(mch merchant.IMerchant, employeeRepo staff.IStaffRepo) staff.IStaffManager {
	return &employeeManagerImpl{}
}

// Create implements staff.IStaffManager.
func (e *employeeManagerImpl) Create(memberId int) error {
	panic("unimplemented")
}
