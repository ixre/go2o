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
	if len(query) > 0 {
		query = "?" + query
	}
	return url.QueryEscape(fmt.Sprintf("http://%s%s%s",
		r.Host, r.URL.Path, query))
}
