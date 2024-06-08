package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
)

// 测试发送验证码
func Test_checkServiceImpl_SendCode(t *testing.T) {
	svc := inject.GetCheckService()
	r := &proto.SendCheckCodeRequest{
		Token:         "testtoken",
		ReceptAccount: "13162222872",
		UserId:        1,
		Operation:     "注册会员",
		MsgTemplateId: "验证手机",
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
