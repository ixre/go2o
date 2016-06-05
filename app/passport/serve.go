/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package passport

import (
	"go2o/app/util"
	"go2o/x/echox"
	//mw "gopkg.in/labstack/echo.v1/middleware"
	"net/http"
	"gopkg.in/labstack/echo.v1"
	"github.com/jsix/gof"
)


//注册路由
func registerRoutes(s *echox.Echo) {
	app := gof.CurrentApp
	mc := &mainC{App:app}
	sc := &ssoC{App:app}
	s.Getx("/",mc.Index)
	s.Getx("/sso/login",sc.Login)
	s.Static("/static/", "./public/static/") //静态资源
	s.Static("/image/", "./uploads/")        //图片资源
}

func getServe(path string, files ...string) *echox.Echo {
	s := echox.New()
	//s.Use(mw.Recover())
	s.Use(echox.StopAttack)
	s.Use(beforeHanding)
	registerRoutes(s)
	s.SetTemplateRender(path, files...)
	return s
}

var _ http.Handler = new(PassportServe)
type PassportServe struct{
	pcServe    *echox.Echo
	mobiServe  *echox.Echo
}

func NewServe()http.Handler {
	return &PassportServe{
		pcServe : getServe("public/views/passport/pc/",
			"include/header.html",
			"include/footer.html"),
		mobiServe : getServe("public/views/passport/mobi/",
			"include/header.html",
			"include/footer.html"),
	}
}

// 处理服务
func (this *PassportServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch util.GetBrownerDevice(r) {
	default:
	case util.DevicePC:
		this.pcServe.ServeHTTP(w, r)
	case util.DeviceTouchPad, util.DeviceMobile:
		this.mobiServe.ServeHTTP(w, r)
	}
}


func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		b := util.GetBrownerDevice(ctx.Request())
		if b != util.DeviceMobile{
			util.SetBrownerDevice(ctx.Response(),ctx.Request(),util.DeviceMobile)
		}
		return h(ctx)
	}
}

