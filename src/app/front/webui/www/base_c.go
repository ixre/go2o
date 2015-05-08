/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package www

import (
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
)

type baseC struct {
	// 是否检查用户登陆状态
	memberChk bool
}

func NewFilterImpl(memberChk bool)*baseC{
	return &baseC{
		memberChk : memberChk,
	}
}

func (this *baseC) Requesting(*web.Context) bool{
	//if(this.GetPartnerId())
	return true
}

func (this *baseC) RequestEnd(*web.Context){
}

// 获取商户编号
func (this *baseC) GetPartnerId(ctx *web.Context) int {
	currHost := ctx.Request.Host
	host := ctx.Session().Get("webui_host")
	pid := ctx.Session().Get("webui_pid")
	if host == nil || pid == nil || host != currHost {
		partnerId := dps.PartnerService.GetPartnerIdByHost(currHost)
		if partnerId != -1 {
			ctx.Session().Set("webui_host", currHost)
			ctx.Session().Set("webui_pid", partnerId)
			ctx.Session().Save()
		}
		return partnerId
	}
	return pid.(int)
}

func (this *baseC) GetPartner(ctx *web.Context) (*partner.ValuePartner, error) {
	return dps.PartnerService.GetPartner(this.GetPartnerId(ctx))
}

// 获取会员
func (this *baseC) GetMember(ctx *web.Context) *member.ValueMember {
	memberIdObj := ctx.Session().Get("member")
	if memberIdObj != nil {
		if o, ok := memberIdObj.(*member.ValueMember); ok {
			return o
		}
	}
	return nil
}
