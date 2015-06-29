/**
 * Copyright 2015 @ S1N1 Team.
 * name : apI_portal
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"github.com/atnet/gof/web"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/partner"
)

func ApiTest(ctx *web.Context) {
	ctx.Response.Write([]byte("It's working!"))
}

// 检查是否有权限
func CheckApiPermission(apiId string, secret string) (ok bool, partnerId int) {
	if len(apiId) != 0 && len(secret) != 0 {
		var partnerId int = cache.GetPartnerIdByApiId(apiId)
		var apiInfo *partner.ApiInfo = cache.GetPartnerApiInfo(partnerId)
		if apiInfo != nil {
			return apiInfo.ApiSecret == secret, partnerId
		}
	}
	return false, partnerId
}
