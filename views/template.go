package views

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
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

	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "there was an error executing the template",
			http.StatusInternalServerError)
		return
	}
}

func ParseFS(filesystem fs.FS, pattern ...string) (Template, error) {
	// create new template named after the first gohtml file
	htmlTpl := template.New(pattern[0])
	// declare the function into the template
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<input type="hidden" />`
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
