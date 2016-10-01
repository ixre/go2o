/**
 * Copyright 2015 @ z3q.net.
 * name : shared.go
 * author : jarryliu
 * date : 2016-07-12 09:32
 * description :
 * history :
 */
package shared

import (
	"github.com/jsix/gof/web"
	"go2o/x/echox"
	"gopkg.in/labstack/echo.v1"
	"html/template"
)

const (
	AppPlatform = 1 << iota
	AppMerchant
	AppPassport
	AppUCenter
	AppShop
)

var (
	// 模板监视更改
	TemplateObserverFlag int = AppPlatform | AppMerchant | AppPassport | AppUCenter | AppShop
)

// 提示页面
func RenderMessagePage(c *echox.Context, msg string, btn string, url string) error {
	d := c.NewData()
	d.Map = map[string]interface{}{
		"Message":    template.HTML(msg),
		"ButtonText": btn,
		"HasButton":  btn != "",
		"Url":        url,
	}
	return c.RenderOK("message_page.html", d)
}

func HandleHttpError(err error, ctx *echo.Context) {
	web.HttpError(ctx.Response(), err)
}
