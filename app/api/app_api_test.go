package api

import (
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : app_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 10:30
 * description :
 * history :
 */

func TestAppApi_Process(t *testing.T) {
	mp := map[string]string{}
	mp["prod_type"] = "android"
	mp["prod_version"] = "1.0.0"
	testApi(t, "app.check", mp, true)
	mp["prod_type"] = "ios"
	mp["prod_version"] = "1.0.10"
	testApi(t, "app.check", mp, true)
}
