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

func TestSettingsApi_RegisterPerm(t *testing.T) {
	mp := map[string]string{}
	testApi(t, "settings.register_settings", mp, true)
}
