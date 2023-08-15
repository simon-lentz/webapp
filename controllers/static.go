package controllers

import (
	"net/http"

	"github.com/simon-lentz/webapp/views"
)

// StaticHandler generates a handler for an html template.
func StaticHandler(tmpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// closure with tmpl
		_ = tmpl.Execute(w, nil)
	}
}
