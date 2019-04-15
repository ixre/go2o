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
	"github.com/ixre/gof"
	"github.com/labstack/echo"
	"go2o/core/service/rsi"
	"net/http"
)

type merchantC struct {
}

// 获取广告数据
func (m *merchantC) Get_ad(c echo.Context) error {
	mchId := getMerchantId(c)
	adName := c.Request().FormValue("ad_name")
	dto := rsi.AdService.GetAdAndDataByKey(mchId, adName)
	if dto != nil {
		return c.JSON(http.StatusOK, dto)
	}
	return c.JSON(http.StatusOK,
		gof.Result{ErrCode: 1, ErrMsg: "没有广告数据"})
}
