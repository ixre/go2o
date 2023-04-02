package dao

import (
	"github.com/ixre/go2o/core/dao/model"
)

type IRbacDao interface {
	// auto generate by gof
	// Get 部门
	GetDepart(primary interface{}) *model.PermDept
	// GetBy 部门
	GetDepartBy(where string, v ...interface{}) *model.PermDept
	// Count 部门 by condition
	CountPermDept(where string, v ...interface{}) (int, error)
	// Select 部门
	SelectPermDept(where string, v ...interface{}) []*model.PermDept
	// Save 部门
	SaveDepart(v *model.PermDept) (int, error)
	// Delete 部门
	DeleteDepart(primary interface{}) error
	// Batch Delete 部门
	BatchDeleteDepart(where string, v ...interface{}) (int64, error)

	// Get 岗位
	GetJob(primary interface{}) *model.PermJob
	// GetBy 岗位
	GetJobBy(where string, v ...interface{}) *model.PermJob
	// Count 岗位 by condition
	CountPermJob(where string, v ...interface{}) (int, error)
	// Select 岗位
	SelectPermJob(where string, v ...interface{}) []*model.PermJob
	// Save 岗位
	SaveJob(v *model.PermJob) (int, error)
	// Delete 岗位
	DeleteJob(primary interface{}) error
	// Batch Delete 岗位
	BatchDeleteJob(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 系统用户
	GetUser(primary interface{}) *model.PermUser
	// GetBy 系统用户
	GetUserBy(where string, v ...interface{}) *model.PermUser
	// Count 系统用户 by condition
	CountPermUser(where string, v ...interface{}) (int, error)
	// Select 系统用户
	SelectPermUser(where string, v ...interface{}) []*model.PermUser
	// Save 系统用户
	SaveUser(v *model.PermUser) (int, error)
	// Delete 系统用户
	DeleteUser(primary interface{}) error
	// Batch Delete 系统用户
	BatchDeleteUser(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermUser(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色
	GetPermRole(primary interface{}) *model.PermRole
	// GetBy 角色
	GetPermRoleBy(where string, v ...interface{}) *model.PermRole
	// Count 角色 by condition
	CountPermRole(where string, v ...interface{}) (int, error)
	// Select 角色
	SelectPermRole(where string, v ...interface{}) []*model.PermRole
	// Save 角色
	SavePermRole(v *model.PermRole) (int, error)
	// Delete 角色
	DeletePermRole(primary interface{}) error
	// Batch Delete 角色
	BatchDeletePermRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get PermRes
	GetPermRes(primary interface{}) *model.PermRes
	// GetBy PermRes
	GetPermResBy(where string, v ...interface{}) *model.PermRes
	// Count PermRes by condition
	CountPermRes(where string, v ...interface{}) (int, error)
	// Select PermRes
	SelectPermRes(where string, v ...interface{}) []*model.PermRes
	// 获取最大的排列序号
	GetMaxResourceSortNum(parentId int) int
	// Save PermRes
	SavePermRes(v *model.PermRes) (int, error)
	// Delete PermRes
	DeletePermRes(primary interface{}) error
	// Batch Delete PermRes
	BatchDeletePermRes(where string, v ...interface{}) (int64, error)

	// Get 用户角色关联
	GetUserRole(primary interface{}) *model.PermUserRole
	// GetBy 用户角色关联
	GetUserRoleBy(where string, v ...interface{}) *model.PermUserRole
	// Count 用户角色关联 by condition
	CountPermUserRole(where string, v ...interface{}) (int, error)
	// Select 用户角色关联
	SelectPermUserRole(where string, v ...interface{}) []*model.PermUserRole
	// Save 用户角色关联
	SaveUserRole(v *model.PermUserRole) (int, error)
	// Delete 用户角色关联
	DeleteUserRole(primary interface{}) error
	// Batch Delete 用户角色关联
	BatchDeleteUserRole(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermUserRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色部门关联
	GetPermRoleDept(primary interface{}) *model.PermRoleDept
	// GetBy 角色部门关联
	GetPermRoleDeptBy(where string, v ...interface{}) *model.PermRoleDept
	// Count 角色部门关联 by condition
	CountPermRoleDept(where string, v ...interface{}) (int, error)
	// Select 角色部门关联
	SelectPermRoleDept(where string, v ...interface{}) []*model.PermRoleDept
	// Save 角色部门关联
	SavePermRoleDept(v *model.PermRoleDept) (int, error)
	// Delete 角色部门关联
	DeletePermRoleDept(primary interface{}) error
	// Batch Delete 角色部门关联
	BatchDeletePermRoleDept(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRoleDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// Get 角色菜单关联
	GetPermRoleRes(primary interface{}) *model.PermRoleRes
	// GetBy 角色菜单关联
	GetPermRoleResBy(where string, v ...interface{}) *model.PermRoleRes
	// Count 角色菜单关联 by condition
	CountPermRoleRes(where string, v ...interface{}) (int, error)
	// Select 角色菜单关联
	SelectPermRoleRes(where string, v ...interface{}) []*model.PermRoleRes
	// Save 角色菜单关联
	SavePermRoleRes(v *model.PermRoleRes) (int, error)
	// Delete 角色菜单关联
	DeletePermRoleRes(primary interface{}) error
	// Batch Delete 角色菜单关联
	BatchDeletePermRoleRes(where string, v ...interface{}) (int64, error)
	// Params paging data
	PagingQueryPermRoleRes(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})

	// 获取角色关联的资源列表
	GetRoleResList(roleId int64) []int64
	// 获取用户的角色
	GetUserRoles(id int64) []*model.PermUserRole
	// 获取角色关联系
	GetRoleResources(roles []int) []*model.PermRes
}
