package mobi

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
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("{result:true}"))
	}
}
