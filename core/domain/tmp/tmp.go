/**
 * Copyright 2015 @ 56x.net.
 * name : tmp
 * author : jarryliu
 * date : 2016-05-27 10:42
 * description :
 * history :
 */
package tmp

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

/**  此包用于临时代码

1. 领域中直接操作数据源,后期再重构到repository

*/

var Orm orm.Orm

// 数据库
func Db() db.Connector {
	return gof.CurrentApp.Db()
}

func SetORM(o orm.Orm) {
	Orm = o
}
