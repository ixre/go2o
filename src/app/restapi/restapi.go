/**
 * Copyright 2015 @ z3q.net.
 * name : base_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"github.com/jsix/gof"
	"go2o/src/app/cache"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/partner"
	"gopkg.in/labstack/echo.v1"
	"net/http"
	"strconv"
)

// 获取存储
func GetStorage() gof.Storage {
	return sto
}

// 获取传入的商户接口编号和密钥
func getUserInfo(ctx *echo.Context) (string, string) {
	r := ctx.Request()
	apiId := r.FormValue("partner_id")
	apiSecret := r.FormValue("secret")
	if len(apiId) == 0 {
		apiId = r.URL.Query().Get("partner_id")
	}
	if len(apiSecret) == 0 {
		apiSecret = r.URL.Query().Get("secret")
	}
	return apiId, apiSecret
}

// 检查是否有权限调用接口(商户)
func chkPartnerApiSecret(ctx *echo.Context) bool {
	i, s := getUserInfo(ctx)
	ok, partnerId := CheckApiPermission(i, s)
	if ok {
		ctx.Set("partner_id", partnerId)
	}
	return ok
}

// 检查会员令牌信息
func checkMemberToken(ctx *echo.Context) bool {
	r := ctx.Request()
	sto := gof.CurrentApp.Storage()
	memberId, _ := strconv.Atoi(r.FormValue("member_id"))
	token := r.FormValue("member_token")

	if util.CompareMemberApiToken(sto, memberId, token) {
		ctx.Set("member_id", memberId)
		return true
	}
	return false
}

// 获取商户编号
func getPartnerId(ctx *echo.Context) int {
	return ctx.Get("partner_id").(int)
}

// 获取会员编号
func GetMemberId(ctx *echo.Context) int {
	return ctx.Get("member_id").(int)
}

func ApiTest(ctx *echo.Context) error {
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
