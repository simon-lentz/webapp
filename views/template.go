package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/simon-lentz/webapp/context"
	"github.com/simon-lentz/webapp/errors"
	"github.com/simon-lentz/webapp/models"
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
	// Here, path.Base allows ParseFS to fetch nested embedded files.
	tmpl := template.New(path.Base(patterns[0])) // For csrf template function called in html.

	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return `<input type="hidden" />`, fmt.Errorf("csrf function not implemented.")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tmpl, err := tmpl.ParseFS(fs, patterns...) // Expand variadic parameter.
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTmpl: tmpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tmpl, err := t.htmlTmpl.Clone() // Prevents race resulting in overwriting of csrf field caused by multiple requests pointing to same template.
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)
	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer // Buffer for potential template function error.
	if err = tmpl.Execute(&buf, data); err != nil {
		log.Printf("execute: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(w, &buf); err != nil {
		log.Printf("copy buffer: %v", err)
		return
	}
}

// Alert with generic msg rather than leaking the internal error log.
func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr errors.PublicError
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err) // log unregistered error
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
