/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package merchant

import (
	"errors"

	"github.com/ixre/go2o/core/infrastructure/domain"
)

var (
	ErrMerchantDisabled = domain.NewError(
		"err_merchant_disabled", "商户已被停用")

	ErrEnabledFxSales = domain.NewError(
		"err_enabled_fx_sales", "系统当前不允许启用分销")

	ErrMerchantExpires = domain.NewError(
		"err_merchant_expires", "商户已过期")

	ErrMissingMerchantUser = domain.NewError(
		"err_missing_merchant_user", "未设置商户用户")

	ErrBindAnotherMerchant = domain.NewError(
		"err_bind_another_merchant", "会员已绑定其他商户")

	ErrMemberBindExists = domain.NewError(
		"err_member_bind_exist", "当前已绑定该会员")

	ErrMerchantUserExists = domain.NewError(
		"err_merchant_user_exists", "商户已存在")

	ErrNoSuchMerchant = domain.NewError(
		"no_such_partner", "商户不存在")

	ErrNoSuchShop = domain.NewError(
		"no_such_shop", "门店不存在")

	ErrMerchantNotMatch = domain.NewError(
		"not_match", "商户不匹配")

	ErrSalesPercent = domain.NewError(
		"err_sales_percent", "销售比例错误")

	ErrTxRate = domain.NewError(
		"err_tx_rate", "交易比例错误")

	ErrAmount = domain.NewError(
		"err_mch_amount", "金额不正确")

	ErrNoMoreAmount = domain.NewError(
		"err_mch_no_more_amount", "余额不足")

	ErrNoSuchSignUpInfo = domain.NewError(
		"err_no_such_sign_up_info", "商户申请信息不存在")

	ErrRequireRejectRemark = domain.NewError(
		"err_mch_require_remark", "请填写退回的原因")

	ErrMissingCompanyName = domain.NewError(
		"err_mch_missing_company_name", "请填写公司名称")

	ErrMissingMerchantName = domain.NewError(
		"err_mch_missing_merchant_name", "请填写商户名称")

	ErrMissingCompanyNo = domain.NewError(
		"err_mch_missing_company_no", "请填写营业执照编号")

	ErrMissingAddress = domain.NewError(
		"err_mch_missing_address", "请填写详细地址")

	ErrMissingPersonName = domain.NewError(
		"err_mch_missing_person_name", "请填写法人姓名")

	ErrMissingPhone = domain.NewError(
		"err_mch_missing_phone", "请填写联系电话")

	ErrMissingPersonId = domain.NewError(
		"err_mch_missing_person_id", "请填写法人身份证")

	ErrPersonCardId = domain.NewError(
		"err_mch_missing_person_card_id", "法人身份证号码不正确")

	ErrMissingCompanyImage = domain.NewError(
		"err_mch_missing_company_image", "请上传营业执照复印件")

	ErrMissingPersonImage = domain.NewError(
		"err_mch_missing_person_image", "请上传法人身份证复印件")
)

// CheckMchStatus 检查商户状态
func CheckMchStatus(status int) error {
	switch status {
	case 0:
		return errors.New("商户未完成审核")
	case 2:
		return errors.New("商户已停用")
	case 3:
		return errors.New("商户已关闭")
	}
	return nil
}
