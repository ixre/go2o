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
	"fmt"
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof/web"
	"github.com/labstack/echo"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/core/variable"
	"html/template"
	"regexp"
	"strings"
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
		"ErrMsg":     template.HTML(msg),
		"ButtonText": btn,
		"HasButton":  btn != "",
		"Url":        url,
	}
	return c.RenderOK("message_page.html", d)
}

// 处理HTTP错误
func HandleHttpError(err error, c echo.Context) {
	web.HttpError(c.Response(), err)
}

var (
	sysIgnoreRegex = regexp.MustCompile(".(gif|jpg|css|js|png|woff|ttf|woff2)$")
)

// 系统状态检测
func SystemCheck(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, _ := rsi.FoundationService.GetPlatformConf(thrift.Context)
		if conf.Suspend {
			rsp := c.Response()
			path := c.Request().URL.Path
			// 访问挂起页面及相关的资源页面不跳转
			if strings.Index(path, "suspend") == -1 &&
				!sysIgnoreRegex.MatchString(path) {
				url := fmt.Sprintf("http://%s%s/suspend", variable.DOMAIN_PREFIX_PORTAL,
					variable.Domain)
				rsp.Header().Add("Location", url)
				rsp.WriteHeader(302)
				return nil
			}
		}
		return h(c)
	}
}
