package notify

/**
 * Copyright 2009-2019 @ to2.net
 * name : sms.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-01 19:33
 * description :
 * history :
 */

// 短信接口
type SmsApiPerm struct {
	//接口编号
	ApiKey string
	//接口密钥
	ApiSecret string
	//接口地址
	ApiUrl string
	//发送内容的编码
	Encoding string
	//发送成功，包含的字符，用于检测是否发送成功
	SuccessChar string
	//是否默认的接口使用
	Default bool
}

// 短信接口设置
type SmsApiSet map[int]*SmsApiPerm