package app

import (
	"net/http"

	"harrybrown.com/web"
)

// HandleHome serves the homepage
func HandleHome(w http.ResponseWriter, r *http.Request) {
	home := &Page{
		Title:        "Harry Brown",
		BodyTemplate: "pages/home.html",
	}

	if err := home.Write(w); err != nil {
		web.NotFound(w, r)
		return
	}
}
