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
	"github.com/labstack/echo"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"net/http"
	"strings"
)

type partnerC struct {
}

// 获取广告数据
func (this *partnerC) Get_ad(ctx *echo.Context) error {
	partnerId := getPartnerId(ctx)
	adName := ctx.Request().FormValue("ad_name")
	adv, data := dps.AdvertisementService.GetAdvertisementAndDataByName(partnerId, adName)
	if data != nil {
		// 图片广告
		if adv.Type == ad.TypeGallery {
			gv := data.(ad.ValueGallery)
			if gv != nil {
				for _, v := range gv {
					if strings.Index(v.ImageUrl, "http://") == -1 {
						v.ImageUrl = format.GetGoodsImageUrl(v.ImageUrl)
					}
				}
			}
			return ctx.JSON(http.StatusOK, gv)
		}
		return ctx.JSON(http.StatusOK, data)
	}
	return ctx.JSON(http.StatusOK, gof.Message{Message: "没有广告数据"})
}
