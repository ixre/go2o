/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package dao

import (
	"github.com/atnet/gof"
	"go2o/src/core"
)

var (
	context  gof.App
	comm_dao *commDao
)

func init() {
	context = glob.CurrContext()
}

func Common() *commDao {
	if comm_dao == nil {
		comm_dao = &commDao{context.Db()}
	}
	return comm_dao
}
