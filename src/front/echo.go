/**
 * Copyright 2015 @ z3q.net.
 * name : echo
 * author : jarryliu
 * date : 2015-12-04 10:51
 * description :
 * history :
 */
package front
import (
	"net/http"
	"strings"
)

type HttpHosts map[string]http.Handler

func (this HttpHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subName := r.Host[:strings.Index(r.Host,".")]
	 if h,ok := this[subName];ok{
		 h.ServeHTTP(w,r)
	 }else if h,ok = this["*"];ok{
		 h.ServeHTTP(w,r)
	 }else{
		 http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	 }
}
