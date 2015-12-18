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
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
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

	s.Static("/static/", "./public/static/") //静态资源
	s.Getx("/", mc.Index)
	s.Getx("/cart", cc.Index)
	s.Getx("/change_device", mc.change_device)
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
	s.Getx("/goods-describe", lc.GoodsDetails)
	s.Getx("/st/*", lc.SaleTagGoodsList)
	s.Getx("/user/jump_m", uc.JumpToMCenter)
	s.Getx("/c-*.htm", lc.List_Index)
	s.Getx("/goods-*.htm", lc.GoodsView)
}

func getServe(path string) *echox.Echo {
	s := echox.New()
	s.SetTemplateRender(path)
	s.Use(mw.Recover())
	s.Use(shopCheck)
	registerRoutes(s)
	return s
}

// 初始化
func init() {
	pcServe = getServe("public/views/shop/ols/pc")
	mobiServe = getServe("public/views/shop/ols/mobi")
	embedServe = getServe("public/views/shop/ols/app_embed")
}

// 获取所有服务
func GetServes() (sPc *echox.Echo, sMobi *echox.Echo, sApp *echox.Echo) {
	return pcServe, mobiServe, embedServe
}

// 处理服务
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func shopCheck(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		// 商户不存在
		s := session.Default(ctx.Response(), ctx.Request())
		partnerId := GetPartnerId(ctx.Request(), s)
		if partnerId <= 0 {
			return ctx.String(http.StatusOK, "No such partner")
		}

		ctx.Set("partner_id", partnerId) // 缓存PartnerId

		// 判断线上商店开通情况
		var conf = cache.GetPartnerSiteConf(partnerId)
		if conf == nil {
			return ctx.String(http.StatusOK, "线上商店未开通")
		}

		if conf.State == enum.PARTNER_SITE_CLOSED {
			if strings.TrimSpace(conf.StateHtml) == "" {
				conf.StateHtml = "网站暂停访问"
			}
			return ctx.String(http.StatusOK, conf.StateHtml)
		}
		return h(ctx)
	}
}
