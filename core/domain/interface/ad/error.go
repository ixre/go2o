/**
 * Copyright 2014 @ to2.net.
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
	ErrInternalDisallow *domain.DomainError = domain.NewError(
		"err_internal_disallow", "不允许修改系统内置广告")

	ErrNoSuchAd *domain.DomainError = domain.NewError(
		"err_no_such_ad", "广告不存在")

	ErrNoSuchAdGroup *domain.DomainError = domain.NewError(
		"err_no_such_ad_group", "广告组不存在")

	ErrNotEmptyGroup *domain.DomainError = domain.NewError(
		"err_ad_not_empty_ad_group", "广告组包含广告位，无法删除")

	ErrNoSuchAdPosition *domain.DomainError = domain.NewError(
		"err_no_such_ad_position", "广告位不存在")

	ErrNotOpened *domain.DomainError = domain.NewError(
		"err_position_not_opened", "广告未开放")
	ErrUserPositionIsBind *domain.DomainError = domain.NewError(
		"err_ad_user_position_is_bind", "该广告位已绑定其他广告")
	ErrAdType *domain.DomainError = domain.NewError(
		"err_ad_type", "请选择广告类型")

	ErrDisallowModifyAdType *domain.DomainError = domain.NewError(
		"err_disallow_modify_ad_type", "广告创建后不允许修改类型")

	ErrKeyExists *domain.DomainError = domain.NewError(
		"err_key_exists", "广告位KEY已存在")

	ErrAdUsed *domain.DomainError = domain.NewError(
		"err_ad_used", " 无法删除已投放的广告")
)
