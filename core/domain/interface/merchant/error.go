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
	ErrMerchantDisabled *domain.DomainError = domain.NewError(
		"err_merchant_disabled", "商户已被停用")

	ErrEnabledFxSales *domain.DomainError = domain.NewError(
		"err_enabled_fx_sales", "系统当前不允许启用分销")

	ErrMerchantExpires *domain.DomainError = domain.NewError(
		"err_merchant_expires", "商户已过期")

	ErrNoSuchMerchant *domain.DomainError = domain.NewError(
		"no_such_partner", "商户不存在")

	ErrNoSuchShop *domain.DomainError = domain.NewError(
		"no_such_shop", "门店不存在")

	ErrMerchantNotMatch *domain.DomainError = domain.NewError(
		"not_match", "商户不匹配")

	ErrSalesPercent *domain.DomainError = domain.NewError(
		"err_sales_percent", "销售比例错误")

	ErrAmount *domain.DomainError = domain.NewError(
		"err_mch_amount", "金额不正确")

	ErrNoMoreAmount *domain.DomainError = domain.NewError(
		"err_mch_no_more_amount", "余额不足")

	ErrNoSuchSignUpInfo *domain.DomainError = domain.NewError(
		"err_no_such_sign_up_info", "商户申请信息不存在")

	ErrRequireRejectRemark *domain.DomainError = domain.NewError(
		"err_mch_require_remark", "请填写退回的原因")

	ErrMissingCompanyName *domain.DomainError = domain.NewError(
		"err_mch_missing_company_name", "请填写公司名称")

	ErrMissingMerchantName *domain.DomainError = domain.NewError(
		"err_mch_missing_merchant_name", "请填写商户名称")

	ErrMissingCompanyNo *domain.DomainError = domain.NewError(
		"err_mch_missing_company_no", "请填写营业执照编号")

	ErrMissingAddress *domain.DomainError = domain.NewError(
		"err_mch_missing_address", "请填写详细地址")

	ErrMissingPersonName *domain.DomainError = domain.NewError(
		"err_mch_missing_person_name", "请填写法人姓名")

	ErrMissingPhone *domain.DomainError = domain.NewError(
		"err_mch_missing_phone", "请填写联系电话")

	ErrMissingPersonId *domain.DomainError = domain.NewError(
		"err_mch_missing_person_id", "请填写法人身份证")

	ErrPersonCardId *domain.DomainError = domain.NewError(
		"err_mch_missing_person_card_id", "法人身份证号码不正确")

	ErrMissingCompanyImage *domain.DomainError = domain.NewError(
		"err_mch_missing_company_image", "请上传营业执照复印件")

	ErrMissingPersonImage *domain.DomainError = domain.NewError(
		"err_mch_missing_person_image", "请上传法人身份证复印件")
)
