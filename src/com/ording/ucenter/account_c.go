package ucenter

import (
	"com/domain/interface/member"
	"com/ording"
	"com/ording/dao"
	"com/ording/entity"
	"encoding/json"
	"net/http"
	"ops/cf/app"
	"ops/cf/web/pager"
	"strconv"
)

type accountC struct {
	app.Context
}

func (this *accountC) Imcomelog(w http.ResponseWriter, r *http.Request,
	mm *member.ValueMember, p *entity.Partner, conf *entity.SiteConf) {

	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		(*m)["conf"] = conf
		(*m)["partner_host"] = conf.Host
		(*m)["record"] = 15
		(*m)["partner"] = p
		(*m)["member"] = mm
	}, "views/ucenter/account/income_log.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *accountC) Imcomelog_post(w http.ResponseWriter, r *http.Request, m *member.ValueMember) {

	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dao.Member().GetIncomeLog(m.Id, page, size, "", "record_time DESC")

	tpage := n / size
	if n%size != 0 {
		tpage = tpage + 1
	}

	p := pager.NewUrlPager(tpage, page, nil)

	pager := &ording.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}

func (this *accountC) ApplyCash(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember, p *entity.Partner, host string) {
	//acc := dao.G
}

/*
def bank_applycash(self):
        '申请提现'
        memberid=self.member['id']
        bankaccount=getbank(memberid)
        account=getaccount(memberid)           #会员账户
        return TPL_USR.bank_applycash(host=HOST,usrtpl=usrtpl,account=account,bank=bankaccount)

    def bank_update_post(self):
        '更新银行帐号'
        request=web.input()
        upbank(self.member['id'],request.bankname,request.bankaccount)
        return '<script>window.parent.location.reload()</script>'

*/
