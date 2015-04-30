/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

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
	this.App.Template().Render(ctx.ResponseWriter, "views/ucenter/login.html", nil)
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
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("{result:true}"))
	}
}

//从partner登录过来的信息
func (this *loginC) Partner_connect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	token := r.URL.Query().Get("token")

	if token == "" {
		w.Write([]byte("<script>location.replace('/login')</script>"))
	} else {
		cookie := &http.Cookie{
			Name:    "ms_token",
			Expires: time.Now().Add(time.Hour * 48),
			Value:   token,
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("<script>location.replace('/')</script>"))
	}
}

//从partner端退出
func (this *loginC) Partner_disconnect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	cookie, err := r.Cookie("ms_token")
	if err == nil {
		cookie.Expires = time.Now().Add(time.Hour * -48)
		http.SetCookie(w, cookie)
	}
	w.Write([]byte("{state:1}"))
}
