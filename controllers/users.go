package controllers

import (
	"fmt"
	"net/http"

	"github.com/simon-lentz/webapp/models"
)

type Users struct {
	Templates struct {
		New Template // This interfaces takes the place of views.Template
	}
	UserService *models.UserService
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
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Created: +%v", user)
}
