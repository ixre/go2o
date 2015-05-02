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
	"strconv"
)

type loginC struct {
	gof.App
}

//登陆
func (this *loginC) Login(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter, "views/ucenter/login.html", nil)
}
func (this *loginC) Login_post(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		r.ParseForm()
		usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
		result,_ := goclient.Member.Login(usr, pwd)
		if !result.Result {
			w.Write([]byte("{result:false,message:'" + result.Message + "'}"))
		} else {
			cookie := &http.Cookie{
				Name:    "ms_token",
				Expires: time.Now().Add(time.Hour * 48),
				Value:   result.Token,
			}
			http.SetCookie(w, cookie)
			w.Write([]byte("{result:true}"))
		}
}

//从partner登录过来的信息
func (this *loginC) Partner_connect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	sessionId := r.URL.Query().Get("sessionId")


	if sessionId == "" {
		// 第三方连接，传入memberId 和 token
		memberId, err := strconv.Atoi(r.URL.Query().Get("mid"))
		token := r.URL.Query().Get("token")
		if err == nil && token != "" {
			m, err := goclient.Member.GetMember(memberId, token)
			if err == nil || m!= nil {
				ctx.Session().Set("member", m)
				ctx.Session().Save()
			}
			w.Write([]byte("<script>location.replace('/')</script>"))
		}
	}else{
		// 从统一平台连接过来（标准版商户PC前端)
		ctx.Session().UseInstead(sessionId)
		w.Write([]byte("<script>location.replace('/')</script>"))
	}

	w.Write([]byte("<script>location.replace('/login')</script>"))
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
