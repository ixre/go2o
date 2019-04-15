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
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/iface/aliyu"
	"go2o/core/infrastructure/iface/cl253"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	SmsHttp  = 1 // HTTP
	SmsAli   = 2 //阿里大鱼
	SmsCl253 = 3 //创蓝253
)

// 发送短信,tpl:短信内容模板
func SendSms(provider int, appKey, appSecret, phoneNum string,
	apiUrl string, enc string, successChar string, tpl string,
	param map[string]interface{}) error {
	switch provider {
	case SmsHttp:
		return sendPhoneMsgByHttpApi(apiUrl, appKey, appSecret, phoneNum,
			compile(tpl, param), enc, successChar)
	case SmsAli:
		return aliyu.SendSms(appKey, appSecret, phoneNum, tpl, param)
	case SmsCl253:
		return cl253.SendMsgToMobile(appKey, appSecret, phoneNum, compile(tpl, param))
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
		case int, int32, int64:
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

// 附加检查手机短信的参数
func AppendCheckPhoneParams(provider int, param map[string]interface{}) map[string]interface{} {
	//todo: 考虑在参数中读取
	if provider == SmsAli {
		param[aliyu.ParamKeyTplName] = ""
		param[aliyu.ParamKeyTplId] = ""
	}
	return param
}

// 检查API接口数据是否正确
func CheckSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
	if provider == SmsHttp {
		if s.ApiUrl == "" {
			return errors.New("HTTP短信接口必须提供API URL")
		}
		if strings.Index(s.ApiUrl, "{key}") == -1 {
			return errors.New("API URL缺少\"{key}\"字段")
		}
		if strings.Index(s.ApiUrl, "{secret}") == -1 {
			return errors.New("API URL缺少\"{secret}\"字段")
		}
		if strings.Index(s.ApiUrl, "{phone}") == -1 {
			return errors.New("API URL缺少\"{phone}\"字段")
		}
		if strings.Index(s.ApiUrl, "{msg}") == -1 {
			return errors.New("API URL缺少\"{msg}\"字段")
		}
		if s.SuccessChar == "" {
			return errors.New("请指定发送成功包含的字符")
		}
	}
	return nil
}

// 通过HTTP-API发送短信,successChar为发送成功包含的字符,enc：编码
func sendPhoneMsgByHttpApi(apiUrl, key, secret, phone, msg,
	enc, successChar string) error {
	//如果指定了编码，则先编码内容
	if enc != "" {
		dst, err := EncodingTransform([]byte(msg), enc)
		if err != nil {
			return err
		}
		msg = string(dst)
	}
	strUrl := compile(apiUrl, map[string]interface{}{
		"key":    key,
		"secret": secret,
		"phone":  phone,
		"msg":    url.QueryEscape(msg),
	})
	rsp, err := http.Get(strUrl)
	if err == nil {
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			err = errors.New("error : " + strconv.Itoa(rsp.StatusCode))
		}
		var data []byte
		data, err = ioutil.ReadAll(rsp.Body)
		if err == nil {
			result := string(data)
			if strings.Index(result, successChar) == -1 {
				err = errors.New("send fail : " + result)
			}
		}
	}
	return err
}

//编码
func EncodingTransform(src []byte, enc string) ([]byte, error) {
	var ec encoding.Encoding
	switch enc {
	default:
		return src, nil
	case "GBK":
		ec = simplifiedchinese.GBK
	case "GB2312":
		ec = simplifiedchinese.HZGB2312
	case "BIG5":
		ec = traditionalchinese.Big5
	}
	dst := make([]byte, len(src)*2)
	n, _, err := ec.NewEncoder().Transform(dst, src, true)
	return dst[:n], err
}
