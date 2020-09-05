/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package delivery

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNotCoveragedArea *domain.DomainError = domain.NewError(
		"not_coveraged_area", "未覆盖的配送区域")
)
