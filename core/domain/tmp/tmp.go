/**
 * Copyright 2015 @ to2.net.
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
)

/**  此包用于临时代码

1. 领域中直接操作数据源,后期再重构到repository

*/

// 数据库
func Db() db.Connector {
	return gof.CurrentApp.Db()
}
