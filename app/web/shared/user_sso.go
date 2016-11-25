/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package shared

import (
    "github.com/jsix/goex/echox"
    "go2o/app/util"
    gu "github.com/jsix/gof/util"
    "net/http"
    "strconv"
    "time"
    "errors"
    "go2o/core/service/thrift"
)

var (
    errCredential = errors.New("error credentital")
)

type UserSync struct {
}

// 同步登陆/登出
func (u *UserSync) Sync(c *echox.Context) (err error) {
    // 获取登陆、登出
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
    // 处理登陆或登出
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
    mKey := c.QueryParam("member_id")
    mId, err := gu.I32Err(strconv.Atoi(mKey))
    // 鉴权，如成功，则存储会话
    token := c.QueryParam("token")
    cli, err := thrift.MemberClient()
    if err == nil {
        defer cli.Transport.Close()
        b, _ := cli.CheckToken(mId, token)
        if b {
            c.Session.Set("member", mId)
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
    d := time.Duration(-1e9)
    expires := time.Now().Local().Add(d)
    list := c.Request().Cookies()
    for _, ck := range list {
        if ck.Name[0] == '_' {
            ck.Expires = expires
            http.SetCookie(rsp, ck)
        }
    }
    return nil
}
