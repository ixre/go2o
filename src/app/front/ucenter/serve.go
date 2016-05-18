/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"go2o/src/app/util"
	"go2o/src/x/echox"
	mw "gopkg.in/labstack/echo.v1/middleware"
	"net/http"
	"sync"
)

var (
	waitInit   bool = true
	serveMux   sync.Mutex
	pcServe    *echox.Echo
	mobiServe  *echox.Echo
	embedServe *echox.Echo
)

func registerRoutes(s *echox.Echo) {
	mc := &mainC{}
	bc := &basicC{}
	oc := &orderC{}
	ac := &accountC{}
	lc := &loginC{}
	gc := &getC{}
	riseC := &personFinanceRiseC{}
	jc := &jsonC{}

	s.Static("/static/", "./public/static/") //静态资源
	s.Getx("/", mc.Index)
	s.Getx("/logout", mc.Logout)
	s.Anyx("/login", lc.Index)
	s.Getx("/change_device", mc.Change_device)
	s.Getx("/msc", mc.Msc)
	s.Getx("/msd", mc.Msd)
	s.Getx("/partner_connect", lc.Merchant_connect)
	s.Getx("/partner_disconnect", lc.Merchant_disconnect)
	s.Aanyx("/basic/:action", bc)
	s.Aanyx("/order/:action", oc)
	s.Aanyx("/account/:action", ac)
	s.Getx("/get/qr/:code/:size", gc.GetQR)
	s.Aanyx("/finance/rise/:action", riseC)
	s.Apostx("/json/:action", jc)
}

func getServe(path string) *echox.Echo {
	s := echox.New()
	s.SetTemplateRender(path)
	s.Use(mw.Recover())
	s.Use(echox.StopAttack)
	s.Use(memberLogonCheck)
	registerRoutes(s)
	return s
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if waitInit {
		serveMux.Lock()
		defer serveMux.Unlock()
		pcServe = getServe("public/views/ucenter/pc")
		mobiServe = getServe("public/views/ucenter/mobi")
		embedServe = getServe("public/views/ucenter/mobi")
		//embedServe = getServe("public/views/ucenter/app_embed")
		waitInit = false
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
