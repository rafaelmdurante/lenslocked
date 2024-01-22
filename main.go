package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rafaelmdurante/lenslocked/controllers"
	"github.com/rafaelmdurante/lenslocked/models"
	"github.com/rafaelmdurante/lenslocked/templates"
	"github.com/rafaelmdurante/lenslocked/views"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// config routes
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml"))))

	// config database
	postgresConfig := models.DefaultPostgresConfig()

	// open connection
	db, err := models.Open(postgresConfig)
	if err != nil {
		panic(err)
	}

	// ensure connection will be closed when main function finishes
	defer db.Close()

	// create user service
	userService := models.UserService{
		DB: db,
	}

	users := controllers.Users{
		UserService: &userService, //TODO: set this!
	}
	users.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", users.New)
	r.Post("/signup", users.Create)

	users.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", users.SignIn)
	r.Post("/signin", users.ProcessSignIn)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
