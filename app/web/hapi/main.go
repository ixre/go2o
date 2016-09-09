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
)

type mainC struct{
    gof.App
}

func (m *mainC) Info(ctx *echo.Context)error{
    return ctx.String(http.StatusOK,`
        release : 2016-09-10
    `)
}
