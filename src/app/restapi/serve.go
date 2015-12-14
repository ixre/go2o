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
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go2o/src/core/variable"
	"log"
	"net/http"
	"strconv"
)

var (
	API_DOMAIN   string
	API_HOST_CHK bool = false // 必须匹配Host
	PathPrefix        = "/go2o_api_v1"
)

func RunRestApi(app gof.App, port int) {
	log.Println("** [ Go2o][ API][ Booted] - Api server running on port " + strconv.Itoa(port))
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	s := echo.New()
	s.Use(mw.Recover())
	s.Use(beforeRequest)
	s.Hook(splitPath) // 获取新的路径,在请求之前发生
	registerRoutes(s)
	s.Run(":" + strconv.Itoa(port)) //启动服务
}

func registerRoutes(s *echo.Echo) {
	pc := &partnerC{}
	mc := &MemberC{}
	gc := &getC{}

	s.Get("/", ApiTest)
	s.Get("/get/invite_qr", gc.Invite_qr) // 获取二维码
	s.Post("/mm_login", mc.Login)         // 会员登陆接口
	s.Post("/mm_register", mc.Register)   // 会员注册接口
	s.Post("/partner/get_ad", pc.Get_ad)  // 商户广告接口
	//s.Post("/member/*",mc)  // 会员接口
}

func beforeRequest(ctx *echo.Context) error {
	host := ctx.Request.URL.Host
	// todo: path compare
	if API_HOST_CHK && host != API_DOMAIN {
		http.Error(ctx.Response, "no such file", http.StatusNotFound)
		ctx.Done()
	}
	return nil
}

func splitPath(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = r.URL.Path[len(PathPrefix):]
}
