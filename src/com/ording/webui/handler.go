package webui

import (
	"fmt"
	"net/http"
)

//处理Webui请求
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello,"+r.Host+"/")
}
