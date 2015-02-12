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
	"github.com/atnet/gof/app"
	"go2o/share/glob"
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
