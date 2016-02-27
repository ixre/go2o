/**
 * Copyright 2015 @ z3q.net.
 * name : error
 * author : jarryliu
 * date : 2016-02-27 20:03
 * description :
 * history :
 */
package domain

import "go2o/src/core/infrastructure/domain"

var (
	ErrState *domain.DomainError = domain.NewDomainError(
		"err_state", "state error")
)
