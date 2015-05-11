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
)

type partnerC struct {
}

func (this *partnerC) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("it's working!"))
}
