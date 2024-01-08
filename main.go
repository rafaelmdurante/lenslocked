package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rafaelmdurante/lenslocked/controllers"
	"github.com/rafaelmdurante/lenslocked/views"
	"log"
	"net/http"
	"path/filepath"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	tpl, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "there was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"))
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"))
}

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf("<p>Get gallery ID %s</p>", id)
	_, err := fmt.Fprint(w, html)
	if err != nil {
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	homeTmpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}

	r.Method(http.MethodGet, "/", controllers.Static{Template: homeTmpl})
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/galleries/{id}", galleriesHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
