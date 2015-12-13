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
	"github.com/jsix/gof/crypto"
	"go2o/src/core/variable"
	"go2o/src/front/master"
	"go2o/src/front/partner"
	"go2o/src/front/shop/ols"
	"go2o/src/front/ucenter"
	"go2o/src/x/echox"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 静态文件
type StaticHandler struct {
}

func (s *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/static"+r.URL.Path)
}

// 图片处理
type ImageFileHandler struct {
	app       gof.App
	upSaveDir string
}

func (i *ImageFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path[1:4] == "res" {
		http.ServeFile(w, r, "public/static"+r.URL.Path)
	} else {
		if len(i.upSaveDir) == 0 {
			i.upSaveDir = i.app.Config().GetString(variable.UploadSaveDir)
		}
		http.ServeFile(w, r, i.upSaveDir+r.URL.Path)
	}
}

// 运行网页
func Run(ch chan bool, app gof.App, addr string) {
	defer func() {
		ch <- true
	}()
	if app.Debug() {
		log.Println("** [ Go2o][ Web][ Booted] - Web server (with debug) running on", addr)
	} else {
		log.Println("** [ Go2o][ Web][ Booted] - Web server running on", addr)
	}

	c := app.Config()
	m := map[string]interface{}{
		"static_serve": c.GetString(variable.StaticServer),
		"img_serve":    c.GetString(variable.ImageServer),
		"domain":       c.GetString(variable.ServerDomain),
		"version":      c.GetString(variable.Version),
		"spam":         crypto.Md5([]byte(strconv.Itoa(int(time.Now().Unix()))))[8:14],
	}
	echox.SetGlobRendData(m)

	hosts := make(echox.HttpHosts)
	hosts["*"] = ols.GetServe()
	hosts[variable.DOMAIN_PREFIX_MEMBER] = ucenter.GetServe()
	hosts[variable.DOMAIN_PREFIX_WEBMASTER] = master.GetServe()
	hosts[variable.DOMAIN_PREFIX_PARTNER] = partner.GetServe()
	hosts[variable.DOMAIN_PREFIX_STATIC] = new(StaticHandler)
	hosts[variable.DOMAIN_PREFIX_IMAGE] = &ImageFileHandler{app: app}

	http.ListenAndServe(addr, hosts)
}
