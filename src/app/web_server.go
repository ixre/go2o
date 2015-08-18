/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:47
 * description :
 * history :
 */

package app

import (
	"fmt"
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/web"
	"go2o/src/app/front/master"
	"go2o/src/app/front/partner"
	"go2o/src/app/front/shop/ols"
	"go2o/src/app/front/ucenter"
	"go2o/src/core/infrastructure"
	"go2o/src/core/variable"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//获取Http请求代理处理程序
func getInterceptor(a gof.App) *web.Interceptor {
	var igor = web.NewInterceptor(a, getHttpExecFunc())
	igor.Except = web.HandleDefaultHttpExcept
	igor.Before = func(ctx *web.Context) bool {

		r := ctx.Request
		//静态资源
		if strings.HasPrefix(r.Host, "static.") {
			http.ServeFile(ctx.Response, r, "./static"+r.URL.Path)
			return false
		} else if strings.HasPrefix(r.Host, "img.") {
			http.ServeFile(ctx.Response, r, "./static/uploads/"+r.URL.Path)
			return false
		}
		return true
	}
	igor.After = nil
	return igor
}

// 获取执行方法
func getHttpExecFunc() web.RequestHandler {
	return func(ctx *web.Context) {

		r, _ := ctx.Request, ctx.Response

		switch host, f := r.Host, strings.HasPrefix; {
		//case host == API_DOMAIN:
		//	apiserv.Handle(ctx)

		//供应商端
		case f(host, variable.DOMAIN_PREFIX_PARTNER):
			partner.Handle(ctx)

		//会员端
		case f(host, variable.DOMAIN_PREFIX_MEMBER):
			ucenter.Handle(ctx)

		//管理中心
		case f(host, variable.DOMAIN_PREFIX_WEBMASTER):
			master.Handle(ctx)

		//线上商店
		default:
			ols.Handle(ctx)
		}
	}
}

// 运行网页
func RunWeb(app gof.App, port int, debug, trace bool) {

	if debug {
		fmt.Println("[ Go2o][ Web][ Booted] - Web server (with debug) running on port " +
			strconv.Itoa(port))
		infrastructure.DebugMode = true
	} else {
		fmt.Println("[ Go2o][ Web][ Booted] - Web server running on port " + strconv.Itoa(port))
	}

	//socket client
	time.Sleep(time.Second * 2) //等待启动Socket
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)

	var in = getInterceptor(app)

	//启动服务
	err := http.ListenAndServe(":"+strconv.Itoa(port), in)

	if err != nil {
		app.Log().Fatalln("ListenAndServer ", err)
	}
}
