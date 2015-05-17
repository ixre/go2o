/**
 * Copyright 2013 @ S1N1 Team.
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
	"github.com/atnet/gof"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
)

// 处理自定义错误
func handleCustomError(w http.ResponseWriter, ctx gof.App, err error) {
	if err != nil {
		ctx.Template().Execute(w, gof.TemplateDataMap{
			"error":  err.Error(),
			"statck": template.HTML(strings.Replace(string(debug.Stack()), "\n", "<br />", -1)),
		},
			"views/shop/ols/error.html")
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
