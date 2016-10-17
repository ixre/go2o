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
	"github.com/jsix/gof"
	"go2o/app/util"
	"go2o/core/service/dps"
	"go2o/x/echox"
	"net/http"
	"strconv"
	"time"
)

type UserC struct {
	gof.App
}

//通过URL参数登录
//@member_id : 会员编号
//@token  :  密钥/令牌
//@device : 设备类型
func (u *UserC) Connect(ctx *echox.Context) error {
	// 获取回调函数方法
	callback := ctx.Query("callback")
	if callback == "" {
		callback = "sso_callback"
	}
	//设置访问设备
	if device := ctx.Query("device"); len(device) > 0 {
		util.SetBrownerDevice(ctx.Response(), ctx.Request(), device)
	}
	// 第三方连接，传入memberId 和 token
	memberId, err := strconv.Atoi(ctx.Query("member_id"))
	if err != nil {
		memberId = 0
	}
	// 鉴权，如成功，则存储会话
	token := ctx.Query("token")
	sto := ctx.App.Storage()
	if util.CompareMemberApiToken(sto, memberId, token) {
		m := dps.MemberService.GetMember(memberId)
		ctx.Session.Set("member", m)
		ctx.Session.Save()
		return ctx.JSONP(http.StatusOK, callback, "success")
	}
	// 鉴权失败
	return ctx.JSONP(http.StatusOK, callback, "error credentital")
}

//同步退出
func (u *UserC) Disconnect(ctx *echox.Context) error {
	// 获取回调函数方法
	callback := ctx.Query("callback")
	if callback == "" {
		callback = "sso_callback"
	}
	// 消除会话
	ctx.Session.Destroy()
	rsp := ctx.Response()
	// 清理以"_"开头的cookie
	d := time.Duration(-1e9)
	expires := time.Now().Local().Add(d)
	list := ctx.Request().Cookies()
	for _, ck := range list {
		if ck.Name[0] == '_' {
			ck.Expires = expires
			http.SetCookie(rsp, ck)
		}
	}
	// 清理购物车
	//ck := &http.Cookie{
	//    Name:     "_cart",
	//    Value:    "",
	//    Path:     "/",
	//    HttpOnly: true,
	//    Expires:  expires,
	//}
	//http.SetCookie(ctx.Response(), ck)
	return ctx.JSONP(http.StatusOK, callback, "success")
}
