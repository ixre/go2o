/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"net/url"
)

func chkLogin(ctx *web.Context) (b bool, partnerId int) {
	//todo:仅仅做了id的检测，没有判断有效性
	// i, err := session.GetLSession().GetPartnerIdFromCookie(ctx.Request)
	// return err == nil, i
	v := ctx.Session().Get("partner_id")
	if v == nil {
		return false, -1
	}
	return true, v.(int)
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
	if b, _ := chkLogin(ctx); !b {
		redirect(ctx)
		return false
	}
	return true
}
func (this *baseC) RequestEnd(ctx *web.Context) {
}

// 获取商户编号
func (this *baseC) GetPartnerId(ctx *web.Context) int {
	v := ctx.Session().Get("partner_id")
	if v == nil {
		this.Requesting(ctx)
		return -1
	}
	return v.(int)
}

func (this *baseC) GetPartner(ctx *web.Context) (*partner.ValuePartner, error) {
	return dps.PartnerService.GetPartner(this.GetPartnerId(ctx))
}

// 输出Json
func (this *baseC) jsonOutput(ctx *web.Context, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		this.errorOutput(ctx, err.Error())
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

// 输出错误信息
func (this *baseC) errorOutput(ctx *web.Context, err string) {
	ctx.ResponseWriter.Write([]byte("{error:\"" + err + "\"}"))
}

// 输出错误信息
func (this *baseC) resultOutput(ctx *web.Context, result gof.Message) {
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf(
		"{result:%v,code:%d,message:\"%s\"}", result.Result, result.Code, result.Message)))
}
