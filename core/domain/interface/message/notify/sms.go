package notify

/**
 * Copyright 2009-2019 @ 56x.net
 * name : sms.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-01 19:33
 * description :
 * history :
 */

// 短信接口参数设置
type SmsApiPerm struct {
	// 短信提供商,1:通用HTTP接口,2:腾讯云短信,3:阿里云短信,4:创蓝短信
	Provider int
	// 接口KEY
	Key string
	// 接口密钥
	Secret string
	/** 签名 */
	Signature string
	// 默认模板编号
	TemplateId string
	// 是否启用
	Enabled bool
	// Http接口
	Extra *SmsExtraSetting
}

/** 短信接口额外信息配置 */
type SmsExtraSetting struct {
	/** 接口地址 */
	ApiUrl string
	/** 请求数据,如: phone={phone}&content={content}*/
	Params string
	/** 请求方式, GET或POST */
	Method string
	/** 编码 */
	Charset string
	/** 发送成功，包含的字符，用于检测是否发送成功 */
	SuccessChars string
}

// 短信接口设置
type SmsApiSet map[int]*SmsApiPerm
