/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package front

import (
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/core/domain/interface/member"
)

type WebC struct {
}

// 获取商户编号
func (this *WebC) GetPartnerId(ctx *web.Context) int {
	currHost := ctx.Request.Host
	host := ctx.Session().Get("webui_host")
	pid := ctx.Session().Get("webui_pid")
	if host == nil || pid == nil || host != currHost {
		partnerId := dps.PartnerService.GetPartnerIdByHost(currHost)
		if partnerId != -1 {
			ctx.Session().Set("webui_host", currHost)
			ctx.Session().Set("webui_pid",partnerId)
			ctx.Session().Save()
		}
		return partnerId
	}
	return pid.(int)
}

func (this *WebC) GetPartner(ctx *web.Context) (*partner.ValuePartner, error) {
	return dps.PartnerService.GetPartner(this.GetPartnerId(ctx))
}

//获取会员
func (this *WebC) GetMember(ctx *web.Context) *member.ValueMember{
	memberIdObj := ctx.Session().Get("member")
	if memberIdObj != nil{
		if o,ok := memberIdObj.(*member.ValueMember);ok{
			return o
		}
	}
	return nil
}