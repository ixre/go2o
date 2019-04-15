/**
 * Copyright 2015 @ z3q.net.
 * name : serve.go
 * author : jarryliu
 * date : 2016-09-09 17:33
 * description :
 * history :
 */
package hapi

import (
	"github.com/ixre/goex/echox"
	"github.com/ixre/gof"
	"github.com/labstack/echo"
	"go2o/app/web/shared"
	"net/http"
)

func NewServe() *echox.Echo {
	s := echox.New()
	//s.Use(mw.Recover())
	// s.Use(echox.StopAttack)
	s.Use(beforeHanding)
	registerRoutes(s)
	return s
}

func registerRoutes(s *echox.Echo) {
	app := gof.CurrentApp
	mc := &mainC{app}
	us := &shared.UserSync{}
	sc := &serviceC{app}
	pc := &presentC{}
	s.GET("/api_info", mc.Info)
	s.GET("/test", mc.Test)
	s.GET("/request_login", mc.RequestLogin)
	s.GET("/r/uc", mc.RedirectUc)
	s.GET("/user/sync_m.p", us.Sync) //同步登录登出
	s.GET("/gad_api", pc.AdApi)
	s.AutoGET("/service", sc)          //服务
	s.AutoGET("/!s", sc)               //服务
	s.AutoGET("/!sp", &shoppingC{app}) //购物
}

func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//todo: 同源验证,受信任的源
		return h(c)
	}
}

//检查会员编号
func getMemberId(c *echox.Context) int {
	return shared.GetMemberId(c)
}

func requestLogin(c *echox.Context) error {
	msg := gof.Result{
		ErrCode: 1,
		ErrMsg:  "not login",
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), msg)
}
