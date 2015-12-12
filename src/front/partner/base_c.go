/**
 * Copyright 2015 @ z3q.net.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import (
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"net/url"
)

func chkLogin(ctx *echox.Context)error(b bool, partnerId int) {
	//todo:仅仅做了id的检测，没有判断有效性
	// i, err := session.GetLSession().GetPartnerIdFromCookie(ctx.Request)
	// return err == nil, i
	v := ctx.Session().Get("partner_id")
	if v == nil {
		return false, -1
	}
	return true, v.(int)
}
func redirect(ctx *echox.Context)error{
	r, w := ctx.Request, ctx.Response
	w.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(r.URL.String()) + "'</script>"))
}

var _ mvc.Filter = new(baseC)

type baseC struct {
}

func (this *baseC) Requesting(ctx *echox.Context)errorbool {
	//验证是否登陆
	if b, _ := chkLogin(ctx); !b {
		redirect(ctx)
		return false
	}
	return true
}
func (this *baseC) RequestEnd(ctx *echox.Context)error{
}

// 获取商户编号
func (this *baseC) GetPartnerId(ctx *echox.Context)errorint {
	v := ctx.Session().Get("partner_id")
	if v == nil {
		this.Requesting(ctx)
		return -1
	}
	return v.(int)
}

func (this *baseC) GetPartner(ctx *echox.Context)error(*partner.ValuePartner, error) {
	return dps.PartnerService.GetPartner(this.GetPartnerId(ctx))
}

// 输出错误信息
func (this *baseC) ErrorOutput(ctx *web.Context, err string) {
	ctx.Response.Write([]byte("{error:\"" + err + "\"}"))
}
