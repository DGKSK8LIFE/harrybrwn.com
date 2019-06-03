package app

import (
	"io"
	"html/template"
	"path/filepath"
)

var (
	tmplPath = "static/templates/"
	
	baseTemplates = []string{
		"static/templates/base.html",
		"static/templates/header.html",
		"static/templates/navbar.html",
	}
)

func execTemplate(w io.Writer, tmpl *template.Template, data interface{}) error {
	return tmpl.ExecuteTemplate(w, "base", data)
}

func newTemplate(files ...string) *template.Template {
	htmlfiles := baseTemplates
	
	for _, f := range files {
		htmlfiles = append(htmlfiles, filepath.Join(tmplPath, f))
	}

	return template.Must(template.ParseFiles(htmlfiles...))
}
