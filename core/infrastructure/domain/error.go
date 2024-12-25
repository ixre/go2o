/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package domain

import (
	"github.com/ixre/go2o/core/initial/provide"
)

// todo: 可以做通过后台设置错误信息
// 处理错误
func HandleError(err error, src string) error {
	debug := provide.GetApp().Debug()

	if err != nil && debug {
		logger := provide.GetApp().Log()
		logger.Println("[ GO2O][ ERROR] - ", err.Error())
	}
	return err
}
