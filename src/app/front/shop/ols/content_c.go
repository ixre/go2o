/**
 * Copyright 2015 @ S1N1 Team.
 * name : list_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/service/dps"
	"strconv"
	"regexp"
	"go2o/src/core/domain/interface/content"
	"net/http"
	"html/template"
)

type ContentC struct {
	*BaseC
}

// 自定义页面
func (this *ContentC) Page(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	mm := this.BaseC.GetMember(ctx)
	siteConf := this.BaseC.GetSiteConf(ctx)
	form := ctx.Request.URL.Query()

	var page *content.ValuePage
	idParam := form.Get("id")
	if b, _ := regexp.MatchString("^\\d+$", idParam); b {
		id, _ := strconv.Atoi(form.Get("id"))
		page = dps.ContentService.GetPage(p.Id, id)
	}else {
		page = dps.ContentService.GetPageByIndent(p.Id, idParam)
	}

	if page == nil {
		http.Error(ctx.ResponseWriter, "404 page not found.", http.StatusNotFound)
		return
	}

	this.BaseC.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":   p,
			"member":    mm,
			"conf":      siteConf,
			"page":page,
			"rawBody":template.HTML(page.Body),
		},
		"views/shop/ols/{device}/page.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}