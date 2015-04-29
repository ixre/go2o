/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type mainC struct {
	gof.App
}

func (this *mainC) Login(ctx *web.Context) {
	this.App.Template().Render(ctx.ResponseWriter, "views/ucenter/login.html", nil)
}

func (this *mainC) Index(ctx *web.Context, mm *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {
	acc, _ := goclient.Member.GetMemberAccount(mm.Id, mm.LoginToken)
	js, _ := json.Marshal(mm)
	info := make(map[string]string)
	info["memName"] = mm.Name

	lv := dps.MemberService.GetLevelById(mm.Level)
	nextLv := dps.MemberService.GetNextLevel(mm.Level)
	if nextLv == nil {
		nextLv = &lv
	}

	this.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		mv := *m
		mv["level"] = lv
		mv["nLevel"] = nextLv
		mv["member"] = mm
		mv["partner"] = p
		mv["conf"] = conf
		mv["partner_host"] = conf.Host
		mv["json"] = template.JS(js)
		mv["acc"] = acc
		mv["regTime"] = time.Unix(mm.RegTime, 0).Format("2006-01-02")
		mv["name"] = gof.BoolString(len(mm.Name) == 0,
			`<span class="red">未填写</span>`,
			mm.Name)

		mv["sex"] = gof.BoolString(mm.Sex == 1, "先生",
			gof.BoolString(mm.Sex == 2, "女士", ""))

	}, "views/ucenter/index.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")

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

func (this *mainC) Profile(ctx *web.Context, mm *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {
	js, _ := json.Marshal(mm)

	this.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		v := *m
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = mm
		v["entity"] = template.JS(js)

	}, "views/ucenter/profile.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *mainC) Pwd(ctx *web.Context, mm *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {

	this.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		v := *m
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = mm

	}, "views/ucenter/pwd.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *mainC) Pwd_post(ctx *web.Context, m *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()
	var oldPwd, newPwd, rePwd string
	oldPwd = r.FormValue("OldPwd")
	newPwd = r.FormValue("NewPwd")
	rePwd = r.FormValue("RePwd")
	var err error
	if newPwd != rePwd {
		err = errors.New("两次密码输入不一致")
	} else {
		err = dps.MemberService.ModifyPassword(m.Id, oldPwd, newPwd)
	}
	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true}
	}
	w.Write(result.Marshal())
}
func (this *mainC) Profile_post(ctx *web.Context, mm *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()
	clientM := new(member.ValueMember)
	web.ParseFormToEntity(r.Form, clientM)
	clientM.Id = mm.Id
	_, err := goclient.Member.SaveMember(clientM, mm.LoginToken)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true}
	}
	w.Write(result.Marshal())
}

func (this *mainC) Deliver(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	this.App.Template().Execute(ctx.ResponseWriter, func(mp *map[string]interface{}) {
		v := *mp
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = m

	}, "views/ucenter/deliver.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *mainC) Deliver_post(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {
	addrs, err := goclient.Member.GetDeliverAddrs(m.Id, m.LoginToken)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("{error:'错误:" + err.Error() + "'}"))
		return
	}
	js, _ := json.Marshal(addrs)
	ctx.ResponseWriter.Write([]byte(`{"rows":` + string(js) + `}`))
}

func (this *mainC) SaveDeliver_post(w http.ResponseWriter,
	r *http.Request, m *member.ValueMember, p *partner.ValuePartner) {
	r.ParseForm()
	var e member.DeliverAddress
	web.ParseFormToEntity(r.Form, &e)
	e.MemberId = m.Id
	b, err := goclient.Member.SaveDeliverAddr(m.Id, m.LoginToken, &e)
	if err == nil {
		if b {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}

func (this *mainC) DeleteDeliver_post(w http.ResponseWriter,
	r *http.Request, m *member.ValueMember) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))

	b, err := goclient.Member.DeleteDeliverAddr(m.Id, m.LoginToken, id)
	if err == nil {
		if b {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}
