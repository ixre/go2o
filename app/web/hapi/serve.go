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
	"go2o/app/web/shared"
	"go2o/core/domain/interface/member"
	"go2o/x/echox"
	"gopkg.in/labstack/echo.v1"
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

	s.Get("/api_info", mc.Info)
	s.Getx("/test", mc.Test)
	s.Getx("/request_login", mc.RequestLogin)
	s.Getx("/r/uc", mc.RedirectUc)
	s.Agetx("/user/:action", uc)    //用户
	s.Agetx("/service/:action", sc) //服务
}

func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		//todo: 同源验证,受信任的源
		return h(ctx)
	}
}

//检查会员编号
func getMemberId(ctx *echox.Context) int {
	v := ctx.Session.Get("member")
	if v != nil {
		m := v.(*member.Member)
		if m != nil {
			return m.Id
		}
	}
	return -1
}

func requestLogin(ctx *echox.Context) error {
	msg := gof.Message{
		Message: "not login",
	}
	return ctx.JSONP(http.StatusOK, ctx.Query("callback"), msg)
}
