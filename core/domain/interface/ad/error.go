/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package ad

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrInternalDisallow *domain.DomainError = domain.NewDomainError(
		"err_internal_disallow", "不允许修改系统内置广告")

	ErrNoSuchAd *domain.DomainError = domain.NewDomainError(
		"err_no_such_ad", "广告不存在")

	ErrNoSuchAdGroup *domain.DomainError = domain.NewDomainError(
		"err_no_such_ad_group", "广告组不存在")

	ErrNoSuchAdPosition *domain.DomainError = domain.NewDomainError(
		"err_no_such_ad_position", "广告位不存在")

	ErrDisallowModifyAdType *domain.DomainError = domain.NewDomainError(
		"err_disallow_modify_ad_type", "广告创建后不允许修改类型")

	ErrKeyExists *domain.DomainError = domain.NewDomainError(
		"err_key_exists", "广告位KEY已存在")

	ErrAdUsed *domain.DomainError = domain.NewDomainError(
		"err_ad_used", " 无法删除已投放的广告")
)
