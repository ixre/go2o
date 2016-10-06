/**
 * Copyright 2015 @ z3q.net.
 * name : tmp
 * author : jarryliu
 * date : 2016-05-27 10:42
 * description :
 * history :
 */
package tmp

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
)

/**  此包用于临时代码

1. 领域中直接操作数据源,后期再重构到repository

*/

// 数据库
func Db() db.Connector {
	return gof.CurrentApp.Db()
}
