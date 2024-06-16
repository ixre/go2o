package sms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/gof/typeconv"
	"github.com/ixre/gof/util"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// 通过HTTP-API发送短信, 短信模板参数在data里指定
func sendPhoneMsgByHttpApi(api *mss.SmsApiPerm, phone, content string, data []string, templateId string) error {
	if api.Extra == nil {
		api.Extra = &mss.SmsExtraSetting{}
	}
	//如果指定了编码，则先编码内容
	if api.Extra.Charset != "" {
		dst, err := EncodingTransform([]byte(content), api.Extra.Charset)
		if err != nil {
			return err
		}
		content = string(dst)
	}
	// 如果GET发送,需要编码
	if api.Extra.Method == "GET" {
		content = url.QueryEscape(content)
	}
	// 请求参数
	params := map[string]string{
		"key":          api.Key,
		"secret":       api.Secret,
		"phone":        phone,
		"content":      content,
		"templateId":   templateId,
		"templateData": strings.Join(data, ","),
		"stamp":        fmt.Sprintf("%s%d", util.RandString(3), time.Now().Unix()),
	}
	body := resolveApiRequestParams(api.Extra.Params, params)

	// 创建请求
	req, err := createHttpRequest(api, body)
	if err != nil {
		return err
	}
	cli := &http.Client{}
	// 忽略证书
	if req.TLS != nil || strings.HasPrefix(api.Extra.ApiUrl, "https://") {
		cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	// 读取响应
	rsp, err := cli.Do(req)
	if err == nil {
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("error : %d", rsp.StatusCode)
		}
		//log.Println("[ GO2O][ Sms]:", body)
		var data []byte
		data, err = io.ReadAll(rsp.Body)
		if err == nil {
			result := string(data)
			if !strings.Contains(result, api.Extra.SuccessChars) {
				return errors.New("send fail : " + result + " message body:" + content)
			}
		}
	}
	return err
}

// 解析HTTP短信中的请求参数
func resolveApiRequestParams(params string, data map[string]string) string {
	for k, v := range data {
		str, _ := typeconv.String(v)
		params = strings.Replace(params, "{"+k+"}",
			str, -1)
	}
	return params
}

// 创建HTTP短信发送请求
func createHttpRequest(api *mss.SmsApiPerm, body string) (*http.Request, error) {
	var req *http.Request
	var err error
	if api.Extra.Method == "POST" {
		req, err = http.NewRequest(api.Extra.Method, api.Extra.ApiUrl, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		url := api.Extra.ApiUrl
		if strings.Contains(api.Extra.ApiUrl, "?") {
			url += "&"
		} else {
			url += "?"
		}
		req, err = http.NewRequest(api.Extra.Method, url+body, nil)
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
