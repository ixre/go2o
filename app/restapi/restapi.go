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
	"github.com/jsix/gof/storage"
	"go2o/app/cache"
	"go2o/app/util"
	"go2o/core/domain/interface/merchant"
	"gopkg.in/labstack/echo.v1"
	"net/http"
	"strconv"
)

// 获取存储
func GetStorage() storage.Interface {
	return sto
}

// 获取传入的商户接口编号和密钥
func getUserInfo(ctx *echo.Context) (string, string) {
	r := ctx.Request()
	apiId := r.FormValue("merchant_id")
	apiSecret := r.FormValue("secret")
	if len(apiId) == 0 {
		apiId = r.URL.Query().Get("merchant_id")
	}

	//todo: 兼容partner_id  ,将删除
	if len(apiId) == 0 {
		apiId = r.FormValue("partner_id")
		if len(apiId) == 0 {
			apiId = r.URL.Query().Get("partner_id")
		}
	}

	if len(apiSecret) == 0 {
		apiSecret = r.URL.Query().Get("secret")
	}
	return apiId, apiSecret
}

// 检查是否有权限调用接口(商户)
func chkMerchantApiSecret(ctx *echo.Context) bool {
	i, s := getUserInfo(ctx)
	ok, merchantId := CheckApiPermission(i, s)
	if ok {
		ctx.Set("merchant_id", merchantId)
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
func getMerchantId(ctx *echo.Context) int {
	return ctx.Get("merchant_id").(int)
}

// 获取会员编号
func GetMemberId(ctx *echo.Context) int {
	return ctx.Get("member_id").(int)
}

func ApiTest(ctx *echo.Context) error {
	return ctx.String(http.StatusOK, "It's working!")
}

// 检查是否有权限
func CheckApiPermission(apiId string, secret string) (ok bool, merchantId int) {
	if len(apiId) != 0 && len(secret) != 0 {
		var merchantId int = cache.GetMerchantIdByApiId(apiId)
		var apiInfo *merchant.ApiInfo = cache.GetMerchantApiInfo(merchantId)
		if apiInfo != nil {
			return apiInfo.ApiSecret == secret, merchantId
		}
	}
	return false, merchantId
}
