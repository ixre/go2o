package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ixre/gof/api"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 10:47
 * description :
 * history :
 */

var serverUrl = "http://localhost:1428/api"

func testApi(t *testing.T, apiName string, paramsMap map[string]string, abortOnFail bool) {
	key := "go2o"
	secret := "123456"
	signType := "sha1"
	params := url.Values{}
	params["key"] = []string{key}
	params["api"] = []string{apiName}
	params["key"] = []string{key}
	params["sign_type"] = []string{signType}
	params["version"] = []string{"1.0.15"}
	for k, v := range paramsMap {
		params[k] = []string{v}
	}
	sign := api.Sign(signType, params, secret)
	//t.Log("-- Sign:", sign)
	params["sign"] = []string{sign}
	cli := http.Client{}
	rsp, err := cli.PostForm(serverUrl, params)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := io.ReadAll(rsp.Body)
	rsp1 := api.Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != api.RSuccessCode {
		println("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		println("接口响应：", string(data))
		if abortOnFail {
			t.FailNow()
		}
	}
	println("接口响应：", string(data))
}

// 测试请求限制
func TestRequestLimit(t *testing.T) {
	mp := map[string]string{}
	mp["prod_type"] = "android"
	mp["prod_version"] = "1.0.0"
	for {
		for i := 0; i < 100; i++ {
			testApi(t, "app.check", mp, false)
		}
		time.Sleep(time.Second)
	}
}

func TestSign(t *testing.T) {
	params := "api=member.login&key=go2o&product=app&pwd=c4ca4238a0b923820dcc509a6f75849b&user=18666398028&version=1.0.0&sign_type=sha1&sign=2933eaffccf9fe49a0ad9a97fe311a41afb6e3b2"
	values, _ := url.ParseQuery(params)
	sign := api.Sign("sha1", values, "131409")
	if sign2 := values.Get("sign"); sign2 != sign {
		println(sign, "/", sign2)
		t.Failed()
	}
	cli := http.Client{}
	rsp, err := cli.PostForm("http://localhost:1428/api", values)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, _ := io.ReadAll(rsp.Body)
	println(string(data))
}
