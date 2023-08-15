package views

import (
	"fmt"
	"html/template"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

func Parse(filepath string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTmpl: tmpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.htmlTmpl.Execute(w, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}
	return nil
}
