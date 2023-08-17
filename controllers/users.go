package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/simon-lentz/webapp/models"
)

type Users struct {
	Templates struct {
		New    Template // This interfaces takes the place of views.Template
		SignIn Template // SignIn contains logic for an existing user.
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}

	data.Email = r.FormValue("email")
	data.CSRFField = csrf.TemplateField(r)
	u.Templates.New.Execute(w, data)
}

// Could use gorilla/schema to parse form.
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Created: %+v", user)
}

// SignIn is the handler func for the /signin route.
func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong.", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true, // No js access, xss protection.
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User Authenticated: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	emailCookie, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "Email cookie could not be read.")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s\n", emailCookie.Value)

}
