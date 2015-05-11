/**
 * Copyright 2015 @ S1N1 Team.
 * name : base_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"encoding/json"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

// 检查是否有权限调用接口(商户)
func chkApiSecret(ctx *web.Context) bool {
	apiId := ctx.Request.FormValue("partner_id")
	apiSecret := ctx.Request.FormValue("secret")
	return CheckApiPermission(apiId,apiSecret)
}

var _ mvc.Filter = new(baseC)

type baseC struct {
}

func (this *baseC) Requesting(ctx *web.Context) bool {
	ctx.Request.ParseForm()

	if !chkApiSecret(ctx) {
		this.errorOutput(ctx, "secret incorrent!")
		return false
	}
	return true
}

func (this *baseC) RequestEnd(ctx *web.Context) {

}

// 输出Json
func (this *baseC) jsonOutput(ctx *web.Context, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		this.errorOutput(ctx, err.Error())
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

func (this *baseC) errorOutput(ctx *web.Context, err string) {
	ctx.ResponseWriter.Write([]byte("{error:\"" + err + "\"}"))
}
