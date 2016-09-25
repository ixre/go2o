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
    "sync"
)

var (
    _globTemplateVar map[string]interface{} = nil
    _ echo.Renderer = new(GoTemplateForEcho)
)

type (
    Echo struct {
        *echo.Echo
        app             gof.App
        mux             sync.RWMutex
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

func (e *Echo) chkApp() {
    if e.app == nil {
        if gof.CurrentApp == nil {
            panic(errors.New("not register or no global app instance for echox!"))
        }
        e.app = gof.CurrentApp
    }
}

// 转换为Echo Handler
func (e *Echo) parseHandler(h Handler) func(ctx *echo.Context) error {
    return func(ctx *echo.Context) error {
        e.chkApp()
        return h(ParseContext(ctx, e.app))
    }
}

// 设置模板
func (e *Echo) SetTemplateRender(basePath string, notify bool, files ...string) {
    e.SetRenderer(newGoTemplateForEcho(basePath, notify, files...))
}

// 注册自定义的GET处理程序
func (e *Echo) Getx(path string, h Handler) {
    e.Get(path, e.parseHandler(h))
}

// 注册自定义的GET/POST处理程序
func (e *Echo) Anyx(path string, h Handler) {
    e.Any(path, e.parseHandler(h))
}

// 注册自定义的GET/POST处理程序
func (e *Echo) Postx(path string, h Handler) {
    e.Post(path, e.parseHandler(h))
}

func (e *Echo) getMvcHandler(route string, c *Context, obj interface{}) Handler {
    e.mux.Lock()
    defer e.mux.Unlock()
    a := c.Param("action")
    k := route + a
    if v, ok := e.dynamicHandlers[k]; ok {
        //查找路由表
        return v
    }
    if v, ok := getHandler(obj, a); ok {
        //存储路由表
        e.dynamicHandlers[k] = v
        return v
    }
    return nil
}

// 注册动态获取处理程序
// todo:?? 应复写Any
func (e *Echo) Aanyx(path string, obj interface{}) {
    h := func(c *Context) error {
        if c.Param("action") == "" {
            return c.String(http.StatusInternalServerError, "route must contain :action")
        }
        if hd := e.getMvcHandler(path, c, obj); hd != nil {
            return hd(c)
        }
        return c.String(http.StatusNotFound, "no such file")
    }
    e.Any(path, e.parseHandler(h))
}

func (e *Echo) Agetx(path string, obj interface{}) {
    h := func(c *Context) error {
        if hd := e.getMvcHandler(path, c, obj); hd != nil {
            return hd(c)
        }
        return c.String(http.StatusNotFound, "no such file")
    }
    e.Get(path, e.parseHandler(h))
}

func (e *Echo) Apostx(path string, obj interface{}) {
    h := func(c *Context) error {
        if hd := e.getMvcHandler(path, c, obj); hd != nil {
            return hd(c)
        }
        return c.String(http.StatusNotFound, "no such file")
    }
    e.Post(path, e.parseHandler(h))
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

func (e *Context) HttpResponse() http.ResponseWriter {
    return e.response
}
func (e *Context) HttpRequest() *http.Request {
    return e.request
}

func (c *Context) IsPost() bool {
    return c.request.Method == "POST"
}

func (c *Context) StringOK(s string) error {
    return c.debug(c.String(http.StatusOK, s))
}

func (c *Context) debug(err error) error {
    if err != nil {
        web.HttpError(c.HttpResponse(), err)
        return nil
    }
    return err
}

func (c *Context) Debug(err error) error {
    return c.debug(err)
}

// 覆写Render方法
func (c *Context) Render(code int, name string, data interface{}) error {
    return c.debug(c.Context.Render(code, name, data))
}

func (c *Context) RenderOK(name string, data interface{}) error {
    return c.debug(c.Render(http.StatusOK, name, data))
}

func (c *Context) NewData() *TemplateData {
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

func (g *GoTemplateForEcho) Render(w io.Writer, name string, data interface{}) error {
    return g.Execute(w, name, data)
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
