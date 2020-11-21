package dao

import (
	"go2o/core/dao/model"
)

type IPermUserDao interface {
	// auto generate by gof
	// Get 系统用户
	Get(primary interface{}) *model.PermUser
	// GetBy 系统用户
	GetBy(where string, v ...interface{}) *model.PermUser
	// Count 系统用户 by condition
	Count(where string, v ...interface{}) (int, error)
	// Select 系统用户
	Select(where string, v ...interface{}) []*model.PermUser
	// Save 系统用户
	Save(v *model.PermUser) (int, error)
	// Delete 系统用户
	Delete(primary interface{}) error
	// Batch Delete 系统用户
	BatchDelete(where string, v ...interface{}) (int64, error)
	// Query paging data
	PagingQuery(begin, end int, where, orderBy string) (num int, rows []map[string]interface{})
}
