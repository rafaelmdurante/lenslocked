package controllers

import (
	"github.com/rafaelmdurante/lenslocked/views"
	"net/http"
)

type Static struct {
	Template views.Template
}

func (s Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Template.Execute(w, nil)
}
