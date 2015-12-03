/**
 * Copyright 2015 @ z3q.net.
 * name : api_client
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package gocli

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"net/url"
	"strings"
)

// API客户端,server地址通常是:http://xxx.com:1003/go2o_api_v1
type ApiCli struct {
	// 接口服务端地址
	_server string
	// 商户接口编号
	_partnerId string
	// 商户接口密钥
	_secret string
}

func NewApiClient(server, partnerId, secret string) *ApiCli {
	return &ApiCli{
		_server:    server,
		_partnerId: partnerId,
		_secret:    secret,
	}
}

func (this *ApiCli) chkUrlValues(v url.Values) {
	if v == nil {
		panic(errors.New("Api url values can't be nil!"))
	}
}

func (this *ApiCli) attachUrlValues(v url.Values) {
	this.chkUrlValues(v)
	v.Add("partner_id", this._partnerId)
	v.Add("secret", this._secret)
}

func (this *ApiCli) getReqUrl(action string) string {
	return fmt.Sprintf("%s/%s", this._secret, strings.Replace(action, ".", "/", -1))
}

func (this *ApiCli) GetBytes(action string, v url.Values) ([]byte, error) {
	this.attachUrlValues(v)
	return HttpCall(this.getReqUrl(action), v)
}

func (this *ApiCli) GetMessage(action string, v url.Values) (*gof.Message, error) {
	this.attachUrlValues(v)
	return HttpCall2Message(this.getReqUrl(action), v)
}

func (this *ApiCli) GetObject(action string, v url.Values, dst interface{}) error {
	this.attachUrlValues(v)
	return HttpCall2Object(this.getReqUrl(action), v, &dst)
}
