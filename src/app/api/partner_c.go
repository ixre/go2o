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
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"strings"
)

type partnerC struct {
	*BaseC
}

func (this *partnerC) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("it's working!"))
}

// 处理请求
func (this *partnerC) handle(ctx *web.Context) {
	mvc.Handle(this, ctx, false)
}

// 获取广告数据
func (this *partnerC) GetAd(ctx *web.Context) {
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
			this.JsonOutput(ctx, gv)
			return
		}
		this.JsonOutput(ctx, data)
	} else {
		this.ErrorOutput(ctx, "没有广告数据")
	}
}
