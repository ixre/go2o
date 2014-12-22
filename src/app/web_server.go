/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-16 21:47
 * description :
 * history :
 */

package app

import (
	"com/infrastructure"
	"com/ording/apiserv"
	"com/ording/partner"
	"com/ording/ucenter"
	"com/ording/webui/mobi"
	"com/ording/webui/weixin"
	"com/ording/webui/www"
	"com/service/goclient"
	"com/share/glob"
	"com/share/variable"
	"flag"
	"fmt"
	"net/http"
	"ops/cf/app"
	"ops/cf/web"
	"os"
	"strconv"
	"strings"
)

var (
	API_DOMAIN string
)

func RunWeb(ctx app.Context, port int, debug bool) {
	var (
		proxy *web.HttpHandleProxy
	)

	flag.Parse()

	if gcx, ok := ctx.(*glob.AppContext); ok {
		if !gcx.Loaded {
			gcx.Init(debug)
		}
	} else {
		fmt.Println("app context err")
		os.Exit(1)
		return
	}

	if debug {
		fmt.Println("[Started]:Web server (with debug) running on port [" +
			strconv.Itoa(port) + "]:")
		infrastructure.DebugMode = true
	} else {
		fmt.Println("[Started]:Web server running on port [" + strconv.Itoa(port) + "]:")
	}

	//socket client
	API_DOMAIN = ctx.Config().Get(variable.ApiDomain)
	goclient.Configure("tcp", ctx.Config().Get(variable.ClientSocketServer), ctx)

	proxy = getHttpProxy()

	//注册路由
	registRoutes(ctx)
	http.HandleFunc("/", proxy.For(nil))

	//启动服务
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if err != nil {
		ctx.Log().Fatalln("ListenAndServer ", err)
	}
}

//获取Http请求代理处理程序
func getHttpProxy() *web.HttpHandleProxy {
	return &web.HttpHandleProxy{
		Before: func(w http.ResponseWriter, r *http.Request) bool {

			//api端
			if r.Host == API_DOMAIN {
				apiserv.HandleRequest(w, r)
				return false
			}

			switch host, f := r.Host, strings.HasPrefix; {

			//供应商端
			case f(r.Host, "partner."):
				partner.HandleRequest(w, r)

				//会员端
			case f(r.Host, variable.DOMAIN_MEMBER_PREFIX):
				ucenter.HandleRequest(w, r)

				//静态资源
			case f(host, "static."):
				file := r.URL.Path
				http.ServeFile(w, r, "./static"+file)
				return false

				//静态资源
			case f(host, "img."):
				file := r.URL.Path
				http.ServeFile(w, r, "./static/uploads/"+file)
				return false

			default:
				//mobi.HandleRequest(weixin, r)
				www.HandleRequest(w, r)
			}

			return true
		},
		After:  nil,
		Except: web.DefaultHttpExceptHandle,
	}
}

func registRoutes(context app.Context) {
	partner.RegistRoutes(context)
	ucenter.RegistRoutes(context)
	www.RegistRoutes(context)
	mobi.RegistRoutes(context)
	weixin.RegistRoutes(context)
	apiserv.RegistRoutes(context)
}
