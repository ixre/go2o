/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package www

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"go2o/app/cache/apicache"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/infrastructure/format"
	"go2o/core/service/goclient"
	"go2o/share/variable"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type mainC struct {
	app.Context
}

// 处理跳转
func (this *mainC) HandleIndexGo(w http.ResponseWriter, r *http.Request) bool {
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

func (this *mainC) Index(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	if this.HandleIndexGo(w, r) {
		return
	}

	if b, siteConf := GetSiteConf(w, p); b {
		shops := apicache.GetShops(this.Context, p.Id, p.Secret)
		if shops == nil {
			shops = []byte("{}")
		}
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

func (this *mainC) Login(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner, mm *member.ValueMember) {
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	if b, siteConf := GetSiteConf(w, p); b {
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			mv := *m
			mv["partner"] = p
			mv["title"] = "会员登录－" + siteConf.SubTitle
			mv["member"] = mm
			mv["conf"] = siteConf
			mv["tipStyle"] = tipStyle
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

func (this *mainC) Register(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
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

func (this *mainC) ValidUsr_post(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	r.ParseForm()
	usr := r.FormValue("usr")
	b := goclient.Partner.UserIsExist(p.Id, p.Secret, usr)
	if !b {
		w.Write([]byte(`{"result":true}`))
	} else {
		w.Write([]byte(`{"result":false}`))
	}
}

func (this *mainC) PostRegistInfo_post(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	r.ParseForm()
	var member member.ValueMember
	web.ParseFormToEntity(r.Form, &member)
	if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
		member.RegIp = r.RemoteAddr[:i]
	}
	b, err := goclient.Partner.RegisterMember(&member, p.Id, 0, "")
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
func (this *mainC) Member(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner, mm *member.ValueMember) {
	var location string
	if mm == nil {
		location = "/login?return_url=/member"
	} else {
		cookie, _ := r.Cookie("ms_token")
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?token=%s",
			variable.DOMAIN_MEMBER_PREFIX,
			this.Context.Config().GetString(variable.ServerDomain),
			cookie.Value,
		)
	}
	w.Write([]byte("<script>window.parent.location.replace('" + location + "')</script>"))
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
		this.Context.Config().GetString(variable.ServerDomain),
	)))
}

func (this *mainC) List(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner, mm *member.ValueMember) {
	if b, siteConf := GetSiteConf(w, p); b {
		categories := apicache.GetCategories(this.Context, p.Id, p.Secret)
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "在线订餐-" + p.Name
			(*m)["categories"] = template.HTML(categories)
			(*m)["member"] = mm
			(*m)["conf"] = siteConf
		},
			"views/web/www/list.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) GetList(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	const getNum int = -1 //-1表示全部
	categoryId, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		w.Write([]byte(`{"error":"yes"}`))
		return
	}
	items, err := goclient.Partner.GetItems(p.Id, p.Secret, categoryId, getNum)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	buf := bytes.NewBufferString("<ul>")

	for _, v := range items {

		buf.WriteString(fmt.Sprintf(`
			<li>
				<div class="gs_goodss">
                        <img src="%s" alt="%s"/>
                        <h3 class="name">%s%s</h3>
                        <span class="srice">原价:￥%s</span>
                        <span class="sprice">优惠价:￥%s</span>
                        <a href="javascript:cart.add(%d,1);" class="add">&nbsp;</a>
                </div>
             </li>
		`, format.GetGoodsImageUrl(v.Image), v.Name, v.Name, v.SmallTitle, format.FormatFloat(v.Price),
			format.FormatFloat(v.SalePrice),
			v.Id))
	}
	buf.WriteString("</ul>")
	w.Write(buf.Bytes())
}
