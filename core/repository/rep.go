/**
 * Copyright 2015 @ z3q.net.
 * name : rep
 * author : jarryliu
 * date : 2016-05-24 10:14
 * description :
 * history :
 */
package repository

import "github.com/jsix/gof"

// 处理错误
func handleError(err error) error {
	if err != nil && gof.CurrentApp.Debug() {
		gof.CurrentApp.Log().Println("[ Go2o][ Rep][ Error] -", err.Error())
	}
	return err
}
