package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/simon-lentz/webapp/views"
)

var log = NewSlogger()

// HTML template http handler function.
func executeTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Debug("Failed to parse template.", slog.Any("err", err))
		http.Error(w, "Failed to parse template.", http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, nil); err != nil {
		log.Debug("Failed to execute template.", slog.Any("err", err))
		http.Error(w, "Failed to execute template.", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "home.html")
	executeTemplate(w, tmplPath)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "about.html")
	executeTemplate(w, tmplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "contact.html")
	executeTemplate(w, tmplPath)
}

// Make new router with Chi, register handler functions,
// listen and serve on port 3000.
func main() {
	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/about", aboutHandler)
	r.Get("/contact", contactHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Debug("http.ListenAndServe failed", slog.Any("err", err))
	}
}
