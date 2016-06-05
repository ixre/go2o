/**
 * Copyright 2015 @ z3q.net.
 * name : sso_c.go
 * author : jarryliu
 * date : 2016-06-06 00:19
 * description :
 * history :
 */
package passport

import (
    "github.com/jsix/gof"
    "go2o/x/echox"
)

type ssoC struct{
    gof.App
}

// 单点登陆
func (this *ssoC) Login(ctx *echox.Context)error {
    r := ctx.HttpRequest()
    var tipStyle string
    var returnUrl string = r.URL.Query().Get("return_url")
    if len(returnUrl) == 0 {
        tipStyle = " hidden"
    }
    d := ctx.NewData()
    d.Map["TipStyle"] = tipStyle
    return ctx.RenderOK("sso_login.html", d)
}
