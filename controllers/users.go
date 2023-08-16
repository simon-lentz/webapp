package controllers

import (
	"net/http"

	"github.com/simon-lentz/webapp/views"
)

type Users struct {
	Templates struct {
		New views.Template // Rather than a separate new type like UsersTemplates, the nested struct limits the number of exported types used.
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	_ = u.Templates.New.Execute(w, nil)
}
