// HTTP API v1.0
// -----------------------
// 约定参数名称:
//	api       : 接口名称
//  api_key  : 接口用户
//  sign      : 签名
//  sign_type : 签名类型
//  app       : 应用编码
// -----------------------
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	RequireVersion = "1.0.0"
	ApiUser        = "go2o"
	ApiSecret      = "123456"
)

// 服务
func NewServe(store storage.Interface, debug bool, requireVer string,
	apiUser string, apiSecret string) http.Handler {
	RequireVersion = requireVer
	ApiUser = apiUser
	ApiSecret = apiSecret
	log.Println(fmt.Sprintf("[ Go2o][ API]: api key is '%s' and secret is '%s'", ApiUser, ApiSecret))
	// 初始化变量
	registry := map[string]interface{}{}
	// 创建上下文工厂
	factory := api.DefaultFactory.Build(registry)
	// 请求限制
	rl := util.NewRequestLimit(store, 200, 10, 600)
	serve := NewService(factory, debug, rl)
	// 创建http处理器
	hs := http.NewServeMux()
	hs.Handle("/api", serve)
	return hs
}

// 服务
func NewService(factory api.ContextFactory, debug bool, rl *util.RequestLimit) *api.ServeMux {
	// 创建服务
	s := api.NewServerMux(factory, swapApiKeyFunc, true)
	// 注册处理器
	s.Register("member", &MemberApi{})
	s.Register("article", &ArticleApi{})
	s.Register("app", NewAppApi())
	s.Register("register", NewRegisterApi())
	s.Register("passport", NewPassportApi())
	s.Register("settings", NewSettingsApi())
	s.Register("res", NewResApi())
	s.Register("goods", NewGoodsApi())
	s.Register("shop", NewShopApi())
	s.Register("account", NewAccountApi())
	// 注册中间键
	serviceMiddleware(s, "[ Go2o][ API][ Log]: ", debug, rl)
	return s
}

// 服务调试跟踪
func serviceMiddleware(s api.Server, prefix string, debug bool, rl *util.RequestLimit) {
	prefix = "[ Api][ Log]"
	if debug {
		// 开启调试
		s.Trace()
		// 输出请求信息
		s.Use(func(ctx api.Context) error {
			apiName := ctx.Form().Get("$api_name").(string)
			log.Println(prefix, "user", ctx.Key(), " calling ", apiName)
			data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
			log.Println(prefix, "request data = [", data, "]")
			// 记录服务端请求时间
			ctx.Form().Set("$rpc_begin_time", time.Now().UnixNano())
			return nil
		})
	}
	// API限流
	s.Use(func(ctx api.Context) error {
		// 验证IP请求限制
		addr := ctx.Form().GetString("$user_addr")
		if len(addr) != 0 && !rl.Acquire(addr, 1) || rl.IsLock(addr) {
			return errors.New("您的网络存在异常,系统拒绝访问")
		}
		return nil
	})
	// 校验版本
	s.Use(func(ctx api.Context) error {
		//prod := ctx.FormData().GetString("product"
		prodVer := ctx.Form().GetString("version")
		if api.CompareVersion(prodVer, RequireVersion) < 0 {
			return errors.New("您当前使用的APP版本较低, 请升级或安装最新版本")
			//return errors.New(fmt.Sprintf("%s,require version=%s",
			//	api.RDeprecated.Message, tarVer))
		}
		return nil
	})

	if debug {
		// 输出响应结果
		s.After(func(ctx api.Context) error {
			form := ctx.Form()
			rsp := form.Get("$api_response").(*api.Response)
			data := ""
			if rsp.Data != nil {
				d, _ := json.Marshal(rsp.Data)
				data = string(d)
			}
			reqTime := int64(ctx.Form().GetInt("$rpc_begin_time"))
			elapsed := float32(time.Now().UnixNano()-reqTime) / 1000000000
			log.Println(prefix, "response : ", rsp.Code, rsp.Message,
				fmt.Sprintf("; elapsed time ：%.4fs ; ", elapsed),
				"result = [", data, "]",
			)
			if rsp.Code == api.RAccessDenied.Code {
				data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
				sortData := api.ParamsToBytes(ctx.Request().Form, form.GetString("$user_secret"))
				log.Println(prefix, "request data = [", data, "]")
				log.Println(" sign not match ! key =", form.Get("key"),
					"\r\n   server_sign=", form.GetString("$server_sign"),
					"\r\n   client_sign=", form.GetString("$client_sign"),
					"\r\n   sort_params=", string(sortData))
			}
			return nil
		})
	}
}

// 交换接口用户凭据
func swapApiKeyFunc(ctx api.Context, key string) (userId int, userSecret string) {
	if key == ApiUser {
		return 1, ApiSecret
	} else if key == "go2o" {
		tt := time.Date(2019, 12, 01, 0, 0, 0, 0, time.Local)
		if time.Now().Unix() < tt.Unix() {
			return 1, "131409"
		}
	}
	//log.Println(fmt.Sprintf("[ UAMS][ API]: 接口用户[%s]交换凭据失败： %s", key, r.ErrMsg))
	return 0, ""
}
