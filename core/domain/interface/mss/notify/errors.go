package notify

import "go2o/core/infrastructure/domain"

/**
 * Copyright 2009-2019 @ to2.net
 * name : errors.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-01 19:27
 * description :
 * history :
 */

var(

	ErrNoSuchNotifyItem = domain.NewError(
		"err_no_such_notify_item", "通知项不存在")
	ErrNotSettingSmsProvider = domain.NewError(
		"err_not_setting_sms_provider", "未配置短信服务商")
	ErrNoSuchSmsProvider = domain.NewError(
		"err_not_such_sms_provider", "不存在短信服务商")

)
