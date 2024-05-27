package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/rafaelmdurante/lenslocked/controllers"
	"github.com/rafaelmdurante/lenslocked/migrations"
	"github.com/rafaelmdurante/lenslocked/models"
	"github.com/rafaelmdurante/lenslocked/templates"
	"github.com/rafaelmdurante/lenslocked/views"
)

func main() {
	// config database
	postgresConfig := models.DefaultPostgresConfig()

	// open connection
	db, err := models.Open(postgresConfig)
	if err != nil {
		panic(err)
	}

	// ensure connection will be closed when main function finishes
	defer db.Close()

	// run the migrations
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// set up services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// set up middleware
	// middleware to set the user from token
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// random 32-byte key
	csrfKey := "A4roiqjosijdfoi145ADSdfoqiwer"
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		// 'false' because it is not https yet, fix before deploy to prod
		csrf.Secure(false))

	// set up controllers
	users := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	users.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	users.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	users.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml"))

	// set up router and routes
	r := chi.NewRouter()

	r.Use(csrfMiddleware)
	r.Use(middleware.Logger)
	r.Use(umw.SetUser)

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

	r.Get("/signup", users.New)
	r.Post("/signup", users.Create)
	r.Get("/signin", users.SignIn)
	r.Post("/signin", users.ProcessSignIn)
	// using POST instead of DELETE because it is quite annoying to create links
	// and forms that performe the verb without the use of JavaScript
	r.Post("/signout", users.ProcessSignOut)

	// this creates sort of a namespace for the routes
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", users.CurrentUser)
	})

	r.Get("/forgot-pw", users.ForgotPassword)
	r.Post("/forgot-pw", users.ProcessForgotPassword)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("starting server on :3000...")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
