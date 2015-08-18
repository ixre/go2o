/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
)

type loginC struct {
	gof.App
}

//登陆
func (this *loginC) Login(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/login.html")
}
func (this *loginC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")
	pt, result, message := this.ValidLogin(usr, pwd)

	if result {
		ctx.Session().Set("partner_id", pt.Id)
		if err := ctx.Session().Save(); err != nil {
			result = false
			message = err.Error()
		}
	}

	if result {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte("{result:false,message:'" + message + "'}"))
	}
}

//验证登陆
func (pb *loginC) ValidLogin(usr string, pwd string) (*partner.ValuePartner, bool, string) {
	var message string
	var result bool
	var pt *partner.ValuePartner
	var err error

	id := dps.PartnerService.Verify(usr, pwd)

	if id == -1 {
		result = false
		message = "用户或密码不正确！"
	} else {
		pt, err = dps.PartnerService.GetPartner(id)
		if err != nil {
			message = err.Error()
			result = false
		} else {
			result = true
		}
	}
	return pt, result, message
}
