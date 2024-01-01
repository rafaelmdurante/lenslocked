package main

import (
	"fmt"
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

func pathHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, r.URL.Path)
	if err != nil {
		return
	}
}
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/path", pathHandler)
	fmt.Println("starting server on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		return
	}
}
