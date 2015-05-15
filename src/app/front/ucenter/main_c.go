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
	"github.com/atnet/gof/web"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"net/http"
	"time"
)

type mainC struct {
	*baseC
}

//todo:bug 当在UCenter登陆，会话会超时
func (this *mainC) Index(ctx *web.Context) {
	if this.Requesting(ctx) {
		mm := this.GetMember(ctx)
		p := this.GetPartner(ctx)
		conf := this.GetSiteConf(p.Id)

		acc, _ := goclient.Member.GetMemberAccount(mm.Id, mm.DynamicToken)
		js, _ := json.Marshal(mm)
		info := make(map[string]string)
		info["memName"] = mm.Name

		lv := dps.MemberService.GetLevelById(mm.Level)
		nextLv := dps.MemberService.GetNextLevel(mm.Level)
		if nextLv == nil {
			nextLv = &lv
		}

		ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
			"level":        lv,
			"nLevel":       nextLv,
			"member":       mm,
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"json":         template.JS(js),
			"acc":          acc,
			"regTime":      time.Unix(mm.RegTime, 0).Format("2006-01-02"),
			"name": gof.BoolString(len(mm.Name) == 0,
				`<span class="red">未填写</span>`,
				mm.Name),

			"sex": gof.BoolString(mm.Sex == 1, "先生",
				gof.BoolString(mm.Sex == 2, "女士", "")),
		}, "views/ucenter/index.html",
			"views/ucenter/inc/header.html",
			"views/ucenter/inc/menu.html",
			"views/ucenter/inc/footer.html")
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
