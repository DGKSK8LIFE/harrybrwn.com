package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"harrybrown.com/web"
)

func init() {
	web.TemplateDir = "static/templates/"
	web.BaseTemplates = []string{"/index.html", "/nav.html"} // included in all pages
}

var (
	// Routes is a list of all the app's routes
	Routes = []web.Route{
		&web.Page{
			Title:     "Harry Brown",
			Template:  "pages/home.html",
			RoutePath: "/",
			Data: struct{ Age string }{
				Age: fmt.Sprintf("%d years", int((time.Since(bday).Hours()/24)/365)),
			},
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

	bday = time.Date(1998, time.August, 4, 4, 0, 0, 0, time.UTC)
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