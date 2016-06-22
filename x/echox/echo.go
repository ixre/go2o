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
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/session"
	"gopkg.in/labstack/echo.v1"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var (
	_globTemplateVar map[string]interface{} = nil
	_                echo.Renderer          = new(GoTemplateForEcho)
)

type (
	Echo struct {
		*echo.Echo
		app             gof.App
		dynamicHandlers map[string]Handler // 动态处理程序
	}
	Context struct {
		*echo.Context
		App      gof.App
		Session  *session.Session
		response http.ResponseWriter
		request  *http.Request
	}
	TemplateData struct {
		Var  map[string]interface{}
		Map  map[string]interface{}
		Data interface{}
	}
	Handler         func(*Context) error
	HandlerProvider interface {
		FactoryHandler(path string) *Handler
	}
)

// new echo instance
func New() *Echo {
	return &Echo{
		Echo:            echo.New(),
		dynamicHandlers: make(map[string]Handler),
	}
}

func (this *Echo) chkApp() {
	if this.app == nil {
		if gof.CurrentApp == nil {
			panic(errors.New("not register or no global app instance for echox!"))
		}
		this.app = gof.CurrentApp
	}
}

// 转换为Echo Handler
func (this *Echo) parseHandler(h Handler) func(ctx *echo.Context) error {
	return func(ctx *echo.Context) error {
		this.chkApp()
		return h(ParseContext(ctx, this.app))
	}
}

// 设置模板
func (this *Echo) SetTemplateRender(basePath string, notify bool, files ...string) {
	this.SetRenderer(newGoTemplateForEcho(basePath, notify, files...))
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

func (this *Echo) getMvcHandler(route string, c *Context, obj interface{}) Handler {
	a := c.Param("action")
	k := route + a
	if v, ok := this.dynamicHandlers[k]; ok {
		//查找路由表
		return v
	}
	if v, ok := getHandler(obj, a); ok {
		//存储路由表
		this.dynamicHandlers[k] = v
		return v
	}
	return nil
}

// 注册动态获取处理程序
// todo:?? 应复写Any
func (this *Echo) Aanyx(path string, obj interface{}) {
	h := func(c *Context) error {
		if hd := this.getMvcHandler(path, c, obj); hd != nil {
			return hd(c)
		}
		return c.String(http.StatusInternalServerError, "no such file")
	}
	this.Any(path, this.parseHandler(h))
}

func (this *Echo) Agetx(path string, obj interface{}) {
	h := func(c *Context) error {
		if hd := this.getMvcHandler(path, c, obj); hd != nil {
			return hd(c)
		}
		return c.String(http.StatusInternalServerError, "no such file")
	}
	this.Get(path, this.parseHandler(h))
}

func (this *Echo) Apostx(path string, obj interface{}) {
	h := func(c *Context) error {
		if hd := this.getMvcHandler(path, c, obj); hd != nil {
			return hd(c)
		}
		return c.String(http.StatusInternalServerError, "no such file")
	}
	this.Post(path, this.parseHandler(h))
}

func ParseContext(ctx *echo.Context, app gof.App) *Context {
	req, rsp := ctx.Request(), ctx.Response()
	s := session.Default(rsp, req)
	return &Context{
		Context:  ctx,
		Session:  s,
		App:      app,
		response: rsp,
		request:  req,
	}
}

func (this *Context) HttpResponse() http.ResponseWriter {
	return this.response
}
func (this *Context) HttpRequest() *http.Request {
	return this.request
}

func (this *Context) StringOK(s string) error {
	return this.debug(this.String(http.StatusOK, s))
}

func (this *Context) debug(err error) error {
	if err != nil {
		web.HttpError(this.HttpResponse(), err)
		return nil
	}
	return err
}

func (this *Context) Debug(err error) error {
	return this.debug(err)
}

// 覆写Render方法
func (this *Context) Render(code int, name string, data interface{}) error {
	return this.debug(this.Context.Render(code, name, data))
}

func (this *Context) RenderOK(name string, data interface{}) error {
	return this.debug(this.Render(http.StatusOK, name, data))
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

// 全局设定ECHO参数
func GlobSet(globVars map[string]interface{}) {
	_globTemplateVar = globVars
}

// 获取全局模版变量
func GetGlobTemplateVars() map[string]interface{} {
	return _globTemplateVar
}

func newGoTemplateForEcho(dir string, notify bool, files ...string) echo.Renderer {
	return &GoTemplateForEcho{
		CachedTemplate: gof.NewCachedTemplate(dir, notify, files...),
	}
}

type GoTemplateForEcho struct {
	*gof.CachedTemplate
}

func (this *GoTemplateForEcho) Render(w io.Writer, name string, data interface{}) error {
	return this.Execute(w, name, data)
}

//
//type InterceptorFunc func(echo.Context) bool
//
//// 拦截器
//func Interceptor(fn echo.HandlerFunc, ifn InterceptorFunc) echo.HandlerFunc {
//	return func(c *echo.Context) error {
//		if ifn(c) {
//			return fn(c)
//		}
//		return nil
//	}
//}

/****************  MIDDLE WAVE ***************/

var (
	requestFilter = map[string]*regexp.Regexp{
		"GET": regexp.MustCompile("'|(and|or)\\b.+?(>|<|=|in|like)|\\/\\*.+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION" +
			".+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|(CREATE|ALTER|DROP|TRUNCATE)\\s+" +
			"(TABLE|DATABASE)"),
		"POST": regexp.MustCompile("\\b(and|or)\\b.{1,6}?(=|>|<|\\bin\\b|\\blike\\b)|\\/\\*" +
			".+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION.+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|" +
			"(CREATE|ALTER|DROP|TRUNCATE)\\s+(TABLE|DATABASE)"),
	}

	/*
	   getFilter = postFilter = cookieFilter = regexp.MustCompile("\\b(and|or)\\b.{1,6}?(=|>|<|\\bin\\b|\\blike\\b)|\\/\\*.+?\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION.+?SELECT|UPDATE.+?SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE).+?FROM|(CREATE|ALTER|DROP|TRUNCATE)\\s+(TABLE|DATABASE)");
	*/
)

// 防SQL注入
func StopAttack(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		badRequest := false
		method := ctx.Request().Method
		switch method {
		case "GET":
			badRequest = requestFilter[method].MatchString(ctx.Request().URL.RawQuery)
		case "POST":
			badRequest = requestFilter["GET"].MatchString(ctx.Request().URL.RawQuery) ||
				requestFilter[method].MatchString(
					ctx.Request().Form.Encode())
		}
		if badRequest {
			return ctx.HTML(http.StatusNotFound,
				"<div style='color:red;'>您提交的参数非法,系统已记录您本次操作!</div>")
		}
		return h(ctx)
	}
}
