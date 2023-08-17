package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/simon-lentz/webapp/controllers"
	"github.com/simon-lentz/webapp/migrations"
	"github.com/simon-lentz/webapp/models"
	"github.com/simon-lentz/webapp/templates"
	"github.com/simon-lentz/webapp/views"
)

func main() {
	log := controllers.NewLogger()

	r := chi.NewRouter()
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = models.MigrateFS(db, migrations.FS, "."); err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersCon := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersCon.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.html", "tailwind.html",
	))
	r.Get("/signup", usersCon.New)
	r.Post("/users", usersCon.Create)
	usersCon.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.html", "tailwind.html",
	))
	r.Get("/signin", usersCon.SignIn)
	r.Post("/signin", usersCon.ProcessSignIn)
	r.Post("/signout", usersCon.ProcessSignOut)
	r.Get("/users/me", usersCon.CurrentUser) // For current user ONLY, otherwise will be /:id

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"home.html", "tailwind.html",
		))))
	r.Get("/about", controllers.About(
		views.Must(views.ParseFS(
			templates.FS,
			"about.html", "tailwind.html",
		))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"contact.html", "tailwind.html",
		))))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")

	csrfKey := "aInWh37hwuGH5JK8ga1fqjbLhgfANH3Q"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // Fix before deploying.
	if err := http.ListenAndServe(":3000", csrfMw(r)); err != nil {
		log.Debug("http.ListenAndServe failed", slog.Any("err", err))
	}
}
