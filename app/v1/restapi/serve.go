/**
 * Copyright 2015 @ to2.net.
 * name : rest_server.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	apiv2 "go2o/app/api"
	"go2o/app/v1/api"
	"go2o/core/variable"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	API_DOMAIN   string
	API_HOST_CHK = false // 必须匹配Host
	PathPrefix   = "/go2o_api_v1"
	store        storage.Interface
	serve        *echo.Echo
)

func newServe(config *gof.Config, store storage.Interface) *echo.Echo {
	serve := echo.New()
	serve.Use(mw.Recover())
	serve.Use(beforeRequest())

	//todo:  echo
	//serve.Hook(splitPath) // 获取新的路径,在请求之前发生
	registerRoutes(serve)
	registerNewApi(serve, config, store)
	registerApiV2(serve, store)
	return serve
}

func registerApiV2(s *echo.Echo, store storage.Interface) {
	prefix := "/a/v2"
	h := apiv2.ServeApiV2(store, prefix, false, "", "", "")
	hf := func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
	route := prefix + "/*"
	s.POST(route, hf)
	s.GET(route, hf)
	s.PUT(route, hf)
	s.PATCH(route, hf)
	s.DELETE(route, hf)
	s.OPTIONS(route, hf)
}

// 注册新的服务接口
func registerNewApi(s *echo.Echo, config *gof.Config, store storage.Interface) {
	reqVer := config.GetString("api_require_version")
	apiUser := config.GetString("api_user")
	apiSecret := config.GetString("api_secret")
	mux := api.NewServe(store, false, reqVer, apiUser, apiSecret)
	s.GET("/api", func(ctx echo.Context) error {
		return ctx.String(200, "go2o api server")
	})
	s.POST("/api", func(ctx echo.Context) error {
		mux.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})
	s.OPTIONS("/api", func(ctx echo.Context) error {
		mux.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})
}

func Run(ch chan bool, app gof.App, port int) {
	store = app.Storage()
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	handler := newServe(app.Config(), store)
	log.Println("[ Go2o][ API]: api gateway serve on port :" + strconv.Itoa(port))

	err := http.ListenAndServe(":"+strconv.Itoa(port), handler)

	//r := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	//s := web.NewService(
	//	//web.RegisterTTL(time.Second*2),
	//	web.Name("Go2oService"),
	//	web.Address(":"+strconv.Itoa(port)),
	//	web.Handler(handler),
	//	web.Registry(r))
	//if err := s.Run();err != nil {
	if err != nil {
		log.Println("** [ Go2o][ API] : " + err.Error())
		os.Exit(1)
	}
	ch <- false
}

func registerRoutes(s *echo.Echo) {
	pc := &merchantC{}
	mc := &MemberC{}
	gc := &getC{}
	s.GET("/", ApiTest)
	s.GET(PathPrefix+"/get/invite_qr", gc.Invite_qr) // 获取二维码
	s.GET(PathPrefix+"/get/gen_qr", gc.GenQr)        //生成二维码
	s.POST(PathPrefix+"/mm_login", mc.Login)         // 会员登录接口
	//s.POST(PathPrefix+"/mm_register", mc.Handle)   // 会员注册接口
	s.POST(PathPrefix+"/merchant/get_ad", pc.Get_ad) // 商户广告接口
	s.POST(PathPrefix+"/partner/get_ad", pc.Get_ad)  // 商户广告接口
	//s.Post("/member/*",mc)  // 会员接口
}

func beforeRequest() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := c.Request().URL.Host
			path := c.Request().URL.Path
			// note: 新接口
			if strings.HasPrefix(path, "/a/v2") {
				return h(c)
			}
			if path == "/api" {
				return h(c)
			}
			// todo: path compare
			if API_HOST_CHK && host != API_DOMAIN {
				return c.String(http.StatusNotFound, "no such file")
			}

			if path != "/" {
				//检查商户接口权限
				c.Request().ParseForm()
				if !chkMerchantApiSecret(c) {
					return c.String(http.StatusOK, "{error:\"incorrent secret\"}")
				}
				//检查会员会话
				if strings.HasPrefix(path, "/member") && !checkMemberToken(c) {
					return c.String(http.StatusOK, "{error:\"incorrent session\"}")
				}
			}

			return h(c)
		}
	}
}
