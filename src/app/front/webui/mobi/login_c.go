/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mobi

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/service/goclient"
	"net/http"
	"time"
)

type loginC struct {
	gof.App
}

//登陆
func (this *loginC) Login(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	this.App.Template().Render(w, "views/ucenter/login.html", nil)
}
func (this *loginC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	b, t, msg := goclient.Member.Login(usr, pwd)
	if !b {
		w.Write([]byte("{result:false,message:'" + msg + "'}"))
	} else {
		cookie := &http.Cookie{
			Name:    "ms_token",
			Expires: time.Now().Add(time.Hour * 48),
			Value:   t,
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("{result:true}"))
	}
}
