package dao

import (
	"github.com/ixre/go2o/core/dao/model"
)

type IRbacDao interface {
	// auto generate by gof
	// Get 部门
	GetDepart(primary interface{}) *model.RbacDepart
	// GetBy 部门
	GetDepartBy(where string, v ...interface{}) *model.RbacDepart
	// Count 部门 by condition
	CountPermDept(where string, v ...interface{}) (int, error)
	// Select 部门
	SelectPermDept(where string, v ...interface{}) []*model.RbacDepart
	// Save 部门
	SaveDepart(v *model.RbacDepart) (int, error)
	// Delete 部门
	DeleteDepart(primary interface{}) error
	// Batch Delete 部门
	BatchDeleteDepart(where string, v ...interface{}) (int64, error)

	// Get 岗位
	GetJob(primary interface{}) *model.RbacJob
	// GetBy 岗位
	GetJobBy(where string, v ...interface{}) *model.RbacJob
	// Count 岗位 by condition
	CountPermJob(where string, v ...interface{}) (int, error)
	// Select 岗位
	SelectPermJob(where string, v ...interface{}) []*model.RbacJob
	// Save 岗位
	SaveJob(v *model.RbacJob) (int, error)
	// Delete 岗位
	DeleteJob(primary interface{}) error
	// Batch Delete 岗位
	BatchDeleteJob(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 系统用户
	GetUser(primary interface{}) *model.RbacUser
	// GetBy 系统用户
	GetUserBy(where string, v ...interface{}) *model.RbacUser
	// Count 系统用户 by condition
	CountPermUser(where string, v ...interface{}) (int, error)
	// Select 系统用户
	SelectPermUser(where string, v ...interface{}) []*model.RbacUser
	// Save 系统用户
	SaveUser(v *model.RbacUser) (int, error)
	// Delete 系统用户
	DeleteUser(primary interface{}) error
	// Batch Delete 系统用户
	BatchDeleteUser(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermUser(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色
	GetRole(primary interface{}) *model.RbacRole
	// GetBy 角色
	GetRoleBy(where string, v ...interface{}) *model.RbacRole
	// Count 角色 by condition
	CountPermRole(where string, v ...interface{}) (int, error)
	// Select 角色
	SelectPermRole(where string, v ...interface{}) []*model.RbacRole
	// Save 角色
	SavePermRole(v *model.RbacRole) (int, error)
	// Delete 角色
	DeletePermRole(primary interface{}) error
	// Batch Delete 角色
	BatchDeletePermRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get PermRes
	GetRbacResource(primary interface{}) *model.RbacRes
	// GetBy PermRes
	GetRbacResourceBy(where string, v ...interface{}) *model.RbacRes
	// Count PermRes by condition
	CountPermRes(where string, v ...interface{}) (int, error)
	// Select PermRes
	SelectPermRes(where string, v ...interface{}) []*model.RbacRes
	// GetMaxResourceSortNum 获取最大的排列序号
	GetMaxResourceSortNum(parentId int) int
	// GetMaxResourceSortNum 获取最大的Key
	GetMaxResouceKey(parentId int) string
	// Save PermRes
	SaveRbacResource(v *model.RbacRes) (int, error)
	// Delete PermRes
	DeleteRbacResource(primary interface{}) error
	// Batch Delete PermRes
	BatchDeleteRbacResource(where string, v ...interface{}) (int64, error)

	// Get 用户角色关联
	GetUserRole(primary interface{}) *model.RbacUserRole
	// GetBy 用户角色关联
	GetUserRoleBy(where string, v ...interface{}) *model.RbacUserRole
	// Count 用户角色关联 by condition
	CountPermUserRole(where string, v ...interface{}) (int, error)
	// Select 用户角色关联
	SelectPermUserRole(where string, v ...interface{}) []*model.RbacUserRole
	// Save 用户角色关联
	SaveUserRole(v *model.RbacUserRole) (int, error)
	// Delete 用户角色关联
	DeleteUserRole(primary interface{}) error
	// Batch Delete 用户角色关联
	BatchDeleteUserRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermUserRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色部门关联
	GetRoleDept(primary interface{}) *model.RbacRoleDept
	// GetBy 角色部门关联
	GetRoleDeptBy(where string, v ...interface{}) *model.RbacRoleDept
	// Count 角色部门关联 by condition
	CountPermRoleDept(where string, v ...interface{}) (int, error)
	// Select 角色部门关联
	SelectPermRoleDept(where string, v ...interface{}) []*model.RbacRoleDept
	// Save 角色部门关联
	SavePermRoleDept(v *model.RbacRoleDept) (int, error)
	// Delete 角色部门关联
	DeletePermRoleDept(primary interface{}) error
	// Batch Delete 角色部门关联
	BatchDeletePermRoleDept(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRoleDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色菜单关联
	GetRoleRes(primary interface{}) *model.RbacRoleRes
	// GetBy 角色菜单关联
	GetRoleResBy(where string, v ...interface{}) *model.RbacRoleRes
	// Count 角色菜单关联 by condition
	CountPermRoleRes(where string, v ...interface{}) (int, error)
	// Select 角色菜单关联
	SelectPermRoleRes(where string, v ...interface{}) []*model.RbacRoleRes
	// Save 角色菜单关联
	SavePermRoleRes(v *model.RbacRoleRes) (int, error)
	// Delete 角色菜单关联
	DeletePermRoleRes(primary interface{}) error
	// Batch Delete 角色菜单关联
	BatchDeletePermRoleRes(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRoleRes(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// 获取角色关联资源信息
	GetRoleResList(roles []int) []*model.RbacRoleRes
	// 获取用户的角色
	GetUserRoles(id int) []*model.RbacUserRole
	// 获取角色拥有的资源
	GetRoleResources(roles []int) []*model.RbacRes

	// PagingQueryLoginLog Query paging data
	PagingQueryLoginLog(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
}
