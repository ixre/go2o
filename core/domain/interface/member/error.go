/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 20:12
 * description :
 * history :
 */
package member

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrIncorrectInfo = domain.NewError(
		"err_incorrect_info", "非法数据")

	ErrRegOff = domain.NewError(
		"err_reg_off", "CODE:1010,系统暂停注册")

	ErrRegMissingInvitationCode = domain.NewError(
		"err_reg_must_invitation", "CODE:1011,请填写邀请码")

	ErrRegOffInvitation = domain.NewError(
		"err_reg_off_invitation", "CODE:1012,系统关闭邀请注册")

	ErrInvitationCode = domain.NewError(
		"member_err_invation_code", "CODE:1013,邀请码错误")

	ErrSessionTimeout = domain.NewError(
		"member_session_time_out", "会员会话超时")

	ErrMemberLocked = domain.NewError(
		"err_member_locked", "会员已被停用")

	ErrMemberHasActive = domain.NewError(
		"err_member_has_active", "会员已激活")

	ErrPremiumValue = domain.NewError(
		"err_member_premium_value", "premium not match")

	ErrCredential = domain.NewError(
		"err_member_credential", "会员用户或密码不正确")

	ErrCheckCodeError = domain.NewError(
		"err_member_check_code_err", "验证码不正确")

	ErrCheckCodeExpires = domain.NewError(
		"err_member_check_code_expires", "验证码已过期")

	ErrInvalidSession = domain.NewError(
		"member_invalid_session", "异常会话")

	ErrPwdLength = domain.NewError(
		"err_member_pwd_length", "密码至少包含6个字符")

	ErrNoSuchAddress = domain.NewError(
		"member_no_such_deliver_address", "配送地址错误")

	ErrLevelUsed = domain.NewError(
		"member_level_used", "此等级已被会员使用")

	ErrMaxLevel = domain.NewError(
		"err_member_max_level", "已经为最高等级")

	ErrLevelUpPass = domain.NewError(
		"err_member_level_up_pass", "会员升级已经审核")

	ErrLevelUpReject = domain.NewError(
		"err_member_level_up_reject", "会员升级已经被退回")

	ErrLevelUpConfirm = domain.NewError(
		"err_member_level_up_confirm", "会员升级已经确认")
	ErrLevelUpLaterConfirm = domain.NewError(
		"err_member_level_up_later_confirm", "请稍后二分钟后确认")

	ErrNoSuchLevelUpLog = domain.NewError(
		"err_member_no_such_level_up_log", "不存在升级信息")

	ErrLevelDisabled = domain.NewError(
		"member_level_disabled", "等级未启用")

	ErrExistsSameProgramSignalLevel = domain.NewError(
		"member_exists_same_program_signal_level", "存在相同可编程签名的等级")

	ErrMustMoreThanMaxLevel = domain.NewError(
		"member_level_must_more_than_max_level", "经验值必须大于最大等级")

	ErrLessThanLevelRequireExp = domain.NewError(
		"member_level_less_than_exp", "经验值必须大于前一等级")

	ErrMoreThanLevelRequireExp = domain.NewError(
		"member_level_more_than_exp", "经验值必须小于后一等级")

	ErrNoSuchMember = domain.NewError(
		"member_no_such_member", "会员不存在")

	ErrNoValidInviter = domain.NewError(
		"member_no_valid_inviter", "邀请人不正确")
	ErrDeliverAddressLen = domain.NewError(
		"err_deliver_address_len", "请填写详细的配送地址")

	ErrDeliverContactPersonName = domain.NewError(
		"err_deliver_contact_person_name", "收货人不正确")

	ErrDeliverContactPhone = domain.NewError(
		"err_deliver_phone_is_null", "联系人电话有误")

	ErrNotSetArea = domain.NewError(
		"err_not_set_area", "地址不正确")

	ErrNoSuchBankInfo = domain.NewError(
		"err_no_such_bank_info", "请完善银行卡信息")

	ErrBankInfo = domain.NewError(
		"err_member_bank_info", "银行卡信息不正确")

	ErrBankName = domain.NewError(
		"err_member_bank_name", "请选择开户银行")

	ErrBankAccountName = domain.NewError(
		"err_member_bank_account_name", "开户银行户名不正确")

	ErrBankAccount = domain.NewError(
		"err_member_bank_account", "开户银行账号不正确")

	ErrBankNetwork = domain.NewError(
		"err_member_bank_network", "开户银行支行不正确")

	ErrBankInfoLocked = domain.NewError(
		"err_bank_info_locked", "银行卡信息已锁定,无法更改")

	ErrBankInfoNoYetSet = domain.NewError(
		"err_bank_info_no_yet_set", "银行卡信息尚未设置")

	ErrIncorrectAmount = domain.NewError(
		"err_balance_amount", "金额错误")

	ErrLessTakeAmount = domain.NewError(
		"err_account_less_take_amount", "单笔最低提现金额为%s")

	ErrOutTakeAmount = domain.NewError(
		"err_account_out_take_amount", "单笔最高提现金额为%s")

	ErrTakeOutLevelNoPerm = domain.NewError(
		"err_account_take_out_level_no_perm", "%s会员无法提现")

	ErrTakeOutNotTrust = domain.NewError(
		"err_account_take_out_not_trust", "必须通过实名认证后才可提现")

	ErrIncorrectQuota = domain.NewError(
		"err_member_incorrent_quote", "金额/数量错误")

	ErrOutOfBalance = domain.NewError(
		"err_out_of_balance", "超出金额")

	ErrAccountOutOfTakeOutTimes = domain.NewError(
		"err_account_out_of_take_out_times", "今日已达到提现次数上限")

	ErrTransferAccountsLevelNoPerm = domain.NewError(
		"err_account_transfer_accounts_level_no_perm", "%s不允许转账")

	ErrUsrLength = domain.NewError(
		"err_user_length", "用户名至少6位",
	)

	ErrUsrValidErr = domain.NewError(
		"err_user_valid_err", "用户名为6位以上字符和数字的组合")

	ErrSameUsr = domain.NewError(
		"err_same_user", "用户名与原来相同")

	ErrUsrExist = domain.NewError(
		"err_member_user_exist", "用户名已存在")

	ErrNilNickName = domain.NewError(
		"err_member_nil_nick_name", "昵称不能为空")

	ErrNullAvatar = domain.NewError(
		"err_member_null_avatar", "请上传头像")

	ErrAddress = domain.NewError(
		"err_member_address", "请填写详细地址")

	ErrEmailValidErr = domain.NewError(
		"err_member_email_valid_err", "邮箱不正确")

	ErrPhoneValidErr = domain.NewError(
		"err_member_phone_valid_err", "手机号码不正确")

	ErrPhoneHasBind = domain.NewError(
		"err_member_phone_has_bind", "手机号码已经使用")

	ErrMissingPhone = domain.NewError(
		"err_member_missing_phone", "请填写手机号码")

	ErrMissingIM = domain.NewError(
		"err_member_missing_im", "请填写IM")

	ErrBadPhoneFormat = domain.NewError(
		"err_bad_phone_format", "手机号码不正确")

	ErrQqValidErr = domain.NewError(
		"err_qq_valid_err", "QQ号码不正确")

	ErrNotSetTradePwd = domain.NewError("err_not_set_trade_pwd", "交易密码未设置")

	ErrIncorrectTradePwd            = domain.NewError("err_incorrect_trade_pwd", "交易密码错误")
	ErrNoSuchLogTitleOrRemark       = domain.NewError("err_member_no_such_log_title_or_remark", "缺少账户变动的标题和备注")
	ErrNotSupportAccountType        = domain.NewError("err_account_not_support", "账户类型不支持此操作")
	ErrNotSupportPaymentAccountType = domain.NewError("err_account_not_support_payment", "不支持支付的账户类型")

	ErrAccountNotEnoughAmount = domain.NewError("err_not_enough_amount", "账户余额不足")

	ErrTakeOutState = domain.NewError("err_member_take_out_state", "提现申请状态错误")

	ErrNotSupportTakeOutBusinessKind = domain.NewError("err_not_support_take_out_business_kind", "不支持的提现业务类型")

	ErrBusinessKind = domain.NewError(
		"err_not_support_business_kind", "不支持的业务类型")

	ErrNotSupportChargeMethod = domain.NewError(
		"err_account_not_support_charge_method", "不支持的充值方式")

	ErrNotSupportTransfer = domain.NewError(
		"err_account_not_support_transfer", "不支持的转账方式")

	ErrNoSuchRelateUser = domain.NewError(
		"err_account_no_such_relate_user", "未提供操作人编号")

	ErrMissingTrustedInfo = domain.NewError(
		"err_missing_trusted_info", "信息不完整、无法完成实名认证")

	ErrEmptyReviewRemark = domain.NewError(
		"err_member_empty_remark", "原因不能为空")

	ErrNotTrusted = domain.NewError(
		"err_member_not_trusted", "尚未实名认证")

	ErrNoChangedTrustInfo = domain.NewError(
		"err_member_no_changed_trust_info", "请修改实名信息后再进行提交")

	ErrRealName = domain.NewError(
		"err_real_name", "请输入真实姓名")

	ErrTrustCardId = domain.NewError(
		"err_member_trust_car_id", "身份证号码不正确")

	ErrCarIdExists = domain.NewError(
		"err_member_trust_car_id", "身份证号码已使用")

	ErrTrustMissingImage = domain.NewError(
		"err_member_trust_missing_image", "请上传认证照片")

	ErrFavored = domain.NewError(
		"err_favored", "已经收藏过了")

	ErrAccountBalanceNotEnough = domain.NewError(
		"err_account_balance_not_enough ", "账户余额不足")

	ErrNoSuchIntegral = domain.NewError(
		"err_account_no_such_integral", "账户积分不足")

	ErrMissingOuterNo = domain.NewError(
		"err_account_missing_outer_no", "缺少订单号")
)
