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
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	API_DOMAIN string
)

func RunWeb(ctx app.Context, port int, debug, trace bool) {
	var (
		proxy *web.HttpHandleProxy
	)

	flag.Parse()

	if gcx, ok := ctx.(*glob.AppContext); ok {
		if !gcx.Loaded {
			gcx.Init(debug, trace)
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
	time.Sleep(time.Second * 2) //等待启动Socket
	API_DOMAIN = ctx.Config().GetString(variable.ApiDomain)
	goclient.Configure("tcp", ctx.Config().GetString(variable.ClientSocketServer), ctx)

	proxy = getHttpProxy()

	//注册路由
	RegisterRoutes(ctx)
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

func RegisterRoutes(context app.Context) {
	partner.RegisterRoutes(context)
	ucenter.RegisterRoutes(context)
	www.RegisterRoutes(context)
	mobi.RegisterRoutes(context)
	weixin.RegisterRoutes(context)
	apiserv.RegisterRoutes(context)
}
