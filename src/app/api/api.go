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
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/cache"
	"go2o/src/core/dto"
)

func HandleApi(ctx *web.Context) {
	//r := ctx.Request
	ctx.ResponseWriter.Write([]byte("It's working!"))
}

func CheckApiPermission(sto gof.Storage, apiId string, secret string) bool {
	if len(apiId) != 0 && len(secret) != 0 {
		var partnerId int = cache.GetPartnerIdByApiId(apiId)
		var apiInfo *dto.PartnerApiInfo = cache.GetPartnerApiInfo(partnerId)
		if apiInfo != nil {
			return apiInfo.ApiSecret == secret
		}
	}
	return false
}
