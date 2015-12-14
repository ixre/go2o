/**
 * Copyright 2015 @ z3q.net.
 * name : apI_portal
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"github.com/jsix/gof/web"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/x/echox"
	"net/http"
)

func ApiTest(ctx *echox.Context) error {
	return ctx.String(http.StatusOK, "It's working!")
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
