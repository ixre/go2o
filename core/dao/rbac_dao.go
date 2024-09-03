package dao

import rbac "github.com/ixre/go2o/core/domain/interface/rabc"

type IRbacDao interface {
	// auto generate by gof
	// Get 部门
	GetDepart(primary interface{}) *rbac.RbacDepart
	// GetBy 部门
	GetDepartBy(where string, v ...interface{}) *rbac.RbacDepart
	// Count 部门 by condition
	CountPermDept(where string, v ...interface{}) (int, error)
	// Select 部门
	SelectPermDept(where string, v ...interface{}) []*rbac.RbacDepart
	// Save 部门
	SaveDepart(v *rbac.RbacDepart) (int, error)
	// Delete 部门
	DeleteDepart(primary interface{}) error
	// Batch Delete 部门
	BatchDeleteDepart(where string, v ...interface{}) (int64, error)

	// Get 岗位
	GetJob(primary interface{}) *rbac.RbacJob
	// GetBy 岗位
	GetJobBy(where string, v ...interface{}) *rbac.RbacJob
	// Count 岗位 by condition
	CountPermJob(where string, v ...interface{}) (int, error)
	// Select 岗位
	SelectPermJob(where string, v ...interface{}) []*rbac.RbacJob
	// Save 岗位
	SaveJob(v *rbac.RbacJob) (int, error)
	// Delete 岗位
	DeleteJob(primary interface{}) error
	// Batch Delete 岗位
	BatchDeleteJob(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 系统用户
	GetUser(primary interface{}) *rbac.RbacUser
	// GetBy 系统用户
	GetUserBy(where string, v ...interface{}) *rbac.RbacUser
	// Count 系统用户 by condition
	CountPermUser(where string, v ...interface{}) (int, error)
	// Select 系统用户
	SelectPermUser(where string, v ...interface{}) []*rbac.RbacUser
	// Save 系统用户
	SaveUser(v *rbac.RbacUser) (int, error)
	// Delete 系统用户
	DeleteUser(primary interface{}) error
	// Batch Delete 系统用户
	BatchDeleteUser(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingPermUser(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色
	GetRole(primary interface{}) *rbac.RbacRole
	// GetBy 角色
	GetRoleBy(where string, v ...interface{}) *rbac.RbacRole
	// Count 角色 by condition
	CountPermRole(where string, v ...interface{}) (int, error)
	// Select 角色
	SelectPermRole(where string, v ...interface{}) []*rbac.RbacRole
	// Save 角色
	SavePermRole(v *rbac.RbacRole) (int, error)
	// Delete 角色
	DeletePermRole(primary interface{}) error
	// Batch Delete 角色
	BatchDeletePermRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingPermRole(begin, end int, where string) (total int, rows []map[string]interface{})

	// Get PermRes
	GetRbacResource(primary interface{}) *rbac.RbacRes
	// GetBy PermRes
	GetRbacResourceBy(where string, v ...interface{}) *rbac.RbacRes
	// Count PermRes by condition
	CountPermRes(where string, v ...interface{}) (int, error)
	// Select PermRes
	SelectPermRes(where string, v ...interface{}) []*rbac.RbacRes
	// GetMaxResourceSortNum 获取最大的排列序号
	GetMaxResourceSortNum(parentId int) int
	// GetMaxResourceSortNum 获取最大的Key
	GetMaxResouceKey(parentId int) string
	// Save PermRes
	SaveRbacResource(v *rbac.RbacRes) (int, error)
	// Delete PermRes
	DeleteRbacResource(primary interface{}) error
	// Batch Delete PermRes
	BatchDeleteRbacResource(where string, v ...interface{}) (int64, error)

	// Get 用户角色关联
	GetUserRole(primary interface{}) *rbac.RbacUserRole
	// GetBy 用户角色关联
	GetUserRoleBy(where string, v ...interface{}) *rbac.RbacUserRole
	// Count 用户角色关联 by condition
	CountPermUserRole(where string, v ...interface{}) (int, error)
	// Select 用户角色关联
	SelectPermUserRole(where string, v ...interface{}) []*rbac.RbacUserRole
	// Save 用户角色关联
	SaveUserRole(v *rbac.RbacUserRole) (int, error)
	// Delete 用户角色关联
	DeleteUserRole(primary interface{}) error
	// Batch Delete 用户角色关联
	BatchDeleteUserRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingPermUserRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色部门关联
	GetRoleDept(primary interface{}) *rbac.RbacRoleDept
	// GetBy 角色部门关联
	GetRoleDeptBy(where string, v ...interface{}) *rbac.RbacRoleDept
	// Count 角色部门关联 by condition
	CountPermRoleDept(where string, v ...interface{}) (int, error)
	// Select 角色部门关联
	SelectPermRoleDept(where string, v ...interface{}) []*rbac.RbacRoleDept
	// Save 角色部门关联
	SavePermRoleDept(v *rbac.RbacRoleDept) (int, error)
	// Delete 角色部门关联
	DeletePermRoleDept(primary interface{}) error
	// Batch Delete 角色部门关联
	BatchDeletePermRoleDept(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingPermRoleDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色菜单关联
	GetRoleRes(primary interface{}) *rbac.RbacRoleRes
	// GetBy 角色菜单关联
	GetRoleResBy(where string, v ...interface{}) *rbac.RbacRoleRes
	// Count 角色菜单关联 by condition
	CountPermRoleRes(where string, v ...interface{}) (int, error)
	// Select 角色菜单关联
	SelectPermRoleRes(where string, v ...interface{}) []*rbac.RbacRoleRes
	// Save 角色菜单关联
	SavePermRoleRes(v *rbac.RbacRoleRes) (int, error)
	// Delete 角色菜单关联
	DeletePermRoleRes(primary interface{}) error
	// Batch Delete 角色菜单关联
	BatchDeletePermRoleRes(where string, v ...interface{}) (int64, error)
	// Params paging data
	QueryPagingPermRoleRes(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// 获取角色关联资源信息
	GetRoleResList(roles []int) []*rbac.RbacRoleRes
	// 获取用户的角色
	GetUserRoles(id int) []*rbac.RbacUserRole
	// 获取角色拥有的资源
	GetRoleResources(roles []int) []*rbac.RbacRes

	// QueryPagingLoginLog Query paging data
	QueryPagingLoginLog(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
}
