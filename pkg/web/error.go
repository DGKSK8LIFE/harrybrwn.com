package web

import (
	"fmt"
	"html/template"
	"net/http"
)

// NotFound handles requests that generate a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
	var tmplNotFound = template.Must(template.ParseFiles(
		"static/templates/pages/404.html",
		"static/templates/index.html",
	))
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := tmplNotFound.ExecuteTemplate(w, "base", nil); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
