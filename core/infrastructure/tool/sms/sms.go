/**
 * Copyright 2015 @ z3q.net.
 * name : sms
 * author : jarryliu
 * date : 2016-06-14 09:52
 * description : 接口中的参数均以模板和数据形式出现
 * history :
 */
package sms

import (
	"errors"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/iface/aliyu"
	"strconv"
	"strings"
	"go2o/core/infrastructure/iface/cl253"
)

const (
	SmsAli     = 1 //阿里大鱼
	SmsNetEasy = 2 //网易
	SmsCl253 = 3
)

// 附加检查手机短信的参数
func AppendCheckPhoneParams(provider int, param map[string]interface{}) map[string]interface{} {
	//todo: 考虑在参数中读取
	if provider == SmsAli {
		param[aliyu.ParamKeyTplName] = ""
		param[aliyu.ParamKeyTplId] = ""
	}
	return param
}

// 发送短信
func SendSms(provider int, appKey, appSecret, phoneNum string,
	tpl string, param map[string]interface{}) error {
	switch provider {
	case SmsAli:
		return aliyu.SendSms(appKey, appSecret, phoneNum, tpl, param)
	case SmsCl253:
		return cl253.SendMsgToMobile(appKey,appSecret,phoneNum,compile(tpl,param))
	}
	return errors.New("未知的短信接口服务商" + strconv.Itoa(provider))
}

// 解析模板中的参数
func compile(tpl string, param map[string]interface{}) string {
	var str string
	for k, v := range param {
		switch v.(type) {
		case string:
			str = v.(string)
		case int,int32, int64:
			str = strconv.Itoa(v.(int))
		case float32, float64:
			str = format.FormatFloat(v.(float32))
		case bool:
			str = strconv.FormatBool(v.(bool))
		default:
			str = "unknown"
		}
		tpl = strings.Replace(tpl, "{"+k+"}", str, -1)
	}
	return tpl
}
