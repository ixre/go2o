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
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/mss/notify"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/aliyu"
	"github.com/ixre/go2o/core/infrastructure/tool/sms/cl253"
	"github.com/ixre/gof/util"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
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
	if api.Signature != "" && strings.Index(content, api.Signature) == -1 {
		content = api.Signature + content
	}
	c := compileArray(content, params)
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

// // 解析模板中的参数
// func compile(tpl string, param map[string]interface{}) string {
// 	var str string
// 	for k, v := range param {
// 		switch v.(type) {
// 		case string:
// 			str = v.(string)
// 		case int, int32, int64:
// 			str = strconv.Itoa(v.(int))
// 		case float32, float64:
// 			str = format.FormatFloat(v.(float32))
// 		case bool:
// 			str = strconv.FormatBool(v.(bool))
// 		default:
// 			str = "unknown"
// 		}
// 		tpl = strings.Replace(tpl, "{"+k+"}", str, -1)
// 	}
// 	return tpl
// }

// 解析模板中的参数
func compileArray(tpl string, param []string) string {
	for k, v := range param {
		tpl = strings.Replace(tpl, fmt.Sprintf("{%d}", k), v, -1)
	}
	return tpl
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
		if strings.Index(s.Params, "{key}") == -1 {
			return errors.New("API Params缺少\"{key}\"字段")
		}
		if strings.Index(s.Params, "{secret}") == -1 {
			return errors.New("API Params缺少\"{secret}\"字段")
		}
		if strings.Index(s.Params, "{phone}") == -1 {
			return errors.New("API Params缺少\"{phone}\"字段")
		}
		if strings.Index(s.Params, "{content}") == -1 {
			return errors.New("API Params缺少\"{content}\"字段")
		}
		if s.SuccessChar == "" {
			return errors.New("未指定发送成功包含的字符")
		}
	}
	return nil
}

// 通过HTTP-API发送短信, 短信模板参数在data里指定
func sendPhoneMsgByHttpApi(api *SmsApi, phone, content string, data []string, templateId string) error {
	//如果指定了编码，则先编码内容
	if api.Charset != "" {
		dst, err := EncodingTransform([]byte(content), api.Charset)
		if err != nil {
			return err
		}
		content = string(dst)
	}
	// 如果GET发送,需要编码
	if api.Method == "GET" {
		content = url.QueryEscape(content)
	}
	// 请求参数
	params := url.Values{
		"key":        {api.Key},
		"secret":     {api.Secret},
		"phone":      {phone},
		"content":    {content},
		"templateId": {templateId},
		"stamp":      {fmt.Sprintf("%s%d", util.RandString(3), time.Now().Unix())},
	}
	// 创建请求
	req, err := createHttpRequest(api, params.Encode())
	if err != nil {
		return err
	}
	cli := &http.Client{}
	// 忽略证书
	if req.TLS != nil || (len(api.ApiUrl) >= 8 && api.ApiUrl[:8] == "https://") {
		cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	// 读取响应
	rsp, err := cli.Do(req)
	if err == nil {
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			err = errors.New("error : " + strconv.Itoa(rsp.StatusCode))
		}
		//log.Println("[ Go2o][ Sms]:", body)
		var data []byte
		data, err = io.ReadAll(rsp.Body)
		if err == nil {
			result := string(data)
			if strings.Index(result, api.SuccessChar) == -1 {
				err = errors.New("send fail : " + result + " message body:" + content)
			}
		}
	}
	return err
}

// 创建HTTP短信发送请求
func createHttpRequest(api *SmsApi, body string) (*http.Request, error) {
	var req *http.Request
	var err error
	if api.Method == "POST" {
		req, err = http.NewRequest(api.Method, api.ApiUrl, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		url := api.ApiUrl
		if strings.Index(api.ApiUrl, "?") == -1 {
			url += "?"
		} else {
			url += "&"
		}
		req, err = http.NewRequest(api.Method, url+body, nil)
	}
	return req, err
}

// 编码
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
