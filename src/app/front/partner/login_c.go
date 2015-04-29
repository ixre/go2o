/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/app/session"
)

type loginC struct {
	gof.App
}

//登陆
func (this *loginC) Login(ctx *web.Context) {
	this.App.Template().Render(ctx.ResponseWriter, "views/partner/login.html", nil)
}
func (this *loginC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	session.GetLSession().WebValidLogin(w, r.Form.Get("uid"), r.Form.Get("pwd"))
}
