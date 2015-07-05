/**
 * Copyright 2015 @ S1N1 Team.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"strings"
)

type partnerC struct {
	*BaseC
}

func (this *partnerC) Index(ctx *web.Context) {
	ctx.Response.Write([]byte("it's working!"))
}


// 获取广告数据
func (this *partnerC) Get_ad(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	adName := ctx.Request.FormValue("ad_name")
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
			ctx.Response.JsonOutput(gv)
			return
		}
		ctx.Response.JsonOutput(data)
	} else {
		this.ErrorOutput(ctx, "没有广告数据")
	}
}
