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
	"go2o/src/app/util"
	"strconv"
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
		this.ErrorOutput(ctx, "secret incorrent!")
		return false
	}
	return true
}

func (this *BaseC) RequestEnd(ctx *web.Context) {

}

// 检查会员令牌信息
func (this *BaseC) CheckMemberToken(ctx *web.Context) bool {
	r := ctx.Request
	memberId, _ := strconv.Atoi(r.FormValue("member_id"))
	token := r.FormValue("member_token")

	if util.CompareMemberApiToken(ctx.App.Storage(), memberId, token) {
		ctx.Items["member_id"] = memberId
		return true
	}
	this.ErrorOutput(ctx, "invalid request!")
	return false
}

// 输出Json
func (this *BaseC) JsonOutput(ctx *web.Context, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		this.ErrorOutput(ctx, err.Error())
	} else {
		ctx.Response.Write(b)
	}
}

// 输出错误信息
func (this *BaseC) ErrorOutput(ctx *web.Context, err string) {
	ctx.Response.Write([]byte("{error:\"" + err + "\"}"))
}

// 获取商户编号
func (this *BaseC) GetPartnerId(ctx *web.Context) int {
	return ctx.Items["partner_id"].(int)
}

// 获取会员编号
func (this *BaseC) GetMemberId(ctx *web.Context) int {
	if v, ok := ctx.Items["member_id"].(int); ok {
		return v
	}
	return 0
}
