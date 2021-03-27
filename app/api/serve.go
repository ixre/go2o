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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof/jwt-api"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/app/api/member"
	"go2o/core/domain/interface/registry"
	"go2o/core/service"
	"go2o/core/service/proto"
	"net/http"
	"time"
)

var (
	RequireVersion = "1.0.0"
	ApiUser        = "go2o"
	ApiSecret      = "123456"
)

// 接口2.0
func ServeApiV2(store storage.Interface, prefix string, debug bool, requireVer string,
	apiUser string, apiSecret string) http.Handler {
	RequireVersion = requireVer
	// 请求限制
	rl := util.NewRequestLimit(store, 100, 10, 600)
	// 创建服务
	s := api.NewServerMux(swapApiKeyFunc, prefix, true)
	// 注册中间键
	serviceMiddleware(s, "[ Go2o][ API][ Log]: ", debug, rl)
	// 注册处理器
	s.HandlePublic(AccessTokenApi{})
	s.Handle(&fdApi{})             // 基础
	s.Handle(member.ComplainApi{}) // 投诉
	return s
}

// 服务调试跟踪
func serviceMiddleware(s api.Server, prefix string, debug bool, rl *util.RequestLimit) {
	prefix = "[ Api][ Log]"
	// 验证IP请求限制
	s.Use(func(ctx api.Context) error {
		addr := ctx.Request().UserAddr
		if len(addr) != 0 && !rl.Acquire(addr, 1) || rl.IsLock(addr) {
			return errors.New("您的网络存在异常,系统拒绝访问")
			//return errors.New("access denied")
		}
		return nil
	})
	//// 校验版本
	//s.Use(func(ctx api.Context) Error {
	//	//prod := ctx.StoredValues().GetString("product"
	//	prodVer := ctx.Params().GetString("version")
	//	if api.CompareVersion(prodVer, RequireVersion) < 0 {
	//		return errors.New("您当前使用的APP版本较低, 请升级或安装最新版本")
	//		//return errors.New(fmt.Sprintf("%s,require version=%s",
	//		//	api.RCDeprecated.Message, tarVer))
	//	}
	//	return nil
	//})

	if debug {
		// 开启调试
		s.Trace()
		// 输出请求信息
		s.Use(func(ctx api.Context) error {
			apiName := ctx.Request().RequestApi
			log.Println(prefix, "user", ctx.UserKey(), " calling ", apiName)
			data, _ := json.Marshal(ctx.Request().Params)
			log.Println(prefix, "request data = [", data, "]")
			// 记录服务端请求时间
			ctx.Request().Params.Set("$rpc_begin_time", time.Now().UnixNano())
			return nil
		})
	}
	if debug {
		// 输出响应结果
		s.After(func(ctx api.Context) error {
			form := ctx.Request().Params
			rsp := form.Get("$api_response").(*api.Response)
			data := ""
			if rsp.Data != nil {
				d, _ := json.Marshal(rsp.Data)
				data = string(d)
			}
			reqTime := int64(form.GetInt("$rpc_begin_time"))
			elapsed := float32(time.Now().UnixNano()-reqTime) / 1000000000
			log.Println(prefix, "Response : ", rsp.Code, rsp.Message,
				fmt.Sprintf("; elapsed time ：%.4fs ; ", elapsed),
				"Result = [", data, "]",
			)
			//if rsp.Code == api.RAccessDenied.Code {
			//	data, _ := url.QueryUnescape(ctx.Request().Form.Encode())
			//	sortData := api.ParamsToBytes(ctx.Request().Form, form.GetString("$user_secret"), true)
			//	log.Println(prefix, "request data = [", data, "]")
			//	log.Println(" sign not match ! key =", form.Get("key"),
			//		"\r\n   server_sign=", form.GetString("$server_sign"),
			//		"\r\n   client_sign=", form.GetString("$client_sign"),
			//		"\r\n   sort_params=", string(sortData))
			//}
			return nil
		})
	}
}

// 交换接口用户凭据
func swapApiKeyFunc(ctx api.Context) (privateKey []byte, err error) {
	return getJWTSecret(), nil
}

var jwtSecret = []byte("")

func getJWTSecret() []byte {
	if len(jwtSecret) == 0 {
		trans, cli, err := service.RegistryServiceClient()
		if err == nil {
			defer trans.Close()
			value, _ := cli.GetValue(context.TODO(), &proto.String{Value: registry.SysJWTSecret})
			if len(value.Value) == 0 {
				log.Println("[ Go2o][ Warning]: jwt secret is empty")
			}
			jwtSecret = []byte(value.Value)
		}else{
			log.Println("[ Go2o][ Warning]: get jwt secret error: ",err.Error())
		}
	}
	return jwtSecret
}
