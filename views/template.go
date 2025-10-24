package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/amghazanfari/pryx/context"
	"github.com/amghazanfari/pryx/models"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl := template.New(pattern[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing filesystem: %v", err)
	}
	viewTpl := Template{
		htmlTpl: tpl,
	}
	return viewTpl, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		fmt.Printf("cloning templte: %v", err)
		http.Error(w, "there was an error rendering the template", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		fmt.Printf("executing templte: %v", err)
		http.Error(w, "there was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}
