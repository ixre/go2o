/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package www

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/web"
	"go2o/src/app/cache/apicache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"go2o/src/app/front"
)

type mainC struct {
	front.WebC
}

// 处理跳转
func (this *mainC) HandleIndexGo(ctx *web.Context) bool {
	r, w := ctx.Request, ctx.ResponseWriter
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

func (this *mainC) Index(ctx *web.Context, p *partner.ValuePartner) {
	_, w := ctx.Request, ctx.ResponseWriter
	if this.HandleIndexGo(ctx) {
		return
	}

	if b, siteConf := GetSiteConf(w, p); b {
		shops := apicache.GetShops(ctx.App, p.Id, p.Secret)
		if shops == nil {
			shops = []byte("{}")
		}
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
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

func (this *mainC) Login(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	p,_ := this.WebC.GetPartner(ctx)
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	if b, siteConf := GetSiteConf(w, p); b {
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
			mv := *m
			mv["partner"] = p
			mv["title"] = "会员登录－" + siteConf.SubTitle
			mv["conf"] = siteConf
			mv["tipStyle"] = tipStyle
		},
			"views/web/www/login.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	result,_ := goclient.Member.Login(usr, pwd)
	if result.Result {
		ctx.Session().Set("member", result.Member)
		ctx.Session().Save()
		fmt.Println("+=====",result.Member)
		w.Write([]byte("{result:true}"))
	}
	w.Write([]byte("{result:false,message:'" + result.Message + "'}"))
}

func (this *mainC) Register(ctx *web.Context, p *partner.ValuePartner) {
	_, w := ctx.Request, ctx.ResponseWriter
	if b, siteConf := GetSiteConf(w, p); b {
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "会员注册－" + siteConf.SubTitle
			(*m)["conf"] = siteConf
		},
			"views/web/www/register.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *mainC) ValidUsr_post(ctx *web.Context, p *partner.ValuePartner) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	usr := r.FormValue("usr")
	b := goclient.Partner.UserIsExist(p.Id, p.Secret, usr)
	if !b {
		w.Write([]byte(`{"result":true}`))
	} else {
		w.Write([]byte(`{"result":false}`))
	}
}

func (this *mainC) PostRegistInfo_post(ctx *web.Context, p *partner.ValuePartner) {
	r, w := ctx.Request, ctx.ResponseWriter
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
func (this *mainC) Member(ctx *web.Context, p *partner.ValuePartner, mm *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	var location string
	if mm == nil {
		location = "/login?return_url=/member"
	} else {
		cookie, _ := r.Cookie("ms_token")
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?token=%s",
			variable.DOMAIN_MEMBER_PREFIX,
			ctx.App.Config().GetString(variable.ServerDomain),
			cookie.Value,
		)
	}
	w.Write([]byte("<script>window.parent.location.replace('" + location + "')</script>"))
}

//退出
func (this *mainC) Logout(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
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
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
}

func (this *mainC) List(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	p, _ := this.WebC.GetPartner(ctx)
	mm := this.WebC.GetMember(ctx)
	if b, siteConf := GetSiteConf(w, p); b {
		categories := apicache.GetCategories(ctx.App, p.Id, p.Secret)
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
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

func (this *mainC) GetList(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	p, _ := this.WebC.GetPartner(ctx)
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
