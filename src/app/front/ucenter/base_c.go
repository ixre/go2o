/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"net/url"
	"strings"
)

var _ mvc.Filter = new(baseC)

type baseC struct {
}

func (this *baseC) Requesting(ctx *web.Context) bool {
	//验证是否登陆
	s := ctx.Session().Get("member")
	if s != nil {
		if m := s.(*member.ValueMember); m != nil {
			ctx.Items["member"] = m
			return true
		}
	}
	ctx.ResponseWriter.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(ctx.Request.URL.String()) + "'</script>"))
	return false
}

func (this *baseC) RequestEnd(ctx *web.Context) {
}

// 获取商户
func (this *baseC) GetPartner(ctx *web.Context) *partner.ValuePartner {
	val := ctx.Session().Get("member:rel_partner")
	if val != nil {
		return cache.GetValuePartnerCache(val.(int))
	}
	return nil
}

// 获取会员
func (this *baseC) GetMember(ctx *web.Context) *member.ValueMember {
	return ctx.Items["member"].(*member.ValueMember)
}

// 获取商户的站点设置
func (this *baseC) GetSiteConf(partnerId int) *partner.SiteConf {
	return dps.PartnerService.GetSiteConf(partnerId)
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

// 输出错误信息
func (this *baseC) errorOutput(ctx *web.Context, err string) {
	ctx.ResponseWriter.Write([]byte("{error:\"" + err + "\"}"))
}

// 输出错误信息
func (this *baseC) resultOutput(ctx *web.Context, result gof.Message) {
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf(
		"{result:%v,code:%d,message:\"%s\"}", result.Result, result.Code, result.Message)))
}


// 执行模板
func (this *baseC) ExecuteTemplate(ctx *web.Context, d gof.TemplateDataMap, files ...string) {
	newFiles := make([]string, len(files))
	for i, v := range files {
		newFiles[i] = strings.Replace(v, "{device}", ctx.Items["device_view_dir"].(string), -1)
	}
	ctx.App.Template().Execute(ctx.ResponseWriter, d, newFiles...)
}
