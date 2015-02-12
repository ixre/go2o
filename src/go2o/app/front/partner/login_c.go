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
	"github.com/atnet/gof/app"
	"go2o/app/session"
	"net/http"
)

type loginC struct {
	app.Context
}

//登陆
func (this *loginC) Login(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/partner/login.html", nil)
}
func (this *loginC) Login_post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session.LSession.WebValidLogin(w, r.Form.Get("uid"), r.Form.Get("pwd"))
}
