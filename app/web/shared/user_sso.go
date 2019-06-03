/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package shared

import (
	"errors"
	"github.com/ixre/goex/echox"
	gu "github.com/ixre/gof/util"
	"go2o/app/util"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"net/http"
	"strconv"
	"time"
)

var (
	errCredential = errors.New("error credentital")
)

// 获取会员编号
func GetMemberId(c *echox.Context) int {
	v := c.Session.Get("member_id")
	switch v.(type) {
	case int64:
		return int(v.(int64))
	case int:
		return v.(int)
	}
	return 0
}

// 获取会员
func GetMember(c *echox.Context) *member_service.SMember {
	memberId := GetMemberId(c)
	if memberId > 0 {
		m, _ := rsi.MemberService.GetMember(thrift.Context, int64(memberId))
		return m
	}
	return nil
}

type UserSync struct {
}

// 同步登录/登出
func (u *UserSync) Sync(c *echox.Context) (err error) {
	// 获取登录、登出
	out := c.QueryParam("out") == "true"
	// 获取回调函数方法
	callback := c.QueryParam("callback")
	if callback == "" {
		callback = "sso_callback"
	}
	// 设置访问设备
	if device := c.QueryParam("device"); len(device) > 0 {
		util.SetBrownerDevice(c.Response(), c.Request(), device)
	}
	// 处理登录或登出
	if out {
		err = u.ssoDisconnect(c, callback)
	} else {
		err = u.ssoConnect(c, callback)
	}
	// 返回结果
	if err == nil {
		return c.JSONP(http.StatusOK, callback, "success")
	}
	return c.JSONP(http.StatusOK, callback, err.Error())
}

//通过URL参数登录
//@member_id : 会员编号
//@token  :  密钥/令牌
//@device : 设备类型
func (u *UserSync) ssoConnect(c *echox.Context, callback string) error {
	// 第三方连接，传入memberId 和 token
	mStr := c.QueryParam("member_id")
	mId, err := gu.I64Err(strconv.Atoi(mStr))
	// 鉴权，如成功，则存储会话
	token := c.QueryParam("token")
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		b, _ := cli.CheckToken(thrift.Context, mId, token)
		if b {
			c.Session.Set("member_id", mId)
			c.Session.Save()
			return nil
		}
	}
	// 鉴权失败
	return errCredential
}

//同步退出
func (u *UserSync) ssoDisconnect(c *echox.Context, callback string) error {
	// 消除会话
	c.Session.Destroy()
	rsp := c.Response()
	// 清理以"_"开头的cookie,假定以"_"开头均为与业务相关重要的cookie信息
	expires := time.Now().Local().Add(time.Hour * -1e5)
	list := c.Request().Cookies()
	for _, ck := range list {
		if ck.Name[0] == '_' {
			ck.Path = "/"
			ck.Expires = expires
			http.SetCookie(rsp, ck)
		}
	}
	return nil
}
