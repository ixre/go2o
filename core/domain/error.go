/**
 * Copyright 2015 @ 56x.net.
 * name : error
 * author : jarryliu
 * date : 2016-02-27 20:03
 * description :
 * history :
 */
package domain

import "github.com/ixre/go2o/core/infrastructure/domain"

var (
	ErrState = domain.NewError(
		"err_state", "state error")

	ErrPwdCannotSame = domain.NewError(
		"Err_Pwd_Can_not_Same", "新密码不能与旧密码相同")

	ErrPwdOldPwdNotRight = domain.NewError(
		"Err_Pwd_Pld_Pwd_Not_Right", "原密码不正确")
)
