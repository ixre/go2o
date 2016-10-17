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
	"fmt"
	"github.com/jsix/gof"
	"go2o/core/variable"
	"go2o/x/echox"
	"gopkg.in/labstack/echo.v1"
	"net/http"
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
		Result: true,
		Data:   memberId,
	}
	return ctx.JSONP(http.StatusOK, ctx.Query("callback"), d)
}

// 请求登录
func (m *mainC) RequestLogin(ctx *echox.Context) error {
	url := ctx.Request().Referer()
	return ctx.Redirect(http.StatusFound, fmt.Sprintf("%s://%s%s/auth?return_url=%s",
		variable.DOMAIN_PASSPORT_PROTO, variable.DOMAIN_PREFIX_PASSPORT,
		variable.Domain, url))
}

// 跳转到用户中心
func (m *mainC) RedirectUc(ctx *echox.Context) error {
	returnUrl := ctx.Query("url")
	if len(returnUrl) > 0 && returnUrl[0] != '/' {
		returnUrl = "/" + returnUrl
	}
	target := fmt.Sprintf("http://%s%s%s", variable.DOMAIN_PREFIX_MEMBER,
		variable.Domain, returnUrl)
	return ctx.Redirect(302, target)
}
