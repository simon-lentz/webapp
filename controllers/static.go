package controllers

import (
	"html/template"
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

func About(tmpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML // Hardcoded HTML does not need to be escaped.
	}{{
		Question: "Question One",
		Answer:   "Answer One",
	},
		{
			Question: "Question Two",
			Answer:   "Answer Two",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmpl.Execute(w, questions)
	}
}
