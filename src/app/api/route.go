/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package api

import (
	"github.com/atnet/gof/web"
)

var (
	Routes web.Route = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	Routes.Handle(ctx)
}

func init() {
	bc := new(BaseC)
	pc := &partnerC{bc}
	mc := &MemberC{bc}
	Routes.Add("/", ApiTest)
	Routes.Add("/go2o_api_v1/mm_login", mc.Login)       // 会员登陆接口
	Routes.Add("/go2o_api_v1/mm_register", mc.Register) // 会员登陆接口
	Routes.Add("/go2o_api_v1/member/*", mc.Handle)      // 会员接口
	Routes.Add("/go2o_api_v1/partner/*", pc.handle)      // 会员接口
}
