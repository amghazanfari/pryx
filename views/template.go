package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/amghazanfari/pryx/context"
	"github.com/amghazanfari/pryx/models"
)

type ContextKey string

const CsrfFieldKey ContextKey = "csrfField"

type Template struct {
	htmlTpl *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl := template.New(pattern[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			// "csrfField": func() (template.HTML, error) {
			// 	return "", fmt.Errorf("csrfField not implemented")
			// },
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
			"firstCharCap": func() (template.HTML, error) {
				return "", fmt.Errorf("firstChar not implemented")
			},
			"showPrice": func() (template.HTML, error) {
				return "", fmt.Errorf("showPrice not implemented")
			},
			"sumActive": func() (template.HTML, error) {
				return "", fmt.Errorf("sumActive not implemented")
			},
			"sumInActive": func() (template.HTML, error) {
				return "", fmt.Errorf("sumInActive not implemented")
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
			// "csrfField": func(c *gin.Context) template.HTML {
			// 	token := csrf.GetToken(c)
			// 	return template.HTML(`<input type="hidden" name="_csrf" value="` + token + `">`)
			// },
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"firstCharCap": func(s string) string {
				r := []rune(s)
				return strings.ToUpper(string(r[0]))
			},
			"showPrice": func(price float64) string {
				var priceS string
				if price == 0 {
					priceS = "free"
				} else {
					priceS = fmt.Sprintf("%.2f$/M tokens", price)
				}
				return priceS
			},
			"sumActive": func(eps []models.Endpoint) int {
				count := 0
				for _, e := range eps {
					if e.Active {
						count++
					}
				}
				return count
			},
			"sumInActive": func(eps []models.Endpoint) int {
				count := 0
				for _, e := range eps {
					if !e.Active {
						count++
					}
				}
				return count
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
