package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	// ensure every incoming http request has their own template to work with
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "there was an error rendering the page", http.StatusInternalServerError)
		return
	}

	// use the gorilla/csrf package to generate the csrf html code
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
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

func ParseFS(filesystem fs.FS, pattern ...string) (Template, error) {
	// create new template named after the first gohtml file
	htmlTpl := template.New(pattern[0])
	// declare the function into the template
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
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
