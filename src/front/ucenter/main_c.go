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
	gfmt "github.com/jsix/gof/util/fmt"
	"github.com/jsix/gof/web"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"time"
)

type mainC struct {
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

	ctx.Response.Header().Add("Location", toUrl)
	ctx.Response.WriteHeader(302)
}

// Member session connect
func (this *mainC) Msc(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session().Get("member")
		var m *member.ValueMember
		if v != nil {
			m = v.(*member.ValueMember)
			if m.Id != memberId { // 如果会话冲突
				m = nil
			}
		}

		if m == nil {
			m = dps.MemberService.GetMember(memberId)
			ctx.Session().Set("member", m)
			ctx.Session().Save()
		}
	})

	if ok {
		ctx.Items["client_member_id"] = memberId
	}

	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
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
}
