package service

import (
	"context"
	"errors"
	"testing"
	"time"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/typeconv"
)

// 测试发送验证码
func Test_checkServiceImpl_SendCode(t *testing.T) {
	svc := inject.GetCheckService()
	r := &proto.SendCheckCodeRequest{
		Token:         "testtoken",
		ReceptAccount: "13162222872",
		UserId:        1,
		Operation:     "注册会员",
		TemplateCode:  mss.SMS_CHECK_CODE,
	}
	ret, _ := svc.SendCode(context.TODO(), r)
	if len(ret.ErrMsg) != 0 {
		t.Error(errors.New(ret.ErrMsg))
		t.FailNow()
	}
	t.Logf("验证码为:%s", ret.CheckCode)

}

// 测试比较验证码
func Test_checkServiceImpl_CompareCode(t *testing.T) {
	s := inject.GetCheckService()
	req := &proto.CompareCheckCodeRequest{
		ReceptAccount: "13162222872",
		CheckCode:     "564992",
		Token:         "testtoken",
	}
	ret, _ := s.CompareCode(context.TODO(), req)
	if len(ret.ErrMsg) != 0 {
		t.Error(errors.New(ret.ErrMsg))
		t.FailNow()
	}
}

// TestGrantMemberAccessToken 测试GrantMemberAccessToken函数
func TestGrantMemberAccessToken(t *testing.T) {
	var memberId int64 = 1
	s := inject.GetCheckService()
	token, _ := s.GrantAccessToken(context.TODO(), &proto.GrantAccessTokenRequest{
		UserId:      memberId,
		UserType:    1,
		ExpiresTime: time.Now().Unix() + 720,
	})
	if len(token.ErrMsg) > 0 {
		t.Error(token.ErrMsg)
		t.FailNow()
	}
	t.Log("token is:", token.AccessToken)
	now := time.Now().Unix()
	token.AccessToken = "Bearer " + token.AccessToken
	accessToken, _ := s.CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken:      token.AccessToken,
		CheckExpireTime:  now + 800,
		RenewExpiresTime: now + 900,
	})
	if accessToken.UserId != memberId {
		t.Error(accessToken.ErrMsg)
		t.Failed()
	}
	t.Logf("user token:%s", token.AccessToken)
}

// TestCheckMemberAccessToken 测试 CheckMemberAccessToken 函数
func TestCheckMemberAccessToken(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI3MTAiLCJleHAiOjE2Nzc3NTM5MjAsImlzcyI6ImdvMm8iLCJzdWIiOiJnbzJvLWFwaS1qd3QifQ.Vrd0NuAT5AfKM-C5NespEyhAiyMVDugoKDDgwv5hr_g"
	ret, _ := inject.GetCheckService().CheckAccessToken(context.TODO(), &proto.CheckAccessTokenRequest{
		AccessToken: accessToken,
	})
	if len(ret.ErrMsg) > 0 {
		t.Log(ret.ErrMsg)
		t.FailNow()
	}
	t.Log(typeconv.MustJson(ret))
	t.Log("会员Id", ret.UserId)
}
