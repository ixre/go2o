/**
 * Copyright 2015 @ 56x.net.
 * name : sms
 * author : jarryliu
 * date : 2016-06-14 09:52
 * description : 接口中的参数均以模板和数据形式出现
 * history :
 */
package sms

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/aliyu"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/cl253"
)

const (
	SmsHttp  = 1 // HTTP
	SmsAli   = 2 //阿里大鱼
	SmsCl253 = 3 //创蓝253
)

// 短信接口
type SmsApi struct {
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

// 发送短信
func SendSms(provider string, api *SmsApi, phoneNum string, content string,
	params []string) error {
	if api.Signature != "" && !strings.Contains(content, api.Signature) {
		content = api.Signature + content
	}
	c := resolveMessage(content, params)
	templateId := ""
	switch getProviderID(provider) {
	case SmsHttp:
		return sendPhoneMsgByHttpApi(api, phoneNum, c, params, templateId)
	case SmsAli:
		templateName := ""
		return aliyu.SendSms(api.Key,
			api.Secret, phoneNum,
			content, params,
			templateName, templateId)
	case SmsCl253:
		return cl253.SendMsgToMobile(api.Key, api.Secret, phoneNum, c)
	}
	return errors.New("未知的短信接口服务商:" + provider)
}

func getProviderID(provider string) int {
	switch provider {
	case "http":
		return SmsHttp
	case "253":
		return SmsCl253
	case "ali":
		return SmsAli
	}
	return -1
}

// 检查API接口数据是否正确
func CheckSmsApiPerm(provider string, s *notify.SmsApiPerm) error {
	id := getProviderID(provider)
	if id == SmsHttp {
		if s.ApiUrl == "" {
			return errors.New("HTTP短信接口必须提供API URL")
		}
		if !strings.Contains(s.Params, "{key}") {
			return errors.New("API Params缺少\"{key}\"字段")
		}
		if !strings.Contains(s.Params, "{secret}") {
			return errors.New("API Params缺少\"{secret}\"字段")
		}
		if !strings.Contains(s.Params, "{phone}") {
			return errors.New("API Params缺少\"{phone}\"字段")
		}
		if !strings.Contains(s.Params, "{content}") {
			return errors.New("API Params缺少\"{content}\"字段")
		}
		if s.SuccessChar == "" {
			return errors.New("未指定发送成功包含的字符")
		}
	}
	return nil
}

// 解析模板中的参数
func resolveMessage(tpl string, param []string) string {
	//　替换字符标签{name}标签为{0}
	re := regexp.MustCompile(`{(.+?)}`)
	holders := re.FindAllString(tpl, -1)
	for i, v := range holders {
		if _, err := strconv.Atoi(v); err != nil {
			tpl = strings.ReplaceAll(tpl, v, fmt.Sprintf("{%d}", i))
		}
	}
	// 替换值
	for k, v := range param {
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{%d}", k), v)
	}
	return tpl
}
