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
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	"go2o/app/web/shared"
	"go2o/core/service/thrift/idl/gen-go/define"
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
	s.GET("/test", mc.Test)
	s.GET("/request_login", mc.RequestLogin)
	s.GET("/r/uc", mc.RedirectUc)
	s.Auto("/user", uc)    //用户
	s.Auto("/service", sc) //服务
}

func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//todo: 同源验证,受信任的源
		return h(c)
	}
}

//检查会员编号
func getMemberId(c *echox.Context) int32 {
	v := c.Session.Get("member")
	if v != nil {
		m := v.(*define.Member)
		if m != nil {
			return m.ID
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
