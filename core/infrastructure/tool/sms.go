/**
 * Copyright 2015 @ z3q.net.
 * name : sms
 * author : jarryliu
 * date : 2016-06-14 09:52
 * description : 接口中的参数均以模板和数据形式出现
 * history :
 */
package tool

import (
    "go2o/core/infrastructure/iface/aliyu"
    "errors"
    "strconv"
    "strings"
    "go2o/core/infrastructure/format"
)

const (
    SmsAli = 1 //阿里大鱼
    SmsNetEasy = 2 //网易
)

func SendSms(provider int,appKey, appSecret, phoneNum string,
    tpl string, param map[string]interface{}) error {
    switch provider {
    case SmsAli:
       return aliyu.SendSms(appKey,appSecret,phoneNum,tpl,param)
    }
    return errors.New("未知的短信接口服务商"+strconv.Itoa(provider))
}

// 解析模板中的参数
func compile(tpl string,param map[string]interface{})string{
    var str string
    for k, v := range param {
        switch v.(type){
        case string:
            str = v.(string)
        case int32,int64:
            str = strconv.Itoa(v.(int))
        case float32,float64:
            str = format.FormatFloat(v.(float32))
        case bool:
            str = strconv.FormatBool(v.(bool))
        default:
            str = "unknown"
        }
        tpl = strings.Replace(tpl, "{" + k + "}", str, -1)
    }
    return tpl
}
