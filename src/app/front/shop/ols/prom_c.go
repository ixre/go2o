/**
 * Copyright 2015 @ z3q.net.
 * name : prom_c.go
 * author : jarryliu
 * date : 2016-04-26 14:16
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof"
	"go2o/src/x/echox"
)

type promC struct {
}

func (this *promC) Coupon(ctx *echox.Context) error {
	p := getMerchant(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"Merchant": p,
		"Conf":     getSiteConf(ctx),
	}
	return ctx.RenderOK("coupon.html", d)
}
