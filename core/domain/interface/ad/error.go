/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package ad

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
)

var (
	ErrInternalDisallow = domain.NewError(
		"err_internal_disallow", "不允许修改系统内置广告")

	ErrNoSuchAd = domain.NewError(
		"err_no_such_ad", "广告不存在")

	ErrNoSuchAdGroup = domain.NewError(
		"err_no_such_ad_group", "广告位分组不存在")

	ErrNotEmptyGroup = domain.NewError(
		"err_ad_not_empty_ad_group", "广告位分组包含广告位，无法删除")

	ErrNoSuchAdPosition = domain.NewError(
		"err_no_such_ad_position", "广告位不存在")

	ErrNotOpened = domain.NewError(
		"err_position_not_opened", "广告未开放")
	ErrUserPositionIsBind = domain.NewError(
		"err_ad_user_position_is_bind", "该广告位已绑定其他广告")
	ErrAdType = domain.NewError(
		"err_ad_type", "请选择广告类型")

	ErrDisallowModifyAdType = domain.NewError(
		"err_disallow_modify_ad_type", "广告创建后不允许修改类型")

	ErrKeyExists = domain.NewError(
		"err_key_exists", "广告位KEY已存在")

	ErrAdUsed = domain.NewError(
		"err_ad_used", " 无法删除已投放的广告")
)
