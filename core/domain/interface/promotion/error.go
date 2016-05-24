/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package promotion

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrCanNotApplied *domain.DomainError = domain.NewDomainError(
		"name_exists", "无法应用此优惠")

	ErrExistsSamePromotionFlag *domain.DomainError = domain.NewDomainError(
		"exists_same_promotion_flag", "已存在相同的促销")

	ErrNoSuchPromotion *domain.DomainError = domain.NewDomainError(
		"no_such_promotion", "促销不存在")

	ErrNoDetailsPromotion *domain.DomainError = domain.NewDomainError(
		"no_details_promotion", "促销信息不完整")
)
