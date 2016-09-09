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
    "gopkg.in/labstack/echo.v1"
    "github.com/jsix/gof"
)

func NewServe() *echo.Echo {
    s := echo.New()
    //s.Use(mw.Recover())
   // s.Use(echox.StopAttack)
    s.Use(beforeHanding)
    registerRoutes(s)
    return s
}

func registerRoutes(s *echo.Echo){
    app := gof.CurrentApp
    mc := &mainC{app}
    s.Get("/api_info",mc.Info)
}

func beforeHanding(h echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx *echo.Context) error {
        return h(ctx)
    }
}