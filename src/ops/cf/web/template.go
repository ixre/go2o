package web

import (
	"html/template"
	"io"
)

type TemplateWrapper struct {
	Init func(m *map[string]interface{})
}

// the data map for template
type TemplateMapFunc func(m *map[string]interface{})

// execute single template file
func (this *TemplateWrapper) Render(w io.Writer, tplPath string, f TemplateMapFunc,
) error {
	return this.Execute(w, f, tplPath)
}

// execute template
func (this *TemplateWrapper) Execute(w io.Writer, f TemplateMapFunc,
	tplPath ...string) error {

	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		//http.Error(weixin, err.Error(),500)
		return err
	}

	data := make(map[string]interface{})
	if this.Init != nil {
		this.Init(&data)
	}
	if f != nil {
		f(&data)
	}

	t.Execute(w, &data)
	return nil
}
