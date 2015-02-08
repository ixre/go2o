/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package dao

import (
	"github.com/atnet/gof/app"
	"go2o/core/share/glob"
)

var (
	context  app.Context
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
