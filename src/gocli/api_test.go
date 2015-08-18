/**
 * Copyright 2015 @ z3q.net.
 * name : api_test.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package gocli

import (
	"net/url"
	"testing"
)

func ApiCall_test(t *testing.T) {
	cli := &NewApiClient("http://localhost:1003/go2o_api_v1", "partner_id", "partner_secret")
	v := url.Values{
		"usr": {"user"},
		"pwd": {"pwd"},
	}
	if msg, err := cli.GetMessage("mm_login", v); err != nil {
		t.Error(err)
	} else if !msg.Result {
		t.Fail()
		t.Error(msg.Message)
	}

}
