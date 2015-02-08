/**
 * Copyright 2014 @ Ops.
 * name :
 * author : newmin
 * date : 2013-10-11 21:04
 * description :
 * history :
 */

package www

import (
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/goclient"
	"net/http"
	"strings"
)

func GetSiteConf(w http.ResponseWriter, p *partner.ValuePartner) (bool, *partner.SiteConf) {
	siteConf, _ := goclient.Partner.GetSiteConf(p.Id, p.Secret)
	if siteConf.State == enum.PARTNER_SITE_CLOSED {
		if strings.TrimSpace(siteConf.StateHtml) == "" {
			siteConf.StateHtml = "网站暂停访问，请联系商家：" + p.Tel
		}
		w.Write([]byte(siteConf.StateHtml))
		return false, siteConf
	}
	return true, siteConf
}
