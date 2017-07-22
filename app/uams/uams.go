package uams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	RPermissionDenied = &Response{
		Result:  -100,
		Message: "没有权限访问该接口",
	}
	RMissingApiParams = &Response{
		Result:  -101,
		Message: "缺少接口参数，请联系技术人员解决",
	}
	RErrApiName = &Response{
		Result:  -102,
		Message: "调用的API名称不正确",
	}
)

var (
	API_SERVER    = "http://localhost:1419/uams_api_v1"
	API_USER      = "< replace your api user >"
	API_TOKEN     = "< replace your api token >"
	API_SIGN_TYPE = "sha1" // [sha1|md5]
)

// 请求接口
func Post(api string, data map[string]string) (string, error) {
	cli := &http.Client{}
	form := url.Values{}
	if data != nil {
		for k, v := range data {
			form[k] = []string{v}
		}
	}
	form["api"] = []string{api}
	form["api_user"] = []string{API_USER}
	form["sign_type"] = []string{API_SIGN_TYPE}
	sign := Sign(API_SIGN_TYPE, form, API_TOKEN)
	form["sign"] = []string{sign}
	rsp, err := cli.PostForm(API_SERVER, form)
	if err == nil {
		result, err := ioutil.ReadAll(rsp.Body)
		if err == nil {
			err = checkApiRespErr(result)
		}
		return string(result), err
	}
	return "", err
}

// 如果返回接口请求错误, 响应状态码以-10开头
func checkApiRespErr(result []byte) error {
	if bytes.Index(result, []byte{'-', '1', '0'}) != -1 {
		msg := Response{}
		json.Unmarshal(result, &msg)
		switch msg.Result {
		case -100:
			msg = *RPermissionDenied
		case -101:
			msg = *RMissingApiParams
		case -102:
			msg = *RErrApiName
		}
		return errors.New(fmt.Sprintf(
			"Error code %d : %s", msg.Result, msg.Message))
	}
	return nil
}
