/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/core/infrastructure/domain"
	"github.com/atnet/gof"
)

type loginC struct {
}

//登陆
func (this *loginC) Index(ctx *web.Context) {
	ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter,nil,"views/master/login.html")
}
func (this *loginC) Index_post(ctx *web.Context) {
	r := ctx.Request
	var msg gof.Message
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")

	if domain.Md5Pwd(pwd,usr) == ctx.App.Config().GetString("webmaster_valid_md5") {
		ctx.Session().Set("master_id", 1)
		if err := ctx.Session().Save(); err != nil {
			msg.Message = err.Error()
		}else{
			msg.Result =true
		}
	}else{
		msg.Message = "用户或密码不正确！"
	}
	ctx.ResponseWriter.Write(msg.Marshal())
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
