package web

import (
	_ "fmt"
	"html/template"
	"net/http"
)

//呈现模板
func RenderTemplate(w http.ResponseWriter, tplPath string, data interface{}) {
	t, err := template.ParseFiles(tplPath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
