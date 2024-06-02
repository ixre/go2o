/**
 * Copyright 2015 @ 56x.net.
 * name : aliyu.go
 * author : jarryliu
 * date : 2016-06-14 09:25
 * description :
 * history :
 */
package aliyu

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/ixre/alidayu"
)

var (
	mux             sync.RWMutex
	ParamKeyTplName = "ali_template"
	ParamKeyTplId   = "ali_template_id"
)

// 发送短信
func SendSms(appKey, appSecret, phoneNum string,
	tpl string, param []string,
	templateName string,
	templateId string) error {
	mux.Lock()
	defer mux.Unlock()
	alidayu.AppKey = appKey
	alidayu.AppSecret = appSecret
	if len(templateName) == 0 || len(templateId) == 0 {
		return errors.New(`
		param must contain "ali_template"
		 and "ali_template_id" keys.`)
	}
	d, _ := json.Marshal(param)
	success, resp := alidayu.SendSMS(phoneNum,
		templateName,
		templateId,
		string(d),
	)
	if success {
		return nil
	}
	return errors.New(resp)
}
