/**
 * Copyright 2015 @ z3q.net.
 * name : api.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package gocli

import (
	"encoding/json"
	"github.com/jsix/gof"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func HttpCall(url string, v *url.Values) ([]byte, error) {
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}
	rsp, err := client.PostForm(url, *v)
	if err == nil {
		return ioutil.ReadAll(rsp.Body)
	}
	return nil, err
}

// Http调用并返回Hash
func HttpCall2Hash(url string, v *url.Values) (map[string]interface{}, error) {
	d, err := HttpCall(url, v)
	if err == nil {
		var m map[string]interface{} = make(map[string]interface{})
		err = json.Unmarshal(d, &m)
		return m, err
	}
	return nil, err
}

// Http调用并反序列化为对象
func HttpCall2Object(url string, v *url.Values, dst interface{}) error {
	d, err := HttpCall(url, v)
	if err == nil {
		err = json.Unmarshal(d, dst)
	}
	return err
}

// Http调用并反序列化为消息
func HttpCall2Message(url string, v *url.Values) (*gof.Message, error) {
	var message *gof.Message = new(gof.Message)
	return message, HttpCall2Object(url, v, message)
}
