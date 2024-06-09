package employee

type (
	// IEmployeeManager 员工管理接口
	IEmployeeManager interface {
		// Create 创建员工
		Create(memberId int) error
	}

	// IEmployeeRepo 员工数据访问接口
	IEmployeeRepo interface {
	}
)
