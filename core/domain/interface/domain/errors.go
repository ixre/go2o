package domain

import "go2o/core/infrastructure/domain"

/**
 * Copyright 2009-2019 @ to2.net
 * name : errors.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-11 17:23
 * description :
 * history :
 */

var (
	ErrNotMD5Format = domain.NewError("err_not_md5_format", "密码非32位(MD5)")

	ErrCredential = domain.NewError(
		"err_credential", "用户或密码不正确")

	ErrPwdLength = domain.NewError(
		"err_pwd_length", "密码至少包含6个字符")

	ErrCheckCodeError = domain.NewError(
		"err_check_code_err", "验证码不正确")

	ErrCheckCodeExpires = domain.NewError(
		"err_check_code_expires", "验证码已过期")
)
