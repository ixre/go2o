package api

import (
	"go2o/core/infrastructure/domain"
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : register_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 10:30
 * description :
 * history :
 */

func TestRegisterApi_SendRegisterCode(t *testing.T) {
	mp := map[string]string{}
	testApi(t, "register.get_token", mp)
}

func TestRegisterApi_SendRegisterCode2(t *testing.T) {
	mp := map[string]string{
		"phone": "13162221120",
		"token": "0L0XIvcUyq",
	}
	testApi(t, "register.send_code", mp)
}

func TestRegisterApi_Register(t *testing.T) {
	mp := map[string]string{
		"phone":       "13162221121",
		"token":       "0L0XIvcUyq",
		"pwd":         domain.Md5("123456"),
		"reg_from":    "app",
		"invite_code": "",
	}
	testApi(t, "register.submit", mp)

}
