package templates

import "html/template"

func Parse(file string) *template.Template {
	return template.Must(template.ParseFiles("internal/templates/" + file + ".html"))
}
