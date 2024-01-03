package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
	if err != nil {
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, "+
		"email me at <a href=\"mailto:my@email.com\">my@email.com</a>.")
	if err != nil {
		return
	}
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, `<h1>FAQ Page</h1>
<ul>
  <li>
    <b>Is there a free version?</b>
    Yes! We offer a free trial for 30 days on any paid plans.
  </li>
  <li>
    <b>What are your support hours?</b>
    We have support staff answering emails 24/7, though response
    times may be a bit slower on weekends.
  </li>
  <li>
    <b>How do I contact support?</b>
    Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
  </li>
</ul>
`)
	if err != nil {
		return
	}
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
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.Get("/galleries/{id}", galleriesHandler)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		return
	}
}
