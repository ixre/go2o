/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 20:12
 * description :
 * history :
 */
package member

import (
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrSessionTimeout *domain.DomainError = domain.NewDomainError(
		"member_session_time_out", "会员会话超时")

	ErrInvalidSession *domain.DomainError = domain.NewDomainError(
		"member_invalid_session", "异常会话")

	ErrNoSuchDeliverAddress *domain.DomainError = domain.NewDomainError(
		"member_no_such_deliver_address", "配送地址错误")
)
