/**
 * Copyright 2015 @ S1N1 Team.
 * name : static
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"github.com/atnet/gof/web"
	"net/http"
)

// 处理静态文件
var HttpStaticFileHandler = func(ctx *web.Context) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "."+ctx.Request.URL.Path)
}
