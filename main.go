package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/simon-lentz/webapp/controllers"
	"github.com/simon-lentz/webapp/views"
)

var log = NewSlogger()

// Make new router with Chi, register handler functions,
// listen and serve on port 3000.
func main() {
	r := chi.NewRouter()

	tmpl, err := views.Parse(filepath.Join("templates", "home.html"))
	if err != nil {
		log.Debug("Error parsing template", slog.Any("err", err))
	}
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse(filepath.Join("templates", "about.html"))
	if err != nil {
		log.Debug("Error parsing template", slog.Any("err", err))
	}
	r.Get("/about", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse(filepath.Join("templates", "contact.html"))
	if err != nil {
		log.Debug("Error parsing template", slog.Any("err", err))
	}
	r.Get("/contact", controllers.StaticHandler(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Debug("http.ListenAndServe failed", slog.Any("err", err))
	}
}
