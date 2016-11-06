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
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	"go2o/app/web/shared"
	"go2o/core/domain/interface/member"
	"go2o/x/echox"
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
	uc := &shared.UserC{App: app}
	sc := &serviceC{app}
	s.GET("/api_info", mc.Info)
	s.Getx("/test", mc.Test)
	s.Getx("/request_login", mc.RequestLogin)
	s.Getx("/r/uc", mc.RedirectUc)
	s.Agetx("/user/:action", uc)    //用户
	s.Agetx("/service/:action", sc) //服务
}

func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//todo: 同源验证,受信任的源
		return h(c)
	}
}

//检查会员编号
func getMemberId(c *echox.Context) int {
	v := c.Session.Get("member")
	if v != nil {
		m := v.(*member.Member)
		if m != nil {
			return m.Id
		}
	}
	return -1
}

func requestLogin(c *echox.Context) error {
	msg := gof.Message{
		Message: "not login",
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), msg)
}
