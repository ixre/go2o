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
<<<<<<< HEAD
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
=======
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

type ContentC struct {
<<<<<<< HEAD
}

// 自定义页面
func (this *ContentC) Page(c *echox.Context) error {
	p := getPartner(c)
	mm := GetMember(c)
	siteConf := getSiteConf(c)
	var page *content.ValuePage
	idParam := c.P(0)
	if b, _ := regexp.MatchString("^\\d+$", idParam); b {
		id, _ := strconv.Atoi(idParam)
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		page = dps.ContentService.GetPage(p.Id, id)
	} else {
		page = dps.ContentService.GetPageByIndent(p.Id, idParam)
	}

	if page == nil {
<<<<<<< HEAD
		return c.HTML(http.StatusNotFound, "Not found")
	}

	d := c.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"member":  mm,
		"conf":    siteConf,
		"page":    page,
		"rawBody": template.HTML(page.Body),
	}
	return c.RenderOK("page.html", d)
=======
		http.Error(ctx.Response, "404 page not found.", http.StatusNotFound)
		return
	}

	this.BaseC.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner": p,
			"member":  mm,
			"conf":    siteConf,
			"page":    page,
			"rawBody": template.HTML(page.Body),
		},
		"views/shop/ols/{device}/page.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
