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
	r := ctx.HttpRequest()
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
func (this *loginC) Partner_connect(c *echox.Context) error {
	r := c.Request()
	sessionId := r.URL.Query().Get("sessionId")
	url := r.URL.Query().Get("return_url")
	if len(url) == 0 {
		url = "/"
	}
	var m *member.ValueMember
	var err error

	if sessionId == "" {
		// 第三方连接，传入memberId 和 token
		memberId, err := strconv.Atoi(r.URL.Query().Get("mid"))
		token := r.URL.Query().Get("token")
		if err == nil && token != "" {
			m = dps.MemberService.GetMember(memberId)
			c.Session.Set("member", m)
			c.Session.Save()
		} else {
			return c.String(http.StatusOK, "会话不正确")
		}
	} else {
		// 从统一平台连接过来（标准版商户PC前端)
		c.Session.UseInstead(sessionId)
		m = c.Session.Get("member").(*member.ValueMember)
	}

	// 设置访问设备
	util.SetBrownerDevice(c.HttpResponse(), c.HttpRequest(), c.Query("device"))

	if err == nil && m != nil {
		rl := dps.MemberService.GetRelation(m.Id)
		if rl.RegisterPartnerId > 0 {
			c.Session.Set("member:rel_partner", rl.RegisterPartnerId)
			c.Session.Save()
			c.Redirect(302, url)
			return nil
		}
	}
	c.Redirect(302, "/login")
	return nil
}

//从partner端退出
func (this *loginC) Partner_disconnect(ctx *echox.Context) error {
	ctx.Session.Destroy()
	ctx.HttpResponse().Write([]byte("{state:1}"))
	return nil
}
