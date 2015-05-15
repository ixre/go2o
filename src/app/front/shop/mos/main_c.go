/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mobi

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
)

type mainC struct {
	gof.App
}

func (this *mainC) Login(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	ctx.App.Template().Execute(w, "views/ucenter/login.html", nil)
}

func (this *mainC) Index(ctx *web.Context, p *partner.ValuePartner) {
	_, w := ctx.Request, ctx.ResponseWriter
	w.Write([]byte(p.Name))
}
