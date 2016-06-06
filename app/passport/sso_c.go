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
    "go2o/core/service/dps"
    "go2o/core/infrastructure/format"
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
    conf := dps.BaseService.GetPlatformConf()
    conf.Logo = format.GetResUrl(conf.Logo)
    d := ctx.NewData()
    d.Map["TipStyle"] = tipStyle
    d.Map["Conf"] = conf
    return ctx.RenderOK("sso_login.html", d)
}
