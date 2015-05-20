/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mos

import (
	"github.com/atnet/gof/web"
)

type mainC struct {
	*baseC
}

func (this *mainC) Login(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	ctx.App.Template().Execute(w, nil,"views/ucenter/login.html")
}

func (this *mainC) Index(ctx *web.Context) {
	if this.Requesting(ctx) {
		_, w := ctx.Request, ctx.ResponseWriter
		p := this.GetPartner(ctx)
		w.Write([]byte(p.Name))
	}
}
