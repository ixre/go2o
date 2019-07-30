package api

import (
	"encoding/json"
	"github.com/ixre/gof/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 10:47
 * description :
 * history :
 */

func testApi(t *testing.T, apiName string, paramsMap map[string]string) {
	key := "go2o"
	secret := "131409"
	signType := "sha1"
	serverUrl := "http://localhost:1428/api"
	params := url.Values{}
	params["key"] = []string{key}
	params["api"] = []string{apiName}
	params["key"] = []string{key}
	params["sign_type"] = []string{signType}
	params["version"] = []string{"1.0.1"}
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
	data, _ := ioutil.ReadAll(rsp.Body)
	rsp1 := api.Response{}
	json.Unmarshal(data, &rsp1)
	if rsp1.Code != api.RSuccessCode {
		t.Log("请求失败：code:", rsp1.Code, "; message:", rsp1.Message)
		t.Log("接口响应：", string(data))
		t.FailNow()
	}
	t.Log("接口响应：", string(data))
}
