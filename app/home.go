package app

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"harrybrown.com/web"
)

func init() {
	web.BaseTemplates = []string{
		"static/templates/base.html",
		"static/templates/header.html",
		"static/templates/navbar.html",
	}
	web.TemplateDir = "static/templates/"
}

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
			Data:      getResumeContent("./static/data/resume.json"),
		},
	}
)

func getResumeContent(file string) *resumeContent {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}
	c := &resumeContent{}

	if err = json.Unmarshal(b, c); err != nil {
		log.Println(err)
		return nil
	}
	return c
}

type resumeContent struct {
	Experience []resumeItem
	Education  []resumeItem
}

type resumeItem struct {
	Name, Title, Date, Content string
	BulletPoints               []string
}
