/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-13 13:11
 * description :
 * history :
 */
package thrift

import (
	"errors"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/thrift/idl/gen-go/define"
	"gp/src/controller"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	cli, err := MemberServeClient()
	if err != nil {
		t.Error(err)
		return
	}

	defer cli.Transport.Close()
	t.Logf("连接开启：%v", cli.Transport.IsOpen())
	pwd := domain.MemberSha1Pwd("123456")
	mp, err := cli.Login("jarry6", pwd, false)
	t.Logf("登陆(1)结果：\n MemberId:%d\n Result:%d", mp["Id"], mp["Result"])

	t.Logf("%#v", mp)
	pwd = domain.MemberSha1Pwd("123000")
	mp, _ = cli.Login("jarry6", pwd, false)
	t.Logf("登陆(2)结果：\n MemberId:%d\n Result:%d", mp["Id"], mp["Result"])

}

func TestSSORegister(t *testing.T) {
	cli, err := FoundationServeClient()
	if err == nil {
		defer cli.Transport.Close()
		sa := &define.SsoApp{
			ID:     1,
			Name:   "gp",
			ApiUrl: "http://localhost:14281/member/sync_m.p",
		}
		s, _ := cli.RegisterSsoApp(sa)
		arr := strings.Split(s, ":")
		if arr[0] != "1" {
			t.Error(errors.New("注册SSO-APP出错：" +
				s + "; api-url:" + sa.ApiUrl))
		} else {
			controller.SetSsoToken(arr[1])
		}
	} else {
		t.Log("连接失败：", err.Error())
	}
}
