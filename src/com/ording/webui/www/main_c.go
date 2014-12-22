package www

import (
	"com/domain/interface/member"
	"com/ording/cache/apicache"
	"com/ording/entity"
	"com/service/goclient"
	"com/share/variable"
	"fmt"
	"html/template"
	"net/http"
	"ops/cf/app"
	"ops/cf/web"
	"strings"
	"time"
)

type mainC struct {
	app.Context
}

func (this *mainC) Index(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	if b, siteConf := GetSiteConf(w, p); b {
		shops := apicache.GetShops(this.Context, p.Id, p.Secret)
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["conf"] = siteConf
			(*m)["title"] = siteConf.IndexTitle
			(*m)["shops"] = template.HTML(shops)
		},
			"views/web/www/index.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) Login(w http.ResponseWriter, r *http.Request, p *entity.Partner, mm *member.ValueMember) {
	if b, siteConf := GetSiteConf(w, p); b {
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "会员登录－" + siteConf.SubTitle
			(*m)["member"] = mm
			(*m)["conf"] = siteConf
		},
			"views/web/www/login.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) Login_post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	b, t, msg := goclient.Member.Login(usr, pwd)
	if !b {
		w.Write([]byte("{result:false,message:'" + msg + "'}"))
	} else {
		cookie := &http.Cookie{
			Name:    "ms_token",
			Expires: time.Now().Add(time.Hour * 48),
			Path:    "/",
			Value:   t,
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("{result:true}"))
	}
}

func (this *mainC) Register(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	if b, siteConf := GetSiteConf(w, p); b {
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "会员注册－" + siteConf.SubTitle
			(*m)["conf"] = siteConf
		},
			"views/web/www/register.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) ValidUsr_post(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	r.ParseForm()
	usr := r.FormValue("usr")
	b := goclient.Partner.UserIsExist(p.Id, p.Secret, usr)
	if !b {
		w.Write([]byte(`{"result":true}`))
	} else {
		w.Write([]byte(`{"result":false}`))
	}
}

func (this *mainC) PostRegistInfo_post(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	r.ParseForm()
	var member member.ValueMember
	web.ParseFormToEntity(r.Form, &member)
	if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
		member.RegIp = r.RemoteAddr[:i]
	}
	b, err := goclient.Partner.RegistMember(&member, p.Id, 0, "")
	if b {
		w.Write([]byte(`{"result":true}`))
	} else {
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
		} else {
			w.Write([]byte(`{"result":false}`))
		}

	}
}

//跳转到会员中心
func (this *mainC) Member(w http.ResponseWriter, r *http.Request, p *entity.Partner, mm *member.ValueMember) {
	var location string
	if mm == nil {
		location = "/login?return_url=/member"
	} else {
		cookie, _ := r.Cookie("ms_token")
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?token=%s",
			variable.DOMAIN_MEMBER_PREFIX,
			this.Context.Config().Get(variable.ServerDomain),
			cookie.Value,
		)
	}
	w.Write([]byte("<script>window.parent.location.href='" + location + "'</script>"))
}

//退出
func (this *mainC) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ms_token")
	if err == nil {
		cookie.Expires = time.Now().Add(time.Hour * -48)
		http.SetCookie(w, cookie)
	}
	w.Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s.%s/login/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_MEMBER_PREFIX,
		this.Context.Config().Get(variable.ServerDomain),
	)))
}
