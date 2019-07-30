package api

import (
	"go2o/core/service/auto_gen/rpc/member_service"
	"strconv"
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : passport_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 17:38
 * description :
 * history :
 */

const phone2 = "13162222872"
const token2 = "XyjmMoDzEm"

func TestRegisterApi_GetToken(t *testing.T) {
	mp := map[string]string{}
	testApi(t, "passport.get_token", mp)
}

func TestPassportApi_SendCode(t *testing.T) {
	mp := map[string]string{
		"account":   phone2,
		"cred_type": strconv.Itoa(int(member_service.ECredentials_Phone)),
		"token":     token2,
		"op":        "1",
	}
	testApi(t, "passport.send_code", mp)
}

func TestPassportApi_CompareCode(t *testing.T) {
	mp := map[string]string{
		"account":    phone2,
		"cred_type":  strconv.Itoa(int(member_service.ECredentials_Phone)),
		"token":      token2,
		"op":         "1",
		"check_code": "7457",
	}
	testApi(t, "passport.compare_code", mp)
}
