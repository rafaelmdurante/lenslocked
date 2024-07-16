package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rafaelmdurante/lenslocked/context"
	"github.com/rafaelmdurante/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

type public interface {
	Public() string
}

// Execute renders the page and takes errors to show to user
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// ensure every incoming http request has their own template to work with
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "there was an error rendering the page", http.StatusInternalServerError)
		return
	}

	errMessages := errMessages(errs...)

	// use the gorilla/csrf package to generate the csrf html code
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMessages
			},
		},
	)

	// set header and execute the template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// avoid superfluous response.WriteHeader error
	// this happens when tpl.Execute begins executing the template and writing
	// the results to http.ResponseWriter, but if no errors are found, http.rw
	// sets the header to 200, later if an error occurs, the status code won't
	// change because it could have already responded to the client
	// one way to avoid the superfluous error message si to buffer the results
	// from the template execution
	var b bytes.Buffer
	err = tpl.Execute(&b, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "there was an error executing the template",
			http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		log.Printf("executing template copying buffer: %v", err)
		http.Error(w, "there was an error executing the template",
			http.StatusInternalServerError)
	}
}

func errMessages(errs ...error) []string {
	var messages []string
	for _, err := range errs {
		var publicError public
		if errors.As(err, &publicError) {
			messages = append(messages, publicError.Public())
		} else {
			fmt.Println(err)
			messages = append(messages, "Something went wrong.")
		}
	}

	return messages
}

func ParseFS(filesystem fs.FS, pattern ...string) (Template, error) {
	// create new template named after the first gohtml file
	htmlTpl := template.New(pattern[0])
	// declare the function into the template
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		})

	htmlTpl, err := htmlTpl.ParseFS(filesystem, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

// Must function is based on the Go's template package same name function
func Must(t Template, err error) Template {
	if err != nil {
		// in this case we do want to panic as there is no point in starting
		// the server if a template fails to be parsed
		panic(err)
	}

	return t
}
