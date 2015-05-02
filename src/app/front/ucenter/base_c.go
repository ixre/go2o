/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"net/url"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/goclient"
)

func chkLogin(ctx *web.Context)bool {
	return ctx.Session().Get("member") != nil
}
func redirect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	w.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(r.URL.String()) + "'</script>"))
}

var _ mvc.Filter = new(baseC)

type baseC struct {
}

func (this *baseC) Requesting(ctx *web.Context) bool {
	//验证是否登陆
	if !chkLogin(ctx) {
		redirect(ctx)
		return false
	}
	return true
}
func (this *baseC) RequestEnd(ctx *web.Context) {
}

// 获取商户
func (this *baseC) GetPartner(ctx *web.Context) *partner.ValuePartner{
	return ctx.Session().Get("member:rel_partner").(*partner.ValuePartner)
}

// 获取会员
func (this *baseC) GetMember(ctx *web.Context) *member.ValueMember{
	memberIdObj := ctx.Session().Get("member")
	if memberIdObj != nil{
		if o,ok := memberIdObj.(*member.ValueMember);ok{
			return o
		}
	}
	return nil
}

// 获取商户的站点设置
func (this *baseC) GetSiteConf(partnerId int,secret string)(*partner.SiteConf,error){
	return goclient.Partner.GetSiteConf(partnerId, secret)
}