package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)
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
