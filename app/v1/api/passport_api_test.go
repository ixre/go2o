package api

import (
	"go2o/core/infrastructure/domain"
	"go2o/core/service/proto"
	"strconv"
	"testing"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : passport_api_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-30 17:38
 * description :
 * history :
 */

const phone2 = "13162222872"
const token2 = "i3C3WvpTmC"

func TestRegisterApi_GetToken(t *testing.T) {
	mp := map[string]string{}
	testApi(t, "passport.get_token", mp, true)
}

func TestPassportApi_SendCode(t *testing.T) {
	mp := map[string]string{
		"account":   phone2,
		"cred_type": strconv.Itoa(int(proto.ECredentials_PHONE)),
		"token":     token2,
		"op":        "1",
	}
	testApi(t, "passport.send_code", mp, true)
}

func TestPassportApi_CompareCode(t *testing.T) {
	mp := map[string]string{
		"account":    phone2,
		"cred_type":  strconv.Itoa(int(proto.ECredentials_PHONE)),
		"token":      token2,
		"op":         "1",
		"check_code": "8799",
	}
	testApi(t, "passport.compare_code", mp, true)
}

func TestPassportApi_ResetPwd(t *testing.T) {
	mp := map[string]string{
		"account":   phone2,
		"cred_type": strconv.Itoa(int(proto.ECredentials_PHONE)),
		"token":     token2,
		"pwd":       domain.Md5("123456"),
	}
	testApi(t, "passport.reset_pwd", mp, true)
}

func TestPassportApi_ModifyPassword(t *testing.T) {
	mp := map[string]string{
		"account":   phone2,
		"cred_type": strconv.Itoa(int(proto.ECredentials_PHONE)),
		"token":     token2,
		"pwd":       domain.Md5("123000"),
		"old_pwd":   domain.Md5("123456"),
	}
	testApi(t, "passport.modify_pwd", mp, true)
}

func TestPassportApi_TradePwd(t *testing.T) {
	mp := map[string]string{
		"account":   phone2,
		"cred_type": strconv.Itoa(int(proto.ECredentials_PHONE)),
		"token":     token2,
		"pwd":       domain.Md5("123000"),
		"old_pwd":   domain.Md5("237561"),
	}
	testApi(t, "passport.trade_pwd", mp, true)
}
