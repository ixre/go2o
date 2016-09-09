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
	//设置访问设备
	if device := ctx.Query("device"); len(device) > 0 {
		util.SetBrownerDevice(ctx.Response(), ctx.Request(), device)
	}
	// 第三方连接，传入memberId 和 token
	memberId, err := strconv.Atoi(ctx.Query("member_id"))
	token := ctx.Query("token")
	if err != nil || token == "" {
		return ctx.String(http.StatusOK, "会话不正确")
	}
	// 存储会话状态
	m := dps.MemberService.GetMember(memberId)
	ctx.Session.Set("member", m)
	ctx.Session.Save()
	return ctx.StringOK("ok")
}

//同步退出
func (u *UserC) Disconnect(ctx *echox.Context) error {
	ctx.Session.Destroy()

	d := time.Duration(-1e9)
	expires := time.Now().Local().Add(d)
	ck := &http.Cookie{
		Name:     "_cart",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	}
	http.SetCookie(ctx.Response(), ck)
	return ctx.StringOK("ok")
}
