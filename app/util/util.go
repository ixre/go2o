/**
 * Copyright 2015 @ z3q.net.
 * name : util
 * author : jarryliu
 * date : 2016-06-06 17:53
 * description :
 * history :
 */
package util

import (
	"github.com/ixre/gof/web"
	"go2o/core/infrastructure/gen"
	"net/http"
	"net/url"
)

// 生成推广二维码,query为附加的参数查询
func GenerateInvitationQr(domain string, code string, query string) []byte {
	url := domain + "/i/" + code + "?device=3&" + query
	return gen.BuildQrCodeForUrl(url, 10)
}

//获取当前地址
func GetRawUrl(r *http.Request) string {
	return url.QueryEscape(web.RequestRawURI(r))
}
