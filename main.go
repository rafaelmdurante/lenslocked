package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/rafaelmdurante/lenslocked/controllers"
	"github.com/rafaelmdurante/lenslocked/migrations"
	"github.com/rafaelmdurante/lenslocked/models"
	"github.com/rafaelmdurante/lenslocked/templates"
	"github.com/rafaelmdurante/lenslocked/views"
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

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: PSQL
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: STMP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}

	// TODO: CSRF
	cfg.CSRF.Key = "A4roiqjosijdfoi145ADSdfoqiwer"
	cfg.CSRF.Secure = false

	// TODO: Server
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// config database
	// open connection
	db, err := models.Open(cfg.PSQL)
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
	pwResetService := models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := models.GalleryService{
		DB: db,
	}

	// set up middleware
	// middleware to set the user from token
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// random 32-byte key
	csrfMiddleware := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		// 'false' because it is not https yet, fix before deploy to prod
		csrf.Secure(cfg.CSRF.Secure))

	// set up controllers
	users := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}
	users.Templates.New = views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	users.Templates.SignIn = views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	users.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml"))
	users.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml"))
	users.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml"))

	// galleries controllers
	galleries := controllers.Galleries{
		GalleryService: &galleryService,
	}
	galleries.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))

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

	r.Get("/reset-pw", users.ResetPassword)
	r.Post("/reset-pw", users.ProcessResetPassword)

	r.Get("/galleries/new", galleries.New)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Printf("starting server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
