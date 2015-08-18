/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package api

import (
	"encoding/json"
	"fmt"
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/net/jsv"
	"github.com/jrsix/gof/web"
	"go2o/src/core/service/goclient"
)

type websocketC struct {
	gof.App
}

func (this *websocketC) Login(ctx *web.Context) {
	ctx.Response.Write([]byte("ok"))
}

func (this *websocketC) Test(ctx *web.Context) {
	w := ctx.Response
	result, _ := goclient.Member.Login("jarryliu", "123000")

	if result.Result {
		w.Write([]byte("[Login]:Sucessfull." + result.Member.DynamicToken))
	} else {
		w.Write([]byte("[Login]:Failed." + result.Message))
	}
}

func (this *websocketC) Partner(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	buffer := goclient.Redirect.Post([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetPartner`,
		r.FormValue("partner_id"), r.FormValue("secret"))), 512)
	w.Write(buffer)
}

func (this *websocketC) Category(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	buffer := goclient.Redirect.Post([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.Category`,
		r.FormValue("partner_id"), r.FormValue("secret"))), 2048)

	var v jsv.Result
	jsv.JsonCodec.Unmarshal(buffer, &v)
	b, _ := json.Marshal(v.Data)
	w.Write(b)
}
