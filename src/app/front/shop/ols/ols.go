/**
 * Copyright 2015 @ z3q.net.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/url"
	"sync"
)

var (
	mux sync.Mutex
)

// 获取商户编号
func GetSessionPartnerId(ctx *echox.Context) int {
	return ctx.Get("partner_id").(int)
}

func getPartner(ctx *echox.Context) *partner.ValuePartner {
	partnerId := ctx.Get("partner_id").(int)
	return cache.GetValuePartnerCache(partnerId)
}

// 获取商户API信息
func getPartnerApi(ctx *echox.Context) *partner.ApiInfo {
	return dps.PartnerService.GetApiInfo(GetSessionPartnerId(ctx))
}

// 获取商户站点设置
func getSiteConf(ctx *echox.Context) *partner.SiteConf {
	conf := ctx.Get("conf").(*partner.SiteConf)
	if conf == nil {
		conf = cache.GetPartnerSiteConf(GetSessionPartnerId(ctx))
		ctx.Set("conf", conf)
	}
	return conf
}

// 获取会员
func GetMember(ctx *echox.Context) *member.ValueMember {
	memberIdObj := ctx.Session.Get("member")
	if memberIdObj != nil {
		if o, ok := memberIdObj.(*member.ValueMember); ok {
			return o
		}
	}
	return nil
}

// 检查会员是否登陆
func CheckMemberLogin(ctx *echox.Context) bool {
	if ctx.Session.Get("member") == nil {
		ctx.Response().Header().Add("Location", "/user/login?return_url="+
			url.QueryEscape(ctx.HttpRequest().RequestURI))
		ctx.Response().WriteHeader(302)
		return false
	}
	return true
}

// 获取商户编号
func GetPartnerId(ctx *echox.Context) int {
	mux.Lock()
	defer mux.Unlock()
	if v := ctx.Get("partner_id"); v != nil {
		return v.(int)
	}
	currHost := ctx.Request().Host
	//ctx.Set("webui_host", currHost)
	partnerId := cache.GetPartnerIdByHost(currHost)
	ctx.Set("partnerId", partnerId)
	return partnerId
}
