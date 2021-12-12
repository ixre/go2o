package api

import (
	"go2o/core/infrastructure/domain"
	"testing"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : register_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 10:30
 * description :
 * history :
 */
const phone = "13163000002"
const token = "sTRz2UhO4h"

func TestRegisterApi_SendRegisterCode(t *testing.T) {
	mp := map[string]string{}
	testApi(t, "register.get_token", mp, true)
}

func TestRegisterApi_SendRegisterCode2(t *testing.T) {
	mp := map[string]string{
		"phone": phone,
		"token": token,
	}
	testApi(t, "register.send_code", mp, true)
}

func TestRegisterApi_Register(t *testing.T) {
	mp := map[string]string{
		"phone":       phone,
		"token":       token,
		"pwd":         domain.Md5("123456"),
		"reg_from":    "app",
		"invite_code": "",
		"check_code":  "5993",
	}
	testApi(t, "register.submit", mp, true)

}
