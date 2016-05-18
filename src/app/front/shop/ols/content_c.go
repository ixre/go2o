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
func (this *ContentC) Page(c *echox.Context) error {
	p := getMerchant(c)
	mm := GetMember(c)
	siteConf := getSiteConf(c)
	var page *content.ValuePage
	idParam := c.P(0)
	if b, _ := regexp.MatchString("^\\d+$", idParam); b {
		id, _ := strconv.Atoi(idParam)
		page = dps.ContentService.GetPage(p.Id, id)
	} else {
		page = dps.ContentService.GetPageByIndent(p.Id, idParam)
	}

	if page == nil {
		return c.HTML(http.StatusNotFound, "Not found")
	}

	d := c.NewData()
	d.Map = gof.TemplateDataMap{
		"Merchant": p,
		"member":   mm,
		"conf":     siteConf,
		"page":     page,
		"rawBody":  template.HTML(page.Body),
	}
	return c.RenderOK("page.html", d)
}
