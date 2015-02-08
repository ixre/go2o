/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package partner

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchPartner *domain.DomainError = domain.NewDomainError(
		"no_such_partner", "商家不存在")
)
