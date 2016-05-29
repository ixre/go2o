/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package merchant

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrMerchantDisabled *domain.DomainError = domain.NewDomainError(
		"err_merchant_disabled", "商户权限已被取消")

	ErrEnabledFxSales *domain.DomainError = domain.NewDomainError(
		"err_enabled_fx_sales", "系统当前不允许启用分销")

	ErrMerchantExpires *domain.DomainError = domain.NewDomainError(
		"err_merchant_expires", "商户合作已过期")

	ErrNoSuchMerchant *domain.DomainError = domain.NewDomainError(
		"no_such_partner", "商户不存在")

	ErrNoSuchShop *domain.DomainError = domain.NewDomainError(
		"no_such_shop", "门店不存在")

	ErrMerchantNotMatch *domain.DomainError = domain.NewDomainError(
		"not_match", "商户不匹配")

	ErrSalesPercent *domain.DomainError = domain.NewDomainError(
		"err_sales_percent", "销售比例错误")
)
