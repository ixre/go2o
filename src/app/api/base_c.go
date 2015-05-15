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
	ok, partnerId := CheckApiPermission(apiId, apiSecret)
	if ok {
		ctx.Items["partner_id"] = partnerId
	}
	return ok
}

var _ mvc.Filter = new(BaseC)

type BaseC struct {
}

func (this *BaseC) Requesting(ctx *web.Context) bool {
	ctx.Request.ParseForm()
	if !chkApiSecret(ctx) {
		this.errorOutput(ctx, "secret incorrent!")
		return false
	}
	return true
}

func (this *BaseC) RequestEnd(ctx *web.Context) {

}

// 输出Json
func (this *BaseC) jsonOutput(ctx *web.Context, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		this.errorOutput(ctx, err.Error())
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

// 输出错误信息
func (this *BaseC) errorOutput(ctx *web.Context, err string) {
	ctx.ResponseWriter.Write([]byte("{error:\"" + err + "\"}"))
}

// 获取商户编号
func (this *BaseC) GetPartnerId(ctx *web.Context) int {
	return ctx.Items["partner_id"].(int)
}
