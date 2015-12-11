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
	"github.com/labstack/echo"
	"html/template"
	"container/list"
	"os"
	"path/filepath"
	"errors"
	"gopkg.in/fsnotify.v1"
	"log"
	"strings"
	"io"
)

var (
	_ echo.Renderer = new(GoTemplateForEcho)
)

func getTemplate(dir, pattern string) (t *template.Template, dirs *list.List) {
	dirs = new(list.List)
	fi, err := os.Lstat(dir)
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		panic(errors.New("path must be directory"))
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
				if event.Op & fsnotify.Write != 0 || event.Op & fsnotify.Create != 0 {
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


