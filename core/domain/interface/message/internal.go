/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name : internal.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2024-08-03 00:33
 * description :
 * history :
 */

package mss

var SMS_CHECK_CODE = "SMS_CHECK_CODE"

// InternalSmsTemplate 内置短信模板
var InternalSmsTemplate = []*NotifyTemplate{
	{
		Code:     SMS_CHECK_CODE,
		TempName: "短信验证码",
		// 验证码模板变量只能为0-6位数字
		// 申请模板内容： 您的验证码为{1},有效期{2}分钟，如非本人操作，请忽略本短信！
		Content:  "您的验证码为${验证码},有效期${有效时间}分钟，如非本人操作，请忽略本短信！",
		TempType: 2,
		SpCode:   "TENCENT",
		SpTid:    "",
		Labels:   "验证码;有效时间",
	},
}

// InternalMailTemplate 内置邮件模板
var InternalMailTemplate = []*NotifyTemplate{
	{
		Code:     "EMAIL_CHECK_CODE",
		TempName: "邮件验证码",
	},
}
