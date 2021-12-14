/**
 * Copyright 2015 @ 56x.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"context"
	"github.com/ixre/gof"
	"github.com/labstack/echo/v4"
	"go2o/core/service"
	"go2o/core/service/proto"
	"net/http"
)

type merchantC struct {
}

// 获取广告数据
func (m *merchantC) Get_ad(c echo.Context) error {
	mchId := getMerchantId(c)
	trans, cli, _ := service.AdvertisementServiceClient()
	defer trans.Close()
	adName := c.Request().FormValue("ad_name")
	dto, _ := cli.GetAdvertisement(context.TODO(),
		&proto.AdIdRequest{
			AdUserId:   mchId,
			AdKey:      adName,
			ReturnData: true,
		})
	if dto != nil {
		return c.JSON(http.StatusOK, dto.Data)
	}
	return c.JSON(http.StatusOK,
		gof.Result{ErrCode: 1, ErrMsg: "没有广告数据"})
}
