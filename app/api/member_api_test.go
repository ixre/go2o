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
