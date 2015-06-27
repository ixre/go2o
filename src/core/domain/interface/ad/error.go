/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package ad

import (
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrNameExists *domain.DomainError = domain.NewDomainError(
		"name_exists", "已经存在相同的名称的广告")

	ErrInternalDisallow *domain.DomainError = domain.NewDomainError(
		"err_internal_disallow", "不允许修改系统内置广告")
)
