/**
 * Copyright 2013 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-03 23:18
 * description :
 * history :
 */
package ols

import (
	"bytes"
	"fmt"
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/web"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
)

// 处理自定义错误
func HandleCustomError(w http.ResponseWriter, ctx *web.Context, err error) {
	if err != nil {
		w.WriteHeader(500)
		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"error": err.Error(),
			"stack": template.HTML(debug.Stack()),
		},
			strings.Replace("views/shop/ols/{device}/error.html", "{device}", ctx.Items["device_view_dir"].(string), -1))
	}
}

func GetShops(c gof.App, partnerId int) []byte {
	//分店
	var buf *bytes.Buffer = bytes.NewBufferString("")

	shops := dps.PartnerService.GetShopsOfPartner(partnerId)
	if len(shops) == 0 {
		return []byte("<div class=\"nodata noshop\">还未添加分店</div>")
	}
	buf.WriteString("<ul class=\"shops\">")
	for i, v := range shops {
		buf.WriteString(fmt.Sprintf(`<li class="s%d">
			<div class="name"><span><strong>%s</strong></div>
			<span class="shop-state shopstate%d">%s</span>
			<div class="phone">%s</div>
			<div class="address">%s</div>
			</li>`, i+1, v.Name, v.State, enum.GetFrontShopStateName(v.State), v.Phone, v.Address))
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}

func GetCategories(c gof.App, partnerId int, secret string) []byte {
	var buf *bytes.Buffer = bytes.NewBufferString("")
	categories := dps.SaleService.GetCategories(partnerId)

	buf.WriteString(`<ul class="categories">
		<li class="s0 current" val="0">
			<div class="name"><span><strong>全部</strong></div>
		</li>
	`)
	for i, v := range categories {
		buf.WriteString(fmt.Sprintf(`<li class="s%d" val="%d">
			<div class="name"><span><strong>%s</strong></div>
			</li>`, i+1, v.Id, v.Name))
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}
