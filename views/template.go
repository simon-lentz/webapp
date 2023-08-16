package views

import (
	"html/template"
	"io/fs"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

// Must wraps template methods that should not produce an error.
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS retrieves and parses embedded html templates.
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tmpl, err := template.ParseFS(fs, patterns...) // Expand variadic parameter.
	if err != nil {
		return Template{}, err
	}

	return Template{
		htmlTmpl: tmpl,
	}, nil
}

/*
func Parse(filepath string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTmpl: tmpl,
	}, nil
}
*/

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.htmlTmpl.Execute(w, data) // Removed error checking to make it match the controllers.Template interface
}
