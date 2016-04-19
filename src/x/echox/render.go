/**
 * Copyright 2015 @ z3q.net.
 * name : render
 * author : jarryliu
 * date : 2015-12-11 17:46
 * description :
 * history :
 */
package echox

import (
	"errors"
	"github.com/labstack/echo"
	"gopkg.in/fsnotify.v1"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ echo.Renderer = new(GoTemplateForEcho)
)

type RenderWatchFunc func(echo.Renderer)

func getTemplate(dir, pattern string) (t *template.Template) {
	fi, err := os.Lstat(dir)
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		panic(errors.New("path must be directory"))
	}
	return template.Must(template.ParseGlob(dir + "/" + pattern))
}

func newGoTemplateForEcho(dir string, onWatch RenderWatchFunc) echo.Renderer {
	g := &GoTemplateForEcho{
		pattern:       "*.html",
		fileDirectory: dir,
		onWatch:       onWatch,
	}
	return g.init()
}

type GoTemplateForEcho struct {
	fileDirectory string
	pattern       string
	templates     *template.Template
	onWatch       RenderWatchFunc
}

func (g *GoTemplateForEcho) init() *GoTemplateForEcho {
	g.templates = getTemplate(g.fileDirectory, g.pattern)
	go g.fsNotify()
	return g
}

func (g *GoTemplateForEcho) fsNotify() {
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
				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					if strings.HasSuffix(event.Name, ".html") {
						log.Println("[ Template][ Update]: file - ", event.Name)
						g.init()
						if g.onWatch != nil {
							g.onWatch(g)
						}
						break
					}
				}
			case err := <-w.Errors:
				log.Println("Error:", err)
			}
		}
	}(g)

	w.Add(g.fileDirectory)
	filepath.Walk(g.fileDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return w.Add(path)
		}
		return nil
	})

	<-ch
}

func (g *GoTemplateForEcho) Render(w io.Writer, name string, data interface{}) error {
	return g.templates.ExecuteTemplate(w, name, data)
}
