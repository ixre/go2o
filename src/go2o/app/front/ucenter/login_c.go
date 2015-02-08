/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof/app"
	"go2o/core/service/goclient"
	"net/http"
	"time"
)

type loginC struct {
	app.Context
}

//登陆
func (this *loginC) Login(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/ucenter/login.html", nil)
}
func (this *loginC) Login_post(w http.ResponseWriter, r *http.Request) {
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
func (this *loginC) Partner_connect(w http.ResponseWriter, r *http.Request) {
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
func (this *loginC) Partner_disconnect(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ms_token")
	if err == nil {
		cookie.Expires = time.Now().Add(time.Hour * -48)
		http.SetCookie(w, cookie)
	}
	w.Write([]byte("{state:1}"))
}
