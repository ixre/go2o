/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-17 17:03
 * description :
 * history :
 */

package infrastructure

import (
	"com/share/glob"
	"ops/cf/app"
)

var (
	context app.Context
)

func GetContext() app.Context {
	if context == nil {
		context = glob.CurrContext()
	}
	return context
}
