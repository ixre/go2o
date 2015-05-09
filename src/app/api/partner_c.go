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
    "go2o/src/core/service/dps"
)

type partnerC struct{
}

func (this *partnerC) Index(ctx *web.Context){
    dps.PartnerService.GetPartner(666888)
    ctx.ResponseWriter.Write([]byte("it's working!"))
}