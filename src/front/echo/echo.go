/**
 * Copyright 2015 @ z3q.net.
 * name : echo
 * author : jarryliu
 * date : 2015-12-04 10:51
 * description :
 * history :
 */
package echo

import (
	"container/list"
	"errors"
	"github.com/labstack/echo"
	"gopkg.in/fsnotify.v1"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HttpHosts map[string]http.Handler

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

type InterceptorFunc func(*echo.Context) bool

var (
	_                 echo.Renderer          = new(GoTemplateForEcho)
	_globTemplateData map[string]interface{} = nil
)

// 拦截器
func Interceptor(fn echo.HandlerFunc, ifn InterceptorFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if ifn(c) {
			return fn(c)
		}
		return nil
	}
}

func getTemplate(dir, pattern string) (t *template.Template, dirs *list.List) {
	dirs = new(list.List)
	fi, err := os.Lstat(dir)
	if err != nil{
		panic(err)
	}
	if !fi.IsDir() {
		panic(errors.New("path must be direction"))
	}
	t = template.Must(template.ParseGlob(dir + "/" + pattern))
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs.PushBack(path)
			if path != dir {
				t.ParseGlob(path + "/" + pattern)
			}
		}
		return nil
	})
	return t, dirs
}

func NewGoTemplateForEcho(dir string) echo.Renderer {
	g := &GoTemplateForEcho{
		pattern:       "*.html",
		fileDirectory: dir,
	}
	return g.init()
}

type GoTemplateForEcho struct {
	fileDirectory string
	pattern       string
	templates     *template.Template
}

func (g *GoTemplateForEcho) init() *GoTemplateForEcho {
	var l *list.List
	g.templates, l = getTemplate(g.fileDirectory, g.pattern)
	go g.fsNotify(l)
	return g
}

func (g *GoTemplateForEcho) fsNotify(l *list.List) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var ch chan bool = make(chan bool)
	go func(g *GoTemplateForEcho) {
		for {
			select {
			case event := <-w.Events:
				log.Println(event)
				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					if strings.HasSuffix(event.Name, ".html") {
						log.Println("[ Template][ Update]: file - ", event.Name)
						g.init()
						break
					}
				}
			case err := <-w.Errors:
				log.Println(err)
				log.Println("Error:", err)
			}
		}
	}(g)

	for itr := l.Front(); itr != nil; itr = itr.Next() {
		err = w.Add(itr.Value.(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	<-ch

}

func (g *GoTemplateForEcho) Render(w io.Writer, name string, data interface{}) error {
	return g.templates.ExecuteTemplate(w, name, data)
}

type TemplateData struct {
	Map  map[string]interface{}
	Data interface{}
}

func SetGlobRendData(m map[string]interface{}) {
	_globTemplateData = m
}

func NewRendData() *TemplateData {
	return &TemplateData{
		Map: _globTemplateData,
	}
}
