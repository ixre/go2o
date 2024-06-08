package mss

var SMS_CHECK_CODE = "SMS_CHECK_CODE"

// InternalSmsTemplate 内置短信模板
var InternalSmsTemplate = []*NotifyTemplate{
	{
		Code:     SMS_CHECK_CODE,
		TempName: "短信验证码",
		Content:  "您好,本次${操作}验证码为${验证码},有效期为${有效时间}分钟。",
		TempType: 2,
		SpCode:   "UCLOUD",
		SpTid:    "",
		Labels:   "操作;验证码;有效时间",
	},
}
