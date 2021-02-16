package api

import (
	"errors"
	"github.com/ixre/gof/crypto"
	api "github.com/ixre/gof/jwt-api"
	http2 "github.com/ixre/gof/util/http"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	tc *api.Client
)

var (
	RInternalError = &api.Response{
		Code:    api.RCInternalError,
		Message: "内部服务器出错",
	}
	RAccessDenied = &api.Response{
		Code:    api.RCAccessDenied,
		Message: "没有权限访问该接口",
	}
	RIncorrectApiParams = &api.Response{
		Code:    api.RCNotAuthorized,
		Message: "缺少接口参数，请联系技术人员解决",
	}
	RUndefinedApi = &api.Response{
		Code:    api.RCUndefinedApi,
		Message: "调用的API名称不正确",
	}
)

func init() {
	server := "http://localhost:1428/a/v2"
	md5Secret := string(crypto.Md5([]byte("123456")))
	tc = api.NewClient(server, "go2o", md5Secret)
	tc.UseToken(func(key, secret string) string {
		r, err1 := http.Get(server + "/access_token?key=" + key + "&secret=" + secret)
		if err1 != nil {
			println("---获取accessToken失败", err1.Error())
			return ""
		}
		bytes, _ := ioutil.ReadAll(r.Body)
		return string(bytes)
	}, 30000)
	tc.HandleError(func(code int, message string) error {
		switch code {
		case api.RCAccessDenied:
			message = RAccessDenied.Message
		case api.RCNotAuthorized:
			message = RIncorrectApiParams.Message
		case api.RCUndefinedApi:
			message = RUndefinedApi.Message
		}
		return errors.New(message)
	})
}

// 测试提交
func testPost(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	params["version"] = "1.0.0"
	rsp, err := tc.Post(apiName, params)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}

// 测试提交
func testPostForm(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	params["version"] = "1.0.0"
	rsp, err := tc.Post(apiName, params)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}

// 测试提交
func testGET(t *testing.T, apiName string, params map[string]string) ([]byte, error) {
	params["version"] = "1.0.0"
	query := http2.ParseUrlValues(params).Encode()
	rsp, err := tc.Get(apiName+"?"+query, nil)
	t.Log("[ Response]:", string(rsp))
	if err != nil {
		t.Error(err)
		//t.FailNow()
	}
	return rsp, err
}
