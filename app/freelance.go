package app

import (
	"net/http"
	
	"harrybrown.com/web"
)

// Freelance handles the freelance section of the website
func Freelance(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title:        "Freelancing",
		BodyTemplate: "pages/freelance.html",
	}
	
	if err := page.Write(w); err != nil {
		web.NotFound(w, r)
		return
	}

	// templates := newTemplate("pages/freelance.html")
	// 	if err := execTemplate(w, templates, nil); err != nil {
	// }
}