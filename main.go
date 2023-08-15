package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/simon-lentz/webapp/controllers"
	"github.com/simon-lentz/webapp/templates"
	"github.com/simon-lentz/webapp/views"
)

func main() {
	log := controllers.NewLogger()

	r := chi.NewRouter()

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
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Debug("http.ListenAndServe failed", slog.Any("err", err))
	}
}
