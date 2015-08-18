/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 20:12
 * description :
 * history :
 */
package member

import (
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrSessionTimeout *domain.DomainError = domain.NewDomainError(
		"member_session_time_out", "会员会话超时")

	ErrInvalidSession *domain.DomainError = domain.NewDomainError(
		"member_invalid_session", "异常会话")

	ErrNoSuchDeliverAddress *domain.DomainError = domain.NewDomainError(
		"member_no_such_deliver_address", "配送地址错误")

	ErrNoSuchMember *domain.DomainError = domain.NewDomainError(
		"member_no_such_member", "会员不存在")

	ErrDeliverAddressLen *domain.DomainError = domain.NewDomainError(
		"err_deliver_address_len", "请填写详细的配送地址")

	ErrDeliverRealNameIsNull *domain.DomainError = domain.NewDomainError(
		"err_deliver_real_name_is_null", "收货人不能为空")

	ErrDeliverPhoneIsNull *domain.DomainError = domain.NewDomainError(
		"err_deliver_phone_is_null", "电话不能为空")

	ErrPwdCannotSame *domain.DomainError = domain.NewDomainError(
		"Err_Pwd_Can_not_Same", "新密码不能与旧密码相同")

	ErrPwdOldPwdNotRight *domain.DomainError = domain.NewDomainError(
		"Err_Pwd_Pld_Pwd_Not_Right", "原密码不正确")

	ErrIncorrectAmount *domain.DomainError = domain.NewDomainError(
		"err_balance_amount", "金额错误")

	ErrOutOfBalance *domain.DomainError = domain.NewDomainError(
		"err_out_of_balance", "超出金额")

	ErrUserLength *domain.DomainError = domain.NewDomainError(
		"err_user_length", "用户名必须大于6位",
	)

	ErrUserValidErr *domain.DomainError = domain.NewDomainError(
		"err_user_valid_err", "用户名为6位以上字符和数字的组合")

	ErrEmailValidErr *domain.DomainError = domain.NewDomainError(
		"err_email_valid_err", "邮箱不正确")

	ErrPhoneValidErr *domain.DomainError = domain.NewDomainError(
		"err_phone_valid_err", "手机号码不正确")

	ErrQqValidErr *domain.DomainError = domain.NewDomainError(
		"err_qq_valid_err", "QQ号码不正确")

	ErrNotSetTradePwd *domain.DomainError = domain.NewDomainError(
		"err_not_set_tarde_pwd", "交易密码未设置")

	ErrIncorrectTradePwd *domain.DomainError = domain.NewDomainError(
		"err_incorrect_tarde_pwd", "交易密码错误")
)
