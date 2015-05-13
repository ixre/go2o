/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package pcshop

import (
	"fmt"
	"github.com/atnet/gof/web"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/dto"
	"go2o/src/core/service/dps"
	"net/url"
)

type baseC struct {
}

func (this *baseC) Requesting(ctx *web.Context) bool {
	// 商户不存在
	partnerId := getPartnerId(ctx)
	if partnerId <= 0 {
		renderError(ctx, true, "No such partner.")
		return false
	}
	ctx.Items["partner_id"] = partnerId // 缓存PartnerId

	// 校验商户并缓存
	pt, err := getPartner(ctx)
	if err != nil {
		renderError(ctx, true, err.Error())
		return false
	}
	ctx.Items["partner_ins"] = pt

	return true
}

func (this *baseC) RequestEnd(*web.Context) {
}

// 获取商户编号
func (this *baseC) GetPartnerId(ctx *web.Context) int {
	return ctx.GetItem("partner_id").(int)
}
func (this *baseC) GetPartner(ctx *web.Context) *partner.ValuePartner {
	return ctx.GetItem("partner_ins").(*partner.ValuePartner)
}

// 获取商户API信息
func (this *baseC) GetPartnerApi(ctx *web.Context) *dto.PartnerApiInfo {
	return dps.PartnerService.GetApiInfo(getPartnerId(ctx))
}

// 获取会员
func (this *baseC) GetMember(ctx *web.Context) *member.ValueMember {
	memberIdObj := ctx.Session().Get("member")
	if memberIdObj != nil {
		if o, ok := memberIdObj.(*member.ValueMember); ok {
			return o
		}
	}
	return nil
}

// 检查会员是否登陆
func (this *baseC) CheckMemberLogin(ctx *web.Context) bool {
	if ctx.Session().Get("member") == nil {
		ctx.ResponseWriter.Header().Add("Location", "/user/login?return_url="+
			url.QueryEscape(ctx.Request.RequestURI))
		ctx.ResponseWriter.WriteHeader(302)
		return false
	}
	return true
}

func renderError(ctx *web.Context, simpleError bool, message string) {
	if simpleError {
		const errTpl string = "<html><body><h1 style='color:red'>%s</h1></body></html>"
		ctx.ResponseWriter.Write([]byte(fmt.Sprintf(errTpl, message)))
	} else {
		//todo: 用模板显示错误
	}
}
func getPartnerId(ctx *web.Context) int {
	currHost := ctx.Request.Host
	host := ctx.Session().Get("webui_host")
	pid := ctx.Session().Get("webui_pid")
	if host == nil || pid == nil || host != currHost {
		partnerId := dps.PartnerService.GetPartnerIdByHost(currHost)
		if partnerId != -1 {
			ctx.Session().Set("webui_host", currHost)
			ctx.Session().Set("webui_pid", partnerId)
			ctx.Session().Save()
		}
		return partnerId
	}
	return pid.(int)
}
func getPartner(ctx *web.Context) (*partner.ValuePartner, error) {
	//todo: 缓存，用Member对应的Partner编号来查询缓存
	var partnerId int = ctx.GetItem("partner_id").(int)
	var err error
	var pt *partner.ValuePartner = cache.GetValuePartnerCache(partnerId)
	if pt == nil {
		if pt, err = dps.PartnerService.GetPartner(getPartnerId(ctx)); err == nil {
			cache.SetValuePartnerCache(partnerId, pt)
		}
	}
	return pt, err
}
