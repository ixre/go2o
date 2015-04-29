package webui

import (
	"fmt"
	"net/http"
)

//处理Webui请求
func Handle(ctx *web.Context) {

	fmt.Fprintf(w, "hello,"+r.Host+"/")
}
