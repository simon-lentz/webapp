package controllers

import (
	"net/http"
)

type Users struct {
	Templates struct {
		New Template // This interfaces takes the place of views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
