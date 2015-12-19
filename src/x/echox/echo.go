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
	"errors"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
	"strings"
)

var (
	_globApp         gof.App
	_globTemplateVar map[string]interface{} = nil
	_globRenderWatch RenderWatchFunc
)

type (
	Echo struct {
		*echo.Echo
		app             gof.App
		dynamicHandlers map[string]Handler // 动态处理程序
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
	Handler   func(*Context) error
	HttpHosts map[string]http.Handler
)

// new echo instance
func New() *Echo {
	if _globApp == nil {
		_globApp = gof.CurrentApp
	}
	return &Echo{
		Echo:            echo.New(),
		app:             _globApp,
		dynamicHandlers: make(map[string]Handler),
	}
}

func (this *Echo) parseHandler(h Handler) func(ctx *echo.Context) error {
	return func(ctx *echo.Context) error {
		s := session.Default(ctx.Response(), ctx.Request())
		return h(&Context{
			Context: ctx,
			Session: s,
			App:     this.app,
		})
	}
}

// 设置模板
func (this *Echo) SetTemplateRender(path string) {
	this.SetRenderer(newGoTemplateForEcho(path, _globRenderWatch))
}

// 注册自定义的GET处理程序
func (this *Echo) Getx(path string, h Handler) {
	this.Get(path, this.parseHandler(h))
}

// 注册自定义的GET/POST处理程序
func (this *Echo) Anyx(path string, h Handler) {
	this.Any(path, this.parseHandler(h))
}

// 注册自定义的GET/POST处理程序
func (this *Echo) Postx(path string, h Handler) {
	this.Post(path, this.parseHandler(h))
}

// 注册动态获取处理程序
func (this *Echo) Danyx(path string, v interface{}) {
	h := func(ctx *Context) error {
		a := ctx.Param("action")
		k := path + a
		if v, ok := this.dynamicHandlers[k]; ok {
			return v(ctx)
		}
		if v, ok := getHandler(v, a); ok {
			this.dynamicHandlers[k] = v
			return v(ctx)
		}
		return errors.New("no such action named : " + a)
	}
	this.Any(path, this.parseHandler(h))
}

func (this *Context) StringOK(s string) error {
	return this.String(http.StatusOK, s)
}

func (this *Context) RenderOK(name string, data interface{}) error {
	return this.Render(http.StatusOK, name, data)
}

func (this *Context) NewData() *TemplateData {
	return &TemplateData{
		Var:  _globTemplateVar,
		Map:  make(map[string]interface{}),
		Data: nil,
	}
}

// get handler by reflect
func getHandler(v interface{}, action string) (Handler, bool) {
	t := reflect.ValueOf(v)
	method := t.MethodByName(strings.Title(action))
	if method.IsValid() {
		v, ok := method.Interface().(func(*Context) error)
		return v, ok
	}
	return nil, false
}

func (this HttpHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subName := r.Host[:strings.Index(r.Host, ".")+1]
	if h, ok := this[subName]; ok {
		h.ServeHTTP(w, r)
	} else if h, ok = this["*"]; ok {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// 全局设定ECHO参数
func GlobSet(globVars map[string]interface{}, w RenderWatchFunc) {
	_globTemplateVar = globVars
	_globRenderWatch = w
}

// 获取全局模版变量
func GetGlobTemplateVars() map[string]interface{} {
	return _globTemplateVar
}

func NewRenderData() *TemplateData {
	return &TemplateData{
		Var:  _globTemplateVar,
		Map:  make(map[string]interface{}),
		Data: nil,
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
