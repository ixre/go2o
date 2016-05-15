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
	"go2o/src/app/front/master"
	"go2o/src/app/front/shop/ols"
	"go2o/src/app/front/ucenter"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"gopkg.in/labstack/echo.v1"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 运行Web,监听到3个端口
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
	w := func(e echo.Renderer) { //当改动文件时,自动创建spam
		m := echox.GetGlobTemplateVars()
		m["spam"] = crypto.Md5([]byte(strconv.Itoa(int(time.Now().Unix()))))[8:14]
	}
	echox.GlobSet(m, w)
	hosts := make(MyHttpHosts)
	hosts[variable.DOMAIN_PREFIX_WEBMASTER] = master.GetServe()
	http.ListenAndServe(addr, hosts)
}

type MyHttpHosts echox.HttpHosts

func (this MyHttpHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subName := r.Host[:strings.Index(r.Host, ".")+1]
	if subName == variable.DOMAIN_PREFIX_MEMBER {
		ucenter.ServeHTTP(w, r)
	} else if h, ok := this[subName]; ok {
		h.ServeHTTP(w, r)
	} else {
		ols.ServeHTTP(w, r)
	}
}
