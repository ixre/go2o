/**
 * Copyright 2015 @ z3q.net.
 * name : main.c
 * author : jarryliu
 * date : 2016-06-06 00:18
 * description :
 * history :
 */
package passport

import (
    "github.com/jsix/gof"
    "go2o/x/echox"
    "net/http"
)

type mainC struct{
    gof.App
}

func (this *mainC) Index(ctx *echox.Context)error{
    return ctx.Redirect(http.StatusTemporaryRedirect,"sso/login")
}
