/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-03 23:18
 * description :
 * history :
 */
package www

import (
	"github.com/atnet/gof"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/goclient"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
	"go2o/src/cache"
)



func GetSiteConf(w http.ResponseWriter, p *partner.ValuePartner) (bool, *partner.SiteConf) {
	var conf = cache.GetPartnerSiteConf(p.Id)
	if conf == nil {
		conf, _ = goclient.Partner.GetSiteConf(p.Id, p.Secret)

		if conf == nil {
			w.Write([]byte("网站访问过程中出现了异常，请重试!"))
			return false, nil
		}

		if conf.State == enum.PARTNER_SITE_CLOSED {
			if strings.TrimSpace(conf.StateHtml) == "" {
				conf.StateHtml = "网站暂停访问，请联系商家：" + p.Tel
			}
			w.Write([]byte(conf.StateHtml))
			return false, conf
		}
		cache.SetPartnerSiteConf(p.Id,conf)
	}
	return true, conf
}

// 处理自定义错误
func handleCustomError(w http.ResponseWriter, ctx gof.App, err error) {
	if err != nil {
		ctx.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["error"] = err.Error()
			(*m)["statck"] = template.HTML(strings.Replace(string(debug.Stack()), "\n", "<br />", -1))
		},
			"views/web/www/error.html")
	}
}
