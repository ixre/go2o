/**
 * Copyright 2015 @ z3q.net.
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
	"go2o/app/api"
	"go2o/core/variable"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	API_DOMAIN   string
	API_HOST_CHK bool = false // 必须匹配Host
	PathPrefix        = "/go2o_api_v1"
	sto          storage.Interface
	serve        *echo.Echo
)

func init() {
	serve = newServe()
}

func newServe() *echo.Echo {
	serve := echo.New()
	serve.Use(mw.Recover())
	serve.Use(beforeRequest())

	//todo:  echo
	//serve.Hook(splitPath) // 获取新的路径,在请求之前发生
	registerRoutes(serve)
	registerNewApi(serve)
	return serve
}

// 注册新的服务接口
func registerNewApi(s *echo.Echo) {
	mux := api.NewServe(false, "1.0.0")
	s.GET("/api", func(ctx echo.Context) error {
		return ctx.String(200, "go2o api server")
	})
	s.POST("/api", func(ctx echo.Context) error {
		mux.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})
}

// 获取服务实例
func GetServe() *echo.Echo {
	return serve
}

func Run(app gof.App, port int) {
	sto = app.Storage()
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	log.Println("** [ Go2o][ API][ Booted] - Api server running on port " +
		strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), serve)
}

func registerRoutes(s *echo.Echo) {
	pc := &merchantC{}
	mc := &MemberC{}
	gc := &getC{}
	s.GET("/", ApiTest)
	s.GET(PathPrefix+"/get/invite_qr", gc.Invite_qr) // 获取二维码
	s.GET(PathPrefix+"/get/gen_qr", gc.GenQr)        //生成二维码
	s.POST(PathPrefix+"/mm_login", mc.Login)         // 会员登录接口
	//s.POST(PathPrefix+"/mm_register", mc.Register)   // 会员注册接口
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
