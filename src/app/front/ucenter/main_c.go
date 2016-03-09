/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
<<<<<<< HEAD
	gfmt "github.com/jsix/gof/util/fmt"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
=======
	"github.com/jsix/gof"
	gfmt "github.com/jsix/gof/util/fmt"
	"github.com/jsix/gof/web"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"time"
)

type mainC struct {
<<<<<<< HEAD
}

//todo:bug 当在UCenter登陆，会话会超时
func (this *mainC) Index(ctx *echox.Context) error {
	mm := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)

	acc := dps.MemberService.GetAccount(mm.Id)
	js, _ := json.Marshal(mm)
	info := make(map[string]string)
	info["memName"] = mm.Name

	lv := dps.PartnerService.GetLevel(p.Id, mm.Level)
	//nextLv := dps.PartnerService.GetNextLevel(p.Id, mm.Level)

	//		if nextLv == nil {
	//			nextLv = lv
	//		}

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"level": lv,
		//"nLevel":       nextLv,
		"member":       mm,
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"json":         template.JS(js),
		"acc":          acc,
		"regTime":      time.Unix(mm.RegTime, 0).Format("2006-01-02"),
		"name": template.HTML(gfmt.BoolString(len(mm.Name) == 0, `<span class="red">未填写</span>`,
			mm.Name)),
		"sex": gfmt.BoolString(mm.Sex == 1, "先生",
			gfmt.BoolString(mm.Sex == 2, "女士", "")),
	}

	return ctx.RenderOK("index.html", d)
}

func (this *mainC) Logout(ctx *echox.Context) error {
	//	r, w := ctx.Request, ctx.Response()
	//	cookie, err := r.Cookie("ms_token")
	//	if err == nil {
	//		cookie.Expires = time.Now().Add(time.Hour * -48)
	//		http.SetCookie(w, cookie)
	//	}
	ctx.Session.Destroy()
	ctx.Response().Write([]byte("<script>location.replace('/login')</script>"))
	return nil
}

// 切换设备
func (this *mainC) Change_device(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.Response(), ctx.Request())
	toUrl := form.Get("return_url")
	if len(toUrl) == 0 {
		toUrl = ctx.Request().Referer()
=======
	base *baseC
}

//todo:bug 当在UCenter登陆，会话会超时
func (this *mainC) Index(ctx *web.Context) {
	if this.base.Requesting(ctx) {
		mm := this.base.GetMember(ctx)
		p := this.base.GetPartner(ctx)

		conf := this.base.GetSiteConf(p.Id)

		acc := dps.MemberService.GetAccount(mm.Id)
		js, _ := json.Marshal(mm)
		info := make(map[string]string)
		info["memName"] = mm.Name

		lv := dps.PartnerService.GetLevel(p.Id, mm.Level)
		//nextLv := dps.PartnerService.GetNextLevel(p.Id, mm.Level)

		//		if nextLv == nil {
		//			nextLv = lv
		//		}

		this.base.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"level": lv,
			//"nLevel":       nextLv,
			"member":       mm,
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"json":         template.JS(js),
			"acc":          acc,
			"regTime":      time.Unix(mm.RegTime, 0).Format("2006-01-02"),
			"name": template.HTML(gfmt.BoolString(len(mm.Name) == 0, `<span class="red">未填写</span>`,
				mm.Name)),
			"sex": gfmt.BoolString(mm.Sex == 1, "先生",
				gfmt.BoolString(mm.Sex == 2, "女士", "")),
		}, "views/ucenter/{device}/index.html",
			"views/ucenter/{device}/inc/header.html",
			"views/ucenter/{device}/inc/menu.html",
			"views/ucenter/{device}/inc/footer.html")
	}
}

func (this *mainC) Logout(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	cookie, err := r.Cookie("ms_token")
	if err == nil {
		cookie.Expires = time.Now().Add(time.Hour * -48)
		http.SetCookie(w, cookie)
	}
	w.Write([]byte("<script>location.replace('/login')</script>"))
}

// 切换设备
func (this *mainC) Change_device(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	toUrl := form.Get("return_url")
	if len(toUrl) == 0 {
		toUrl = ctx.Request.Referer()
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		if len(toUrl) == 0 {
			toUrl = "/"
		}
	}
<<<<<<< HEAD
	ctx.Response().Header().Add("Location", toUrl)
	ctx.Response().WriteHeader(302)
	return nil
}

// Member session connect
func (this *mainC) Msc(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.Response(), ctx.Request())
	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session.Get("member")
=======

	ctx.Response.Header().Add("Location", toUrl)
	ctx.Response.WriteHeader(302)
}

// Member session connect
func (this *mainC) Msc(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session().Get("member")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		var m *member.ValueMember
		if v != nil {
			m = v.(*member.ValueMember)
			if m.Id != memberId { // 如果会话冲突
				m = nil
			}
		}
<<<<<<< HEAD
		if m == nil {
			m = dps.MemberService.GetMember(memberId)
			ctx.Session.Set("member", m)
			ctx.Session.Save()
		}
	})
	if ok {
		ctx.Session.Set("client_member_id", memberId)
		ctx.Session.Save()
	}
=======

		if m == nil {
			m = dps.MemberService.GetMember(memberId)
			ctx.Session().Set("member", m)
			ctx.Session().Save()
		}
	})

	if ok {
		ctx.Items["client_member_id"] = memberId
	}

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
<<<<<<< HEAD
	ctx.Response().Header().Add("Location", rtu)
	ctx.Response().WriteHeader(302)
	return nil
}

// Member session disconnect
func (this *mainC) Msd(ctx *echox.Context) error {
	if util.MemberHttpSessionDisconnect(ctx) {
		ctx.Session.Set("member", nil)
		ctx.Session.Save()
		return ctx.StringOK("disconnect success")
	}
	return ctx.StringOK("disconnect fail")
=======
	ctx.Response.Header().Add("Location", rtu)
	ctx.Response.WriteHeader(302)
}

// Member session disconnect
func (this *mainC) Msd(ctx *web.Context) {
	if util.MemberHttpSessionDisconnect(ctx) {
		ctx.Session().Set("member", nil)
		ctx.Session().Save()
		ctx.Response.Write([]byte("disconnect success"))
	} else {
		ctx.Response.Write([]byte("disconnect fail"))
	}
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
