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
	"go2o/src/app/util"
	aputil "go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"time"
)

type mainC struct {
}

func (this *mainC) mobileIndex(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"AliasGrowAccount":    template.HTML(variable.AliasGrowAccount),
		"AliasPresentAccount": template.HTML(variable.AliasPresentAccount),
	}
	return ctx.RenderOK("index.html", d)
}

func (this *mainC) Index(ctx *echox.Context) error {

	switch aputil.GetBrownerDevice(ctx.HttpRequest()) {
	default:
	case aputil.DevicePC:
	case aputil.DeviceTouchPad, aputil.DeviceMobile, aputil.DeviceAppEmbed:
		return this.mobileIndex(ctx)
	}

	mm := getMember(ctx)
	p := getMerchant(ctx)
	conf := getSiteConf(p.Id)

	acc := dps.MemberService.GetAccount(mm.Id)
	js, _ := json.Marshal(mm)
	info := make(map[string]string)
	info["memName"] = mm.Name

	lv := dps.MerchantService.GetLevel(p.Id, mm.Level)
	//nextLv := dps.MerchantService.GetNextLevel(p.Id, mm.Level)

	//		if nextLv == nil {
	//			nextLv = lv
	//		}

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"AliasGrowAccount":    variable.AliasGrowAccount,
		"AliasPresentAccount": variable.AliasPresentAccount,
		"level":               lv,
		//"nLevel":       nextLv,
		"member":       mm,
		"Merchant":     p,
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
	ctx.HttpResponse().Write([]byte("<script>location.replace('/login')</script>"))
	return nil
}

// 切换设备
func (this *mainC) Change_device(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.HttpResponse(), ctx.HttpRequest())
	toUrl := form.Get("return_url")
	if len(toUrl) == 0 {
		toUrl = ctx.Request().Referer()
		if len(toUrl) == 0 {
			toUrl = "/"
		}
	}
	ctx.Response().Header().Add("Location", toUrl)
	ctx.Response().WriteHeader(302)
	return nil
}

// Member session connect
func (this *mainC) Msc(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.HttpResponse(), ctx.HttpRequest())
	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session.Get("member")
		var m *member.ValueMember
		if v != nil {
			m = v.(*member.ValueMember)
			if m.Id != memberId { // 如果会话冲突
				m = nil
			}
		}
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
	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
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
}
