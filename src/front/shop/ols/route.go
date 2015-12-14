/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ols

import (
	mw "github.com/labstack/echo/middleware"
	"go2o/src/app/util"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/x/echox"
	"net/http"
	"strings"
)

var (
	waitInit   bool = true
	pcServe    *echox.Echo
	mobiServe  *echox.Echo
	embedServe *echox.Echo
)

//注册路由
func registerRoutes(s *echox.Echo) {
	mc := &MainC{}
	sp := &ShoppingC{}
	pc := &PaymentC{}
	cc := &CartC{}
	uc := &UserC{}
	lc := &ListC{}
	ctc := &ContentC{}

	s.Danyx("/main/:action", mc)
	s.Danyx("/buy/:action", sp)

	s.Danyx("/shopping/:action", sp)
	s.Danyx("/list/:action", lc)
	s.Danyx("/cart/:action", cc)
	s.Danyx("/user/:action", uc)
	s.Danyx("/content/:action", ctc)
	s.Danyx("/pay/:action", pc)

	// 购物车接口
	s.Postx("/cart_api_v1", cc.CartApiHandle)

	// 支付异步提醒
	s.Postx("/pay/notify/*", pc.Notify_post)

	// 首页
	s.Getx("/", mc.Index)
	s.Getx("/goods-describe", lc.GoodsDetails)
	s.Getx("/st/*", lc.SaleTagGoodsList)
	s.Getx("/user/jump_m", uc.JumpToMCenter)
	s.Getx("^/c-*.htm", lc.List_Index)
	s.Getx("^/goods-*.htm", lc.GoodsView)
}

func getServe(path string) *echox.Echo {
	s := echox.New()
	s.SetTemplateRender(path)
	s.Use(mw.Recover())
	s.Use(shopCheck)
	registerRoutes(s)
	return s
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if waitInit {
		pcServe = getServe("public/views/shop/ols/pc")
		mobiServe = getServe("public/views/shop/ols/mobi")
		embedServe = getServe("public/views/shop/ols/app_embed")
		waitInit = true
	}
	switch util.GetBrownerDevice(r) {
	default:
	case util.DevicePC:
		pcServe.ServeHTTP(w, r)
	case util.DeviceTouchPad, util.DeviceMobile:
		mobiServe.ServeHTTP(w, r)
	case util.DeviceAppEmbed:
		embedServe.ServeHTTP(w, r)
	}
}

func shopCheck(ctx *echox.Context) error {
	// 商户不存在
	partnerId := getPartnerId(ctx)
	if partnerId <= 0 {
		err := ctx.StringOK("No such partner.")
		ctx.Done()
		return err
	}

	ctx.Set("partner_id", partnerId) // 缓存PartnerId

	// 判断线上商店开通情况
	var conf = cache.GetPartnerSiteConf(partnerId)
	if conf == nil {
		err := ctx.StringOK("线上商店未开通")
		ctx.Done()
		return err
	}

	if conf.State == enum.PARTNER_SITE_CLOSED {
		if strings.TrimSpace(conf.StateHtml) == "" {
			conf.StateHtml = "网站暂停访问"
		}
		err := ctx.StringOK(conf.StateHtml)
		ctx.Done()
		return err
	}

	return true
}
