package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/simon-lentz/webapp/controllers"
	"github.com/simon-lentz/webapp/migrations"
	"github.com/simon-lentz/webapp/models"
	"github.com/simon-lentz/webapp/templates"
	"github.com/simon-lentz/webapp/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, _ = strconv.Atoi(portStr)
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.CSRF.Key = os.Getenv("CSRF_KEY") // 32-byte string
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func run(cfg config) error {
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err = models.MigrateFS(db, migrations.FS, "."); err != nil {
		return err
	}
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService, _ := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}
	usersCon := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	galleriesCon := controllers.Galleries{
		GalleryService: galleryService,
	}
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"))
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
	usersCon.Templates.SignUp = views.Must(views.ParseFS(
		templates.FS,
		"sign-up.html", "tailwind.html",
	))
	usersCon.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"sign-in.html", "tailwind.html",
	))
	usersCon.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-password.html", "tailwind.html",
	))
	usersCon.Templates.CheckEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-email.html", "tailwind.html",
	))
	usersCon.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-password.html", "tailwind.html",
	))
	galleriesCon.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.html", "tailwind.html",
	))
	galleriesCon.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.html", "tailwind.html",
	))
	galleriesCon.Templates.Index = views.Must(views.ParseFS(
		templates.FS,
		"galleries/index.html", "tailwind.html",
	))
	galleriesCon.Templates.Show = views.Must(views.ParseFS(
		templates.FS,
		"galleries/show.html", "tailwind.html",
	))
	r := chi.NewRouter()
	r.Use(csrfMw, umw.SetUser)
	r.Get("/signup", usersCon.SignUp)
	r.Post("/users", usersCon.Create)
	r.Get("/signin", usersCon.SignIn)
	r.Post("/signin", usersCon.ProcessSignIn)
	r.Post("/signout", usersCon.ProcessSignOut)
	r.Get("/", homeCon)
	r.Get("/about", aboutCon)
	r.Get("/contact", contactCon)
	r.Get("/forgot-pw", usersCon.ForgotPassword)
	r.Post("/forgot-pw", usersCon.ProcessForgotPassword)
	r.Get("/reset-pw", usersCon.ResetPassword)
	r.Post("/reset-pw", usersCon.ProcessResetPassword)
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesCon.Show) // Public access route.
		r.Get("/{id}/images/{filename}", galleriesCon.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser) // Must be signed in to access these routes.
			r.Get("/", galleriesCon.Index)
			r.Post("/", galleriesCon.Create)
			r.Get("/new", galleriesCon.New)
			r.Post("/{id}", galleriesCon.Update)
			r.Get("/{id}/edit", galleriesCon.Edit)
			r.Post("/{id}/images", galleriesCon.UploadImage)
			r.Post("/{id}/delete", galleriesCon.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesCon.DeleteImage)
		})

	})
	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})
	fmt.Printf("Starting server on %s...\n", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	if err = run(cfg); err != nil {
		panic(err)
	}
}
