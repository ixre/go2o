/**
 * Copyright 2015 @ z3q.net.
 * name : main.go
 * author : jarryliu
 * date : 2016-09-09 17:41
 * description :
 * history :
 */
package hapi

import (
    "github.com/jsix/gof"
    "gopkg.in/labstack/echo.v1"
    "net/http"
    "go2o/x/echox"
)

type mainC struct {
    gof.App
}

func (m *mainC) Info(ctx *echo.Context) error {
    return ctx.String(http.StatusOK, `
        release : 2016-09-10
    `)
}

// 测试HAPI
func (m *mainC) Test(ctx *echox.Context) error {
    memberId := getMemberId(ctx)
    if memberId <= 0 {
        return requestLogin(ctx)
    }
    d := gof.Message{
        Result:true,
        Data:memberId,
    }
    return ctx.JSONP(http.StatusOK, ctx.Query("callback"), d)
}
