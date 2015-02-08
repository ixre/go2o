/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-03 23:18
 * description :
 * history :
 */
package www

import (
	"net/http"
	"net/url"
)

// 跳转到登录页面
func RedirectLoginPage(w http.ResponseWriter, returnUrl string) {
	var header http.Header = w.Header()
	header.Add("Location", "/login?return_url="+url.QueryEscape(returnUrl))
	w.WriteHeader(302)
}
