package api

import "testing"

/**
 * Copyright 2009-2019 @ to2.net
 * name : res_api_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-26 17:22
 * description :
 * history :
 */

func TestAdApi(t *testing.T) {
	mp := map[string]string{
		"pos_keys": "mobi-index-scroller",
		"user_id":  "0",
	}
	testApi(t, "res.ad_api", mp, true)
}
