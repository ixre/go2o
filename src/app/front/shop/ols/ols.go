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
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/url"
	"sync"
)

var (
	mux sync.Mutex
)

// 获取商户编号
func GetSessionMerchantId(ctx *echox.Context) int {
	return ctx.Get("merchant_id").(int)
}

func getPartner(ctx *echox.Context) *merchant.MerchantValue {
	merchantId := ctx.Get("merchant_id").(int)
	return cache.GetValuePartnerCache(merchantId)
}

// 获取商户API信息
func getPartnerApi(ctx *echox.Context) *merchant.ApiInfo {
	return dps.PartnerService.GetApiInfo(GetSessionMerchantId(ctx))
}

// 获取商户站点设置
func getSiteConf(ctx *echox.Context) *merchant.SiteConf {
	conf := ctx.Get("conf").(*merchant.SiteConf)
	if conf == nil {
		conf = cache.GetPartnerSiteConf(GetSessionMerchantId(ctx))
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
func GetMerchantId(ctx *echox.Context) int {
	mux.Lock()
	defer mux.Unlock()
	if v := ctx.Get("merchant_id"); v != nil {
		return v.(int)
	}
	currHost := ctx.Request().Host
	//ctx.Set("webui_host", currHost)
	merchantId := cache.GetMerchantIdByHost(currHost)
	ctx.Set("merchantId", merchantId)
	return merchantId
}
