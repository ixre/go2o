/**
 * Copyright 2015 @ z3q.net.
 * name : static
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"github.com/jsix/gof/web"
	"net/http"
	"strings"
)

// 处理静态文件
var HttpStaticFileHandler = func(ctx *web.Context) {
	http.ServeFile(ctx.Response, ctx.Request, "."+ctx.Request.URL.Path)
}

var HttpImageFileHandler = func(ctx *web.Context) {
	path := strings.Replace(ctx.Request.URL.Path,"/img/","",1)
	http.ServeFile(ctx.Response, ctx.Request, "./static/uploads/"+path)
}