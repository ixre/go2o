package uams

import (
	"errors"
	"github.com/jsix/gof/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	RPermissionDenied = &api.Response{
		ErrCode: api.RPermissionDenied.ErrCode,
		ErrMsg:  "没有权限访问该接口",
	}
	RMissingApiParams = &api.Response{
		ErrCode: api.RMissingApiParams.ErrCode,
		ErrMsg:  "缺少接口参数，请联系技术人员解决",
	}
	RErrApiName = &api.Response{
		ErrCode: api.RErrUndefinedApi.ErrCode,
		ErrMsg:  "调用的API名称不正确",
	}
	RNoSuchApp = &api.Response{
		ErrCode: 10096,
		ErrMsg:  "no such app",
	}
)

var (
	API_SERVER    = "http://localhost:1419/uams_api_v1"
	API_USER      = "< replace your api user >"
	API_TOKEN     = "< replace your api token >"
	API_APP       = "< replace your serve code >"
	API_SIGN_TYPE = "sha1" // [sha1|md5]
)

// 请求接口
func Post(apiName string, data map[string]string) ([]byte, error) {
	cli := &http.Client{}
	form := url.Values{}
	if data != nil {
		for k, v := range data {
			form[k] = []string{v}
		}
	}
	form["api"] = []string{apiName}
	form["key"] = []string{API_USER}
	form["sign_type"] = []string{API_SIGN_TYPE}
	form["app"] = []string{API_APP}
	form["version"] = []string{"1.2.0.100"}
	sign := api.Sign(API_SIGN_TYPE, form, API_TOKEN)
	form["sign"] = []string{sign}
	rsp, err := cli.PostForm(API_SERVER, form)
	if err == nil {
		data, err := ioutil.ReadAll(rsp.Body)
		if err == nil {
			content := string(data)
			if strings.HasPrefix(content, "!") {
				arr := strings.Split(content[1:], ":")
				code, _ := strconv.Atoi(arr[0])
				return data, checkApiRespErr(code, arr[1])
			}
			return data, nil
		}
	}
	return []byte{}, err
}

// 如果返回接口请求错误, 响应状态码以-10开头
func checkApiRespErr(code int, text string) error {
	switch int64(code) {
	case api.RPermissionDenied.ErrCode:
		text = RPermissionDenied.ErrMsg
	case api.RMissingApiParams.ErrCode:
		text = RMissingApiParams.ErrMsg
	case api.RErrUndefinedApi.ErrCode:
		text = RErrApiName.ErrMsg
	case RNoSuchApp.ErrCode:
		text = RNoSuchApp.ErrMsg
	}
	return errors.New(text)
}
