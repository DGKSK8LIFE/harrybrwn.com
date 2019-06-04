package app

import (
	"harrybrown.com/web"
)

var (
	// Routes is a list of all the app's routes
	Routes = []web.Route{
		&web.Page{
			Title:     "Harry Brown",
			Template:  "pages/home.html",
			RoutePath: "/",
		},
		&web.Page{
			Title:     "Freelancing",
			Template:  "pages/freelance.html",
			RoutePath: "/freelance",
		},
		&web.Page{
			Title:     "Resume",
			Template:  "pages/resume.html",
			RoutePath: "/resume",
		},
	}
)

func init() {
	web.TemplateDir = "static/templates/"

	web.BaseTemplates = []string{
		"static/templates/base.html",
		"static/templates/header.html",
		"static/templates/navbar.html",
	}
}
