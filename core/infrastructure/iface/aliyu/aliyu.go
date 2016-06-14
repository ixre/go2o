/**
 * Copyright 2015 @ z3q.net.
 * name : aliyu.go
 * author : jarryliu
 * date : 2016-06-14 09:25
 * description :
 * history :
 */
package aliyu

import (
    "github.com/jsix/alidayu"
    "sync"
    "encoding/json"
    "errors"
)

var (
    mux sync.RWMutex
    ParamKeyTplName string = "ali_template"
    ParamKeyTplId string = "ali_template_id"
)

// 发送短信
func SendSms(appKey, appSecret, phoneNum string,
tpl string, param map[string]interface{}) error {
    mux.Lock()
    defer mux.Unlock()
    alidayu.AppKey = appKey
    alidayu.AppSecret = appSecret
    tplName, ok := param[ParamKeyTplName]
    tplId, ok1 := param[ParamKeyTplId]
    if !ok || !ok1 {
        panic(errors.New("param must contain ali_template and ali_template_id keys."))
    }
    delete(param, ParamKeyTplName)
    delete(param, ParamKeyTplId)

    d, _ := json.Marshal(param)
    success, resp := alidayu.SendSMS(phoneNum,
        tplName.(string),
        tplId.(string),
        string(d),
    )
    if success {
        return nil
    }
    return errors.New(resp)
}
