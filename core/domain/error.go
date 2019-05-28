/**
 * Copyright 2015 @ z3q.net.
 * name : error
 * author : jarryliu
 * date : 2016-02-27 20:03
 * description :
 * history :
 */
package domain

import "go2o/core/infrastructure/domain"

var (
	ErrState *domain.DomainError = domain.NewError(
		"err_state", "state error")

	ErrPwdCannotSame *domain.DomainError = domain.NewError(
		"Err_Pwd_Can_not_Same", "新密码不能与旧密码相同")

	ErrPwdOldPwdNotRight *domain.DomainError = domain.NewError(
		"Err_Pwd_Pld_Pwd_Not_Right", "原密码不正确")
)
