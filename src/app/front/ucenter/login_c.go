/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/jsix/gof"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
	"strings"
)

type loginC struct {
}

//登陆
func (this *loginC) Index(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.index_post(ctx)
	}
	return ctx.RenderOK("login.html", ctx.NewData())
}
func (this *loginC) index_post(ctx *echox.Context) error {
	r := ctx.Request()
	r.ParseForm()
	var result gof.Message
	usr, pwd := r.FormValue("usr"), r.FormValue("pwd")

	pwd = strings.TrimSpace(pwd)

	m, err := dps.MemberService.TryLogin(-1, usr, pwd, true)
	if err == nil {
		ctx.Session.Set("member", m)
		ctx.Session.Save()
		result.Result = true
	} else {
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Message = "登陆失败"
		}
	}
	return ctx.JSON(http.StatusOK, result)

}

//从partner登录过来的信息
func (this *loginC) Partner_connect(ctx *echox.Context) error {
	r := ctx.Request()
	sessionId := r.URL.Query().Get("sessionId")
	var m *member.ValueMember
	var err error

	if sessionId == "" {
		// 第三方连接，传入memberId 和 token
		memberId, err := strconv.Atoi(r.URL.Query().Get("mid"))
		token := r.URL.Query().Get("token")
		if err == nil && token != "" {
			m := dps.MemberService.GetMember(memberId)
			ctx.Session.Set("member", m)
			ctx.Session.Save()
		}
	} else {
		// 从统一平台连接过来（标准版商户PC前端)
		ctx.Session.UseInstead(sessionId)
		m = ctx.Session.Get("member").(*member.ValueMember)
	}

	// 设置访问设备
	util.SetBrownerDevice(ctx.Response(), ctx.Request(), ctx.Query("device"))

	if err == nil || m != nil {
		rl := dps.MemberService.GetRelation(m.Id)
		if rl.RegisterPartnerId > 0 {
			ctx.Session.Set("member:rel_partner", rl.RegisterPartnerId)
			ctx.Session.Save()
			ctx.Redirect(302, "/")
			return nil
		}
	}
	ctx.Redirect(302, "/login")
	return nil
}

//从partner端退出
func (this *loginC) Partner_disconnect(ctx *echox.Context) error {
	ctx.Session.Destroy()
	ctx.Response().Write([]byte("{state:1}"))
	return nil
}
