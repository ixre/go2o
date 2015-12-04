/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:47
 * description :
 * history :
 */

package app

import (
	"github.com/jsix/gof"
	"go2o/src/core/variable"
	"net/http"
	"go2o/src/front"
	"go2o/src/external/ols"
	"go2o/src/front/ucenter"
	"go2o/src/front/master"
	"go2o/src/front/partner"
)

// 静态文件
type StaticHandler struct{

}
func (s *StaticHandler) ServeHTTP(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w, r, "./static"+r.URL.Path)
}

// 图片处理
type ImageFileHandler struct{
}
func (i *ImageFileHandler) ServeHTTP(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w, r, "./static/uploads/"+r.URL.Path)
}

// 运行网页
func Run(ch chan bool,app gof.App, addr string) {
	defer func(){
		ch <- true
	}()
	if app.Debug() {
		app.Log().Println("** [ Go2o][ Web][ Booted] - Web server (with debug) running on",addr)
	} else {
		app.Log().Println("** [ Go2o][ Web][ Booted] - Web server running on",addr)
	}

	hosts := make(front.HttpHosts)
 	hosts["*"] = ols.GetServe()
	hosts[variable.DOMAIN_PREFIX_MEMBER] = ucenter.GetServe()
	hosts[variable.DOMAIN_PREFIX_WEBMASTER] = master.GetServe()
	hosts[variable.DOMAIN_PREFIX_PARTNER] = partner.GetServe()
	hosts[variable.DOMAIN_PREFIX_STATIC] =new(StaticHandler)
	hosts[variable.DOMAIN_PREFIX_IMAGE] = new(ImageFileHandler)

	http.ListenAndServe(addr,hosts)
}
