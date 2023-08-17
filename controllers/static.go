package controllers

import (
	"html/template"
	"net/http"
)

// StaticHandler generates a handler for an html template.
func StaticHandler(tmpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// closure with tmpl
		tmpl.Execute(w, r, nil)
	}
}

func About(tmpl Template) http.HandlerFunc {
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
		tmpl.Execute(w, r, questions)
	}
}
