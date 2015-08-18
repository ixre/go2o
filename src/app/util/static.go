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
	"github.com/jrsix/gof/web"
	"net/http"
)

// 处理静态文件
var HttpStaticFileHandler = func(ctx *web.Context) {
	http.ServeFile(ctx.Response, ctx.Request, "."+ctx.Request.URL.Path)
}
