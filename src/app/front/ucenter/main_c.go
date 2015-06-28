/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"github.com/atnet/gof"
	gfmt "github.com/atnet/gof/util/fmt"
	"github.com/atnet/gof/web"
	"go2o/src/app/util"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"time"
	"fmt"
)

type mainC struct {
	*baseC
}

//todo:bug 当在UCenter登陆，会话会超时
func (this *mainC) Index(ctx *web.Context) {
	if this.Requesting(ctx) {
		mm := this.GetMember(ctx)
		p := this.GetPartner(ctx)


		fmt.Printf("--%+v   --- %+v\n",mm,p)
		conf := this.GetSiteConf(p.Id)


		acc := dps.MemberService.GetAccount(mm.Id)
		js, _ := json.Marshal(mm)
		info := make(map[string]string)
		info["memName"] = mm.Name

		lv := dps.PartnerService.GetLevel(p.Id, mm.Level)
		//nextLv := dps.PartnerService.GetNextLevel(p.Id, mm.Level)

		//		if nextLv == nil {
		//			nextLv = lv
		//		}

		this.ExecuteTemplate(ctx, gof.TemplateDataMap{
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
	r, w := ctx.Request, ctx.ResponseWriter
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
		if len(toUrl) == 0 {
			toUrl = "/"
		}
	}

	ctx.ResponseWriter.Header().Add("Location", toUrl)
	ctx.ResponseWriter.WriteHeader(302)
}

// Member session connect
func (this *mainC) Msc(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		if ctx.Session().Get("member") == nil {
			m := dps.MemberService.GetMember(memberId)
			ctx.Session().Set("member", m)
		}
	})

	if ok {
		ctx.Items["client_member_id"] = memberId
	}

	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
	ctx.ResponseWriter.Header().Add("Location", rtu)
	ctx.ResponseWriter.WriteHeader(302)
}

// Member session disconnect
func (this *mainC) Msd(ctx *web.Context) {
	if util.MemberHttpSessionDisconnect(ctx) {
		ctx.Session().Set("member", nil)
		ctx.Session().Save()
		ctx.ResponseWriter.Write([]byte("disconnect success"))
	} else {
		ctx.ResponseWriter.Write([]byte("disconnect fail"))
	}
}
