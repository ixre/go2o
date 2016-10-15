/**
 * Copyright 2015 @ z3q.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"github.com/jsix/gof"
	"go2o/core/service/dps"
	"gopkg.in/labstack/echo.v1"
	"net/http"
)

type merchantC struct {
}

// 获取广告数据
func (m *merchantC) Get_ad(ctx *echo.Context) error {
	merchantId := getMerchantId(ctx)
	adName := ctx.Request().FormValue("ad_name")
	dto := dps.AdService.GetAdAndDataByKey(merchantId, adName)
	if dto != nil {
		return ctx.JSON(http.StatusOK, dto)
	}
	return ctx.JSON(http.StatusOK,
		gof.Message{Message: "没有广告数据"})
}
