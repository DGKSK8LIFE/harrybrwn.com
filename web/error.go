package web

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles(
	"static/templates/pages/404.html",
	"static/templates/index.html",
))

// NotFound handles requests that generate a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
