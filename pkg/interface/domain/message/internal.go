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
var EMAIL_MCH_REGISTER = "EMAIL_MCH_REGISTER"

// InternalSmsTemplate 内置短信模板
var InternalSmsTemplate = []*NotifyTemplate{
	{
		TplCode: SMS_CHECK_CODE,
		TplName: "短信验证码",
		TplFlag: TplFlagSystem,
		// 验证码模板变量只能为0-6位数字
		// 申请模板内容： 您的验证码为{1},有效期{2}分钟，如非本人操作，请忽略本短信！
		Content: "您的验证码为${验证码},有效期${有效时间}分钟，如非本人操作，请忽略本短信！",
		TplType: TypeSMS,
		SpCode:  "TENCENT",
		SpTid:   "",
		Labels:  "验证码;有效时间",
	},
}

// InternalMailTemplate 内置邮件模板
var InternalMailTemplate = []*NotifyTemplate{
	{
		TplCode: EMAIL_MCH_REGISTER,
		TplName: "商户注册验证",
		TplFlag: TplFlagSystem,
		TplType: TypeEmail,
		Content: `您正在注册成为商户，请点击以下链接完成注册： <br /><a href="${注册链接}">${注册链接}</a>
		 			<br />
		 		此链接有效期为${有效时间}分钟`,
		SpCode: "",
		SpTid:  "",
		Labels: "注册链接;有效时间",
	},
}
