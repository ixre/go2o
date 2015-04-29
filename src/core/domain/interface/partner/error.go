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
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrNoSuchPartner *domain.DomainError = domain.NewDomainError(
		"no_such_partner", "商家不存在")

	ErrNoSuchShop *domain.DomainError = domain.NewDomainError(
		"no_such_shop", "门店不存在")
)
