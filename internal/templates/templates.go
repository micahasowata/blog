package templates

import (
	"embed"
	"html/template"
)

//go:embed emails
var templates embed.FS

func Parse(file string) *template.Template {
	return template.Must(template.ParseFS(templates, "emails/"+file+".html"))
}
