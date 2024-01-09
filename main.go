package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rafaelmdurante/lenslocked/controllers"
	"github.com/rafaelmdurante/lenslocked/views"
	"net/http"
	"path/filepath"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(
			views.Parse(filepath.Join("templates", "home.gohtml")))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(
			views.Parse(filepath.Join("templates", "contact.gohtml")))))

	r.Get("/faq", controllers.StaticHandler(
		views.Must(
			views.Parse(filepath.Join("templates", "faq.gohtml")))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000...")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
