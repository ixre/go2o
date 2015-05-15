/**
 * Copyright 2014 @ S1N1 Team.
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
)

type loginC struct {
	gof.App
}

//登陆
func (this *loginC) Login(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	ctx.App.Template().Execute(w, "views/ucenter/login.html", nil)
}
func (this *loginC) Login_post(ctx *web.Context) {
	//todo:
	//	r, w := ctx.Request, ctx.ResponseWriter
	//	r.ParseForm()
	//	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	//	result,_ := goclient.Member.Login(usr, pwd)
	//	if !result.Result {
	//		w.Write([]byte("{result:false,message:'" + result.Message + "'}"))
	//	} else {
	//		cookie := &http.Cookie{
	//			Name:    "ms_token",
	//			Expires: time.Now().Add(time.Hour * 48),
	//			Value:   result.Token,
	//		}
	//		http.SetCookie(w, cookie)
	//		w.Write([]byte("{result:true}"))
	//	}
}
