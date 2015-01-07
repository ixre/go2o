package partner

import (
	"com/ording/session"
	"net/http"
	"github.com/newmin/gof/app"
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
