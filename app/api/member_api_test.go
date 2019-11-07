package api

import (
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : article_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-09-04 18:13
 * description :
 * history :
 */

func TestMemberApi_Process(t *testing.T) {
	mp := map[string]string{
		"code": "eNe6FR",
	}
	testApi(t, "member.invites", mp, true)
}

func TestMemberCheckToken_Process(t *testing.T) {
	mp := map[string]string{
		"code":  "m00U41",
		"token": "4e3fa6045473d5e44017558150",
	}
	testApi(t, "member.checkToken", mp, true)
}
