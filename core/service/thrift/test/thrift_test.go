/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-13 13:11
 * description :
 * history :
 */
package test

import (
	"errors"
	"github.com/ixre/gof/log"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/auto_gen/rpc/foundation_service"
	"go2o/core/service/thrift"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	trans, cli, err := thrift.MemberServeClient()
	if err != nil {
		t.Error(err)
		return
	}

	defer trans.Close()
	pwd := domain.MemberSha1Pwd("123456")
	r, _ := cli.CheckLogin(thrift.Context, "jarry6", pwd, false)
	t.Logf("登录(1)结果：\n MemberId:%d\n Result:%v", r.ID, r.Result_)

	pwd = domain.MemberSha1Pwd("123000")
	r, _ = cli.CheckLogin(thrift.Context, "jarry6", pwd, false)
	t.Logf("登录(2)结果：\n MemberId:%d\n Result:%v", r.ID, r.Result_)

	arr, _ := cli.InviterArray(thrift.Context, 16893, 5)
	t.Log("邀请人：", arr)
}

func TestSSORegister(t *testing.T) {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		sa := &foundation_service.SSsoApp{
			ID:     1,
			Name:   "gp",
			ApiUrl: "http://localhost:14281/member/sync_m.p",
		}
		s, _ := cli.RegisterApp(thrift.Context, sa)
		arr := strings.Split(s, ":")
		if arr[0] != "1" {
			t.Error(errors.New("注册SSO-APP出错：" +
				s + "; api-url:" + sa.ApiUrl))
		} else {
			log.Println("得到的token为：", arr[1])
		}
	} else {
		t.Log("连接失败：", err.Error())
	}
}
