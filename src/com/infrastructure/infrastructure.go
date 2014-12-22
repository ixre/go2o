/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-10 16:23
 * description :
 * history :
 */

package infrastructure

import (
	"com/share/glob"
	"ops/cf/app"
)

var (
	DebugMode bool
	Context   app.Context
)

func init() {
	Context = glob.CurrContext()
}
