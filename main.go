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

	// Set up DB.
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg.String())
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = models.MigrateFS(db, migrations.FS, "."); err != nil {
		panic(err)
	}

	// Set up services, associate with user controller.
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

	// Set up middleware.
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	csrfKey := "aInWh37hwuGH5JK8ga1fqjbLhgfANH3Q"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // Fix before deploying.

	// Set up controller handler functions.
	homeCon := controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"home.html", "tailwind.html",
		)))
	aboutCon := controllers.About(
		views.Must(views.ParseFS(
			templates.FS,
			"about.html", "tailwind.html",
		)))
	contactCon := controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"contact.html", "tailwind.html",
		)))

	// Register html templates.
	usersCon.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.html", "tailwind.html",
	))
	usersCon.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.html", "tailwind.html",
	))
	usersCon.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.html", "tailwind.html",
	))

	// Set up router, associate routes with their respective handler functions.
	r := chi.NewRouter()
	r.Use(csrfMw, umw.SetUser)
	r.Get("/signup", usersCon.New)
	r.Post("/users", usersCon.Create)
	r.Get("/signin", usersCon.SignIn)
	r.Post("/signin", usersCon.ProcessSignIn)
	r.Post("/signout", usersCon.ProcessSignOut)
	r.Get("/", homeCon)
	r.Get("/about", aboutCon)
	r.Get("/contact", contactCon)
	r.Get("/forgot-pw", usersCon.ForgotPassword)
	r.Post("/forgot-pw", usersCon.ProcessForgotPassword)
	// r.Get("/users/me", usersCon.CurrentUser)
	// Use subrouting for the context-dependent middleware.
	r.Route("/users/me", func(r chi.Router) { // Subroute that requires user to be signed in.
		r.Use(umw.RequireUser)
		r.Get("/", usersCon.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	// Start server.
	fmt.Println("Starting server on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Debug("http.ListenAndServe failed", slog.Any("err", err))
	}
}
