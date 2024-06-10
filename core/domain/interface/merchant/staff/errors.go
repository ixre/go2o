package staff

import "github.com/ixre/go2o/core/infrastructure/domain"

var (
	ErrStaffAlreadyExists = domain.NewError("staff_already_exists", "员工已存在")
)
