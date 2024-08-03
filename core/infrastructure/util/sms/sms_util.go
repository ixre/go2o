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
	"strings"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/go2o/core/infrastructure/util/sms/aliyu"
	"github.com/ixre/go2o/core/infrastructure/util/sms/cl253"
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
func SendSms(setting *mss.SmsApiPerm, phoneNum string, content string,
	params []string) error {
	if setting.Signature != "" && !strings.Contains(content, setting.Signature) {
		content = setting.Signature + content
	}
	c := util.ResolveMessage(content, params)
	switch mss.SmsProvider(setting.Provider) {
	case mss.HTTP:
		return sendPhoneMsgByHttpApi(setting, phoneNum, c, params, "")
	case mss.TECENT_CLOUD:
		return errors.New("not implemented")
	case mss.ALIYUN:
		templateName := ""
		return aliyu.SendSms(setting.Key,
			setting.Secret, phoneNum,
			content, params,
			templateName, "")
	case mss.CHUANGLAN:
		return cl253.SendMsgToMobile(setting.Key, setting.Secret, phoneNum, c)
	}
	return nil
}

// 检查API接口数据是否正确
func CheckSmsApiPerm(s *mss.SmsApiPerm) error {
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
