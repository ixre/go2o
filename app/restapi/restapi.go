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
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"github.com/labstack/echo"
	"go2o/app/cache"
	"go2o/core/domain/interface/merchant"
	"go2o/core/service/thrift"
	"net/http"
	"strconv"
)

// 获取存储
func GetStorage() storage.Interface {
	return sto
}

// 获取传入的商户接口编号和密钥
func getUserInfo(c echo.Context) (string, string) {
	r := c.Request()
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
func chkMerchantApiSecret(c echo.Context) bool {
	i, s := getUserInfo(c)
	ok, mchId := CheckApiPermission(i, s)
	if ok {
		c.Set("merchant_id", mchId)
	}
	return ok
}

// 检查会员令牌信息
func checkMemberToken(c echo.Context) bool {
	r := c.Request()
	memberId, _ := util.I64Err(strconv.Atoi(r.FormValue("member_id")))
	token := r.FormValue("member_token")
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		if b, _ := cli.CheckToken(thrift.Context, memberId, token); b {
			c.Set("member_id", memberId)
			return true
		}
	}
	return false
}

// 获取商户编号
func getMerchantId(c echo.Context) int32 {
	return c.Get("merchant_id").(int32)
}

// 获取会员编号
func GetMemberId(c echo.Context) int64 {
	return c.Get("member_id").(int64)
}

func ApiTest(c echo.Context) error {
	return c.String(http.StatusOK, "It's working!")
}

// 检查是否有权限
func CheckApiPermission(apiId string, secret string) (bool, int32) {
	if len(apiId) != 0 && len(secret) != 0 {
		mchId := cache.GetMerchantIdByApiId(apiId)
		var apiInfo *merchant.ApiInfo = cache.GetMerchantApiInfo(mchId)
		if apiInfo != nil {
			return apiInfo.ApiSecret == secret, mchId
		}
	}
	return false, 0
}
