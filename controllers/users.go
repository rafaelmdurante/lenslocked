package controllers

import (
	"fmt"
	"net/http"

	"github.com/rafaelmdurante/lenslocked/context"
	"github.com/rafaelmdurante/lenslocked/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong when creating a user",
			http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: long term, we should show a warning about not being able to sign the user in
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

// ProcessSignIn authenticates users and sets the necessary cookies for
// authentication
func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong authenticating", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong creating session", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong deleting session", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

// CurrentUser gets the current user from the cookie
// SetUser and RequireUser middleware are required
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "current user: %s\n", user.Email)
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First try to read the cookie. If we run into an error reading it,
		// proceed with the request. The goal of this middleware isn't to limit
		// access. It only sets the user in the context if it can.
		token, err := readCookie(r, CookieSession)
		if err != nil {
			// cannot lookup the user with no cookie, so proceed without a user
			// being set, then return
			next.ServeHTTP(w, r)
			return
		}

		// if we have a token, try to lookup the user with that token
		user, err := umw.SessionService.User(token)
		if err != nil {
			// invalid or expired token, then proceed without setting a user
			next.ServeHTTP(w, r)
			return
		}

		// if we get to this point, we have a user that we can store in the ctx
		ctx := r.Context()
		// we need to derive a new context to store values in it. be certain
		// thath we import our own context package, and not the one from the
		// standard library.
		ctx = context.WithUser(ctx, user)
		// get a request that uses our new context. this is done in a way
		// similar to how contexts work - we call a withcontext function
		// and it returns us a new request with the context set
		r = r.WithContext(ctx)
		// call the handler that our middleware was applied to with the updated
		next.ServeHTTP(w, r)
	})
}

// RequireUser is a middleware to get the user from the context.
// It assumes the SetUser middleware has been run BEFORE, so it doesn't need
// to perform the same database lookups.
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
        next.ServeHTTP(w, r)
	})
}
