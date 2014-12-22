package ucenter

import (
	"com/domain/interface/member"
	"com/ording/dproxy"
	"com/ording/entity"
	"com/service/goclient"
	"encoding/json"
	"html/template"
	"net/http"
	"ops/cf"
	"ops/cf/app"
	"ops/cf/web"
	"time"
)

type mainC struct {
	app.Context
}

func (this *mainC) Login(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/ucenter/login.html", nil)
}

func (this *mainC) Index(w http.ResponseWriter, r *http.Request, mm *member.ValueMember,
	p *entity.Partner, conf *entity.SiteConf) {
	acc, _ := goclient.Member.GetMemberAccount(mm.Id, mm.LoginToken)
	js, _ := json.Marshal(mm)
	info := make(map[string]string)
	info["memName"] = mm.Name

	lv := dproxy.MemberService.GetLevelById(mm.Level)
	nextLv := dproxy.MemberService.GetNextLevel(mm.Level)
	if nextLv == nil {
		nextLv = &lv
	}

	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		mv := *m
		mv["level"] = lv
		mv["nLevel"] = nextLv
		mv["title"] = "会员中心"
		mv["member"] = mm
		mv["partner"] = p
		mv["conf"] = conf
		mv["partner_host"] = conf.Host
		mv["json"] = template.JS(js)
		mv["acc"] = acc
		mv["regTime"] = mm.RegTime.Format("2006-01-02")
		mv["name"] = cf.BoolString(len(mm.Name) == 0,
			`<span class="red">未填写</span>`,
			mm.Name)

		mv["sex"] = cf.BoolString(mm.Sex == 1, "先生",
			cf.BoolString(mm.Sex == 2, "女士", ""))

	}, "views/ucenter/index.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")

}

func (this *mainC) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ms_token")
	if err == nil {
		cookie.Expires = time.Now().Add(time.Hour * -48)
		http.SetCookie(w, cookie)
	}
	w.Write([]byte("<script>location.replace('/login')</script>"))
}

func (this *mainC) Profile(w http.ResponseWriter, r *http.Request, mm *member.ValueMember,
	p *entity.Partner, conf *entity.SiteConf) {
	js, _ := json.Marshal(mm)

	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		(*m)["partner"] = p
		(*m)["conf"] = conf
		(*m)["partner_host"] = conf.Host
		(*m)["member"] = mm
		(*m)["entity"] = template.JS(js)

	}, "views/ucenter/profile.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *mainC) Profile_post(w http.ResponseWriter, r *http.Request, mm *member.ValueMember) {
	var result cf.JsonResult
	r.ParseForm()
	clientM := new(member.ValueMember)
	web.ParseFormToEntity(r.Form, clientM)
	clientM.Id = mm.Id
	_, err := goclient.Member.SaveMember(clientM, mm.LoginToken)

	if err != nil {
		result = cf.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true}
	}
	w.Write(result.Marshal())
}
