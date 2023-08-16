package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template // This interfaces takes the place of views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

// Could use gorilla/schema to parse POST form vals to struct...
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprintln(w)
	fmt.Fprint(w, "Password: ", r.FormValue("password"))
}
