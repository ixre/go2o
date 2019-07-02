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
	//接口地址
	ApiUrl string
	//接口编号
	Key string
	//接口密钥
	Secret string
	// 请求数据,如: phone={phone}&content={content}
	Params string
	// 请求方式, GET或POST
	Method string
	//发送内容的编码
	Charset string
	// 签名
	Signature string
	//发送成功，包含的字符，用于检测是否发送成功
	SuccessChar string
}

// 短信接口设置
type SmsApiSet map[int]*SmsApiPerm
