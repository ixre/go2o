/**
 * Copyright 2015 @ z3q.net.
 * name : list_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

type ContentC struct {
}

// 自定义页面
func (this *ContentC) Page(ctx *echox.Context) error {
	p := getPartner(ctx)
	mm := GetMember(ctx)
	siteConf := getSiteConf(ctx)
	form := ctx.Request().URL.Query()

	var page *content.ValuePage
	idParam := form.Get("id")
	if b, _ := regexp.MatchString("^\\d+$", idParam); b {
		id, _ := strconv.Atoi(form.Get("id"))
		page = dps.ContentService.GetPage(p.Id, id)
	} else {
		page = dps.ContentService.GetPageByIndent(p.Id, idParam)
	}

	if page == nil {
		http.Error(ctx.Response(), "404 page not found.", http.StatusNotFound)
		return nil
	}

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"member":  mm,
		"conf":    siteConf,
		"page":    page,
		"rawBody": template.HTML(page.Body),
	}
	return ctx.RenderOK("page.html", d)
}
