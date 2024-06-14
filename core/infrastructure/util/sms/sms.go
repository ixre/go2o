package sms

import (
	"errors"
	"fmt"

	"github.com/ixre/go2o/core/infrastructure/util/collections"
)

// 短信服务商
var smsProviders []ISmsProvider

// 短信服务商
type ISmsProvider interface {
	// 提供商的名称
	Name() string

	// 发送短信
	Send(templateId string, content string, args ...string) error
}

// 短信模板
type Template struct {
	// 服务商代码
	ProviderCode string
	// 模板内容,如果传递了TemplateId,使用短信服务商申请的模板内容发送
	TemplateContent string
	// 模板ID,如果自定义短信，则不用传ID
	TemplateId string
}

// 注册短信服务商
func RegisterProvider(p ISmsProvider) {
	if collections.AnyArray(smsProviders, func(e ISmsProvider) bool {
		return e.Name() == p.Name()
	}) {
		panic(fmt.Errorf("providerSms provider %s already registered", p.Name()))
	}
	smsProviders = append(smsProviders, p)
}

func getProvider(providerName string) (ISmsProvider, error) {
	if len(smsProviders) == 0 {
		return nil, errors.New("no any provider registered")
	}
	p := collections.FindArray(smsProviders, func(e ISmsProvider) bool {
		return e.Name() == providerName
	})
	if p == nil {
		return nil, fmt.Errorf("未注册短信服务商:%s", providerName)
	}
	return p, nil
}

// 发送短信
func Send(t Template, phoneNum string, params ...string) error {
	if len(t.ProviderCode) == 0 {
		return errors.New("未指定短信服务商或模板ID")
	}
	if len(t.TemplateContent) == 0 || len(t.TemplateId) == 0 {
		return errors.New("未指定短信内容或短信服务商模板ID")
	}
	p, err := getProvider(t.ProviderCode)
	if err != nil {
		return err
	}
	c := ResolveMessage(t.TemplateContent, params)
	return p.Send(t.TemplateId, c, params...)
	// todo: 旧的短信需要重新实现
	// if setting.Signature != "" && !strings.Contains(content, setting.Signature) {
	// 	content = setting.Signature + content
	// }
	// c := ResolveMessage(content, params)
	// switch mss.SmsProvider(setting.Provider) {
	// case mss.HTTP:
	// 	return sendPhoneMsgByHttpApi(setting, phoneNum, c, params, "")
	// case mss.TECENT_CLOUD:
	// 	return errors.New("not implemented")
	// case mss.ALIYUN:
	// 	templateName := ""
	// 	return aliyu.SendSms(setting.Key,
	// 		setting.Secret, phoneNum,
	// 		content, params,
	// 		templateName, "")
	// case mss.CHUANGLAN:
	// 	return cl253.SendMsgToMobile(setting.Key, setting.Secret, phoneNum, c)
	// }
	// return nil
	return nil
}
