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

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/aliyu"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/cl253"
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
	// 默认模板编号
	TemplateId string
	//发送成功，包含的字符，用于检测是否发送成功
	SuccessChar string
}

// 发送短信
func SendSms(setting *notify.SmsApiPerm, phoneNum string, content string,
	params []string) error {
	if setting.Signature != "" && !strings.Contains(content, setting.Signature) {
		content = setting.Signature + content
	}
	c := resolveMessage(content, params)
	switch mss.SmsProvider(setting.Provider) {
	case mss.HTTP:
		return sendPhoneMsgByHttpApi(setting, phoneNum, c, params, setting.TemplateId)
	case mss.TECENT_CLOUD:
		return errors.New("not implemented")
	case mss.ALIYUN:
		templateName := ""
		return aliyu.SendSms(setting.Key,
			setting.Secret, phoneNum,
			content, params,
			templateName, setting.TemplateId)
	case mss.CHUANGLAN:
		return cl253.SendMsgToMobile(setting.Key, setting.Secret, phoneNum, c)
	}
	return nil
}

// 检查API接口数据是否正确
func CheckSmsApiPerm(s *notify.SmsApiPerm) error {
	if s.Provider == int(mss.HTTP) {
		if s.Extra.ApiUrl == "" {
			return errors.New("HTTP短信接口必须提供API URL")
		}
		if !strings.Contains(s.Extra.Params, "{key}") {
			return errors.New("API Params缺少\"{key}\"字段")
		}
		if !strings.Contains(s.Extra.Params, "{secret}") {
			return errors.New("API Params缺少\"{secret}\"字段")
		}
		if !strings.Contains(s.Extra.Params, "{phone}") {
			return errors.New("API Params缺少\"{phone}\"字段")
		}
		if !strings.Contains(s.Extra.Params, "{content}") {
			return errors.New("API Params缺少\"{content}\"字段")
		}
		if s.Extra.SuccessChars == "" {
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
