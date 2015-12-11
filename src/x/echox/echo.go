/**
 * Copyright 2015 @ z3q.net.
 * name : echo
 * author : jarryliu
 * date : 2015-12-04 10:51
 * description :
 * history :
 */
package echox

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web/session"
)

var (
	globalApp gof.App
	_globTemplateData map[string]interface{} = nil
)

type(
	Echo struct {
		*echo.Echo
		app gof.App
	}
	Context struct {
		*echo.Context
		App     gof.App
		Session *session.Session
	}
	TemplateData struct {
		Var  map[string]interface{}
		Map  map[string]interface{}
		Data interface{}
	}
	Handler func(*Context) error
	HttpHosts map[string]http.Handler
)



// new echo instance
func New() *Echo {
	if globalApp == nil {
		globalApp = gof.CurrentApp
	}
	return &Echo{
		Echo:echo.New(),
		app:globalApp,
	}
}


func (this *Echo) parseHandler(h Handler) func(ctx *echo.Context) error {
	return func(ctx *echo.Context) error {
		s := session.Default(ctx.Response(), ctx.Request())
		return h(&Context{
			Context:ctx,
			Session:s,
			App:this.app,
		})
	}
}

func (this *Echo) Getx(path string, h Handler) {
	this.Get(path, this.parseHandler(h))
}


func (this HttpHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subName := r.Host[:strings.Index(r.Host, ".") + 1]
	if h, ok := this[subName]; ok {
		h.ServeHTTP(w, r)
	} else if h, ok = this["*"]; ok {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}




func SetGlobRendData(m map[string]interface{}) {
	_globTemplateData = m
}

func NewRendData() *TemplateData {
	return &TemplateData{
		Var: _globTemplateData,
		Map:make(map[string]interface{}),
		Data:nil,
	}
}


type InterceptorFunc func(*echo.Context) bool

// 拦截器
func Interceptor(fn echo.HandlerFunc, ifn InterceptorFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if ifn(c) {
			return fn(c)
		}
		return nil
	}
}
