/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package promotion

import (
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrCanNotApplied *domain.DomainError = domain.NewDomainError(
		"name_exists", "无法应用此优惠")
)
