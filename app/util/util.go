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
	"fmt"
	"net/http"
	"net/url"
)

//获取当前地址
func GetRawUrl(r *http.Request) string {
	query := r.URL.RawQuery
	proto := "http"
	if len(query) > 0 {
		query = "?" + query
	}
	if r.Proto == "HTTPS" {
		proto = "https"
	}
	return url.QueryEscape(fmt.Sprintf("%s://%s%s%s",
		proto, r.Host, r.URL.Path, query))
}
