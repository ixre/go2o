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
	"go2o/core/infrastructure/domain"
	"testing"
)

func TestLogin(t *testing.T) {
	cli, err := MemberClient()
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
