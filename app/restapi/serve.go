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
	"github.com/ixre/gof/storage"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"go2o/app/api"
	"go2o/core/variable"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	return serve
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

// 获取服务实例
func GetServe() *echo.Echo {
	return serve
}

func Run(app gof.App, port int) {
	store = app.Storage()
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	handler :=  newServe(app.Config(), store)
	log.Println("** [ Go2o][ API] - Api server running on port " +
		strconv.Itoa(port))

	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	s := web.NewService(
		web.Name("go2o"),
		web.Address(":"+strconv.Itoa(port)),
		web.Handler(handler),
		web.Registry(reg),
		web.RegisterTTL(time.Second*10),
		)
	//err := http.ListenAndServe(":"+strconv.Itoa(port),handler)
	if err := s.Run();err != nil {
		log.Println("** [ Go2o][ API] : " + err.Error())
		os.Exit(1)
	}
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
